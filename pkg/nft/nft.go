package nft

import (
	"encoding/json"
	"fmt"

	binarycodec "github.com/CreatureDev/xrpl-go/binary-codec"
	"github.com/CreatureDev/xrpl-go/client"
	"github.com/CreatureDev/xrpl-go/keypairs"
	"github.com/CreatureDev/xrpl-go/model/client/account"
	"github.com/CreatureDev/xrpl-go/model/client/common"
	ledgerCl "github.com/CreatureDev/xrpl-go/model/client/ledger"
	"github.com/CreatureDev/xrpl-go/model/client/path"
	txCl "github.com/CreatureDev/xrpl-go/model/client/transactions"
	"github.com/CreatureDev/xrpl-go/model/ledger"
	"github.com/CreatureDev/xrpl-go/model/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
	"golang.org/x/exp/slices"
)

type NFTokenURIGenerator func() types.NFTokenURI

func ConstantURIGenerator(uri types.NFTokenURI) NFTokenURIGenerator {
	return func() types.NFTokenURI {
		return uri
	}
}

type NFTokenMintConfig struct {
	NFTokenTaxon uint
	TransferFee  uint16
	URIGenerator NFTokenURIGenerator
	Flags        *types.Flag
	Timeout      uint
	Price        types.CurrencyAmount
}

type NFTokenMinter struct {
	client  *client.XRPLClient
	account types.Address
	privKey string
	pubKey  string
}

func (m *NFTokenMinter) GetAccount() types.Address {
	return m.account
}

func EmptyMinter() *NFTokenMinter {
	return &NFTokenMinter{}
}

func (m *NFTokenMinter) WithClient(cl *client.XRPLClient) *NFTokenMinter {
	m.client = cl
	return m
}

func (m *NFTokenMinter) WithAccount(acc types.Address) *NFTokenMinter {
	m.account = acc
	return m
}

func (m *NFTokenMinter) WithKeys(priv, pub string) *NFTokenMinter {
	m.privKey = priv
	m.pubKey = pub
	return m
}

func (m *NFTokenMinter) WithPublicKey(pub string) *NFTokenMinter {
	m.pubKey = pub
	return m
}

func (m *NFTokenMinter) WithPrivateKey(priv string) *NFTokenMinter {
	m.privKey = priv
	return m
}

func NewNFTokenMinter(cl *client.XRPLClient, acc types.Address, priv string, pub string) *NFTokenMinter {
	return &NFTokenMinter{
		client:  cl,
		account: acc,
		privKey: priv,
		pubKey:  pub,
	}
}

func (minter *NFTokenMinter) NFTIsValid(id types.NFTokenID) bool {
	offers, _ := minter.SellOffersForNFT(id)
	return len(offers) == 0
}

func (minter *NFTokenMinter) SellOffersForNFT(id types.NFTokenID) ([]path.NFTokenOffer, error) {
	offersRequest := path.NFTokenSellOffersRequest{
		NFTokenID: id,
	}
	offersResponse, _, err := minter.client.Path.NFTokenSellOffers(&offersRequest)
	if err != nil {
		return nil, fmt.Errorf("fetch sell offers: %w", err)
	}
	offers := offersResponse.Offers
	for offersResponse.Marker != nil {
		offersRequest.Marker = offersResponse.Marker
		offersResponse, _, err := minter.client.Path.NFTokenSellOffers(&offersRequest)
		if err != nil {
			return nil, fmt.Errorf("fetch pages sell offers: %w", err)
		}
		offers = append(offers, offersResponse.Offers...)
	}
	return offers, nil
}

func (minter *NFTokenMinter) CancelExpiredNFTSales(configs []NFTokenMintConfig) error {
	nftsRequest := account.AccountNFTsRequest{
		Account: minter.account,
	}

	nftsResponse, _, err := minter.client.Account.AccountNFTs(&nftsRequest)
	if err != nil {
		return err
	}
	nfts := nftsResponse.AccountNFTs
	for nftsResponse.Marker != nil {
		nftsRequest.Marker = nftsResponse.Marker
		nftsResponse, _, err = minter.client.Account.AccountNFTs(&nftsRequest)
		nfts = append(nfts, nftsResponse.AccountNFTs...)
	}
	var taxons []uint
	for _, c := range configs {
		taxons = append(taxons, c.NFTokenTaxon)
	}

	var expired []types.Hash256
	for _, n := range nfts {
		if !slices.Contains(taxons, n.NFTokenTaxon) {
			continue
		}
		offers, err := minter.SellOffersForNFT(n.NFTokenID)
		if err != nil {
			return err
		}
		for _, nftOffer := range offers {
			entryRequest := ledgerCl.LedgerEntryRequest{
				Offer: ledgerCl.EntryString(nftOffer.NFTokenOfferIndex),
			}
			entryResponse, _, err := minter.client.Ledger.LedgerEntry(&entryRequest)
			if err != nil {
				return err
			}
			obj, ok := entryResponse.Node.(*ledger.NFTokenOffer)
			if !ok {
				return fmt.Errorf("unexpected nft offer format")
			}
			if obj.Expiration < common.CurrentRippleTime() {
				expired = append(expired, types.Hash256(nftOffer.NFTokenOfferIndex))
			}
		}
	}
	cancelOfferTx := transactions.NFTokenCancelOffer{
		BaseTx: transactions.BaseTx{
			Account: minter.account,
		},
		NFTokenOffers: expired,
	}
	// TODO resubmit on failed TX attempts
	// make submit helper function (autofills/signs until success or too many failures)
	minter.client.AutofillTx(minter.account, &cancelOfferTx)
	blob, _ := binarycodec.EncodeForSigning(&cancelOfferTx)
	sig, _ := keypairs.Sign(blob, minter.privKey)
	cancelOfferTx.BaseTx.TxnSignature = sig
	tx, _ := binarycodec.Encode(&cancelOfferTx)
	submitReq := txCl.SubmitRequest{
		TxBlob: tx,
	}
	_, _, err = minter.client.Transaction.Submit(&submitReq)
	return err
}

func (minter *NFTokenMinter) GetValidNFT(conf NFTokenMintConfig) types.NFTokenID {
	cl := minter.client
	accNftsReq := &account.AccountNFTsRequest{
		Account: minter.account,
	}

	accNftsRes, _, err := cl.Account.AccountNFTs(accNftsReq)
	if err != nil {
		return ""
	}

	for _, nft := range accNftsRes.AccountNFTs {
		if nft.Issuer != minter.account {
			continue
		}
		if nft.NFTokenTaxon == conf.NFTokenTaxon && minter.NFTIsValid(nft.NFTokenID) {
			return nft.NFTokenID
		}
	}

	return ""
}

func (minter *NFTokenMinter) GetOrMintNFT(conf NFTokenMintConfig) (types.NFTokenID, error) {
	if id := minter.GetValidNFT(conf); id != "" {
		return id, nil
	}
	cl := minter.client

	mintTx := transactions.NFTokenMint{
		BaseTx: transactions.BaseTx{
			Account:         minter.account,
			TransactionType: transactions.NFTokenMintTx,
			Flags:           conf.Flags,
			SigningPubKey:   minter.pubKey,
		},
		NFTokenTaxon: conf.NFTokenTaxon,
		TransferFee:  conf.TransferFee,
		URI:          conf.URIGenerator(),
	}
	// TODO resubmit on failed TX attempts
	if err := cl.AutofillTx(minter.account, &mintTx); err != nil {
		return "", fmt.Errorf("autofill mint tx: %w", err)
	}
	bin, err := binarycodec.EncodeForSigning(&mintTx)
	if err != nil {
		return "", fmt.Errorf("encoding mint tx for signing: %w", err)
	}
	sig, _ := keypairs.Sign(bin, minter.privKey)
	mintTx.BaseTx.TxnSignature = sig
	blob, err := binarycodec.Encode(&mintTx)
	if err != nil {
		return "", fmt.Errorf("encoding mint tx: %w", err)
	}
	submitReq := txCl.SubmitRequest{
		TxBlob: blob,
	}
	submitRes, _, err := cl.Transaction.Submit(&submitReq)
	if err != nil {
		return "", fmt.Errorf("minting nft: %w", err)
	}
	if !submitRes.Accepted {
		return "", fmt.Errorf("mint request was not accepted")
	}
	id := minter.GetValidNFT(conf)
	if id == "" {
		return "", fmt.Errorf("minted nft not found")
	}
	return id, nil
}

func (minter *NFTokenMinter) SellNFT(conf NFTokenMintConfig, destination types.Address) (types.NFTokenID, error) {
	nft, err := minter.GetOrMintNFT(conf)
	if err != nil {
		return "", err
	}
	if nft == "" {
		return "", fmt.Errorf("failed to retrieve valid nft")
	}
	offerTx := transactions.NFTokenCreateOffer{
		BaseTx: transactions.BaseTx{
			Account:       minter.account,
			Flags:         types.SetFlag(types.FtfSellNFToken),
			SigningPubKey: minter.pubKey,
		},
		Amount:      conf.Price,
		NFTokenID:   nft,
		Destination: destination,
	}
	if conf.Timeout > 0 {
		offerTx.Expiration = common.CurrentRippleTime() + conf.Timeout
	}
	minter.client.AutofillTx(minter.account, &offerTx)
	blob, _ := binarycodec.EncodeForSigning(&offerTx)
	sig, _ := keypairs.Sign(blob, minter.privKey)
	offerTx.BaseTx.TxnSignature = sig
	tx, _ := binarycodec.Encode(&offerTx)
	submitReq := txCl.SubmitRequest{
		TxBlob: tx,
	}
	_, _, err = minter.client.Transaction.Submit(&submitReq)
	if err != nil {
		return "", fmt.Errorf("submit sell offer: %w", err)
	}
	return nft, nil
}

func (minter *NFTokenMinter) BurnNFT(id types.NFTokenID) error {

	// TODO proper CLIO support in xrpl-go
	/*
		infoReq := &clio.NFTInfoRequest{
			NFTokenID: id,
		}

		resp, _, err := minter.client.Clio.NFTInfo(infoReq)
		if err != nil {
			return fmt.Errorf("fetch nft info: %w", err)
		}
		if resp.IsBurned {
			return nil
		}
	*/
	burnTx := &transactions.NFTokenBurn{
		BaseTx: transactions.BaseTx{
			Account:         minter.account,
			TransactionType: transactions.NFTokenBurnTx,
			SigningPubKey:   minter.pubKey,
		},
		NFTokenID: id,
		// TODO Fetch owner from CLIO endpoint if possible
		//	Owner:     resp.Owner,
	}
	minter.client.AutofillTx(minter.account, burnTx)
	blob, _ := binarycodec.EncodeForSigning(burnTx)
	sig, _ := keypairs.Sign(blob, minter.privKey)
	burnTx.BaseTx.TxnSignature = sig
	s, _ := json.MarshalIndent(burnTx, "", "\t")
	fmt.Printf("Tx %s\n", string(s))
	tx, _ := binarycodec.Encode(burnTx)
	submitReq := txCl.SubmitRequest{
		TxBlob: tx,
	}
	_, _, err := minter.client.Transaction.Submit(&submitReq)
	if err != nil {
		return fmt.Errorf("submit burn transaction: %w", err)
	}
	return nil
}
