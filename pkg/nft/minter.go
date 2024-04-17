package nft

import (
	"fmt"

	binarycodec "github.com/CreatureDev/xrpl-go/binary-codec"
	"github.com/CreatureDev/xrpl-go/client"
	"github.com/CreatureDev/xrpl-go/keypairs"
	"github.com/CreatureDev/xrpl-go/model/client/account"
	"github.com/CreatureDev/xrpl-go/model/client/common"
	txCl "github.com/CreatureDev/xrpl-go/model/client/transactions"
	"github.com/CreatureDev/xrpl-go/model/ledger"
	"github.com/CreatureDev/xrpl-go/model/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
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
	offers, _ := SellOffersForNFT(minter.client, id)
	return len(offers) == 0
}

func (minter *NFTokenMinter) getAccountObjects(addr types.Address, kind account.AccountObjectType) ([]ledger.LedgerObject, error) {
	objsReq := &account.AccountObjectsRequest{
		Account: addr,
		Type:    kind,
	}
	objsRes, _, err := minter.client.Account.AccountObjects(objsReq)
	if err != nil {
		return nil, err
	}
	return objsRes.AccountObjects, nil
}

func (minter *NFTokenMinter) CancelExpiredNFTSales() error {
	objs, err := minter.getAccountObjects(minter.account, account.NFTOfferObject)
	if err != nil {
		return err
	}

	curTime := common.CurrentRippleTime()
	var expired []types.Hash256
	for _, obj := range objs {
		offer, ok := obj.(*ledger.NFTokenOffer)
		if !ok {
			continue
		}
		if offer.Expiration != 0 && offer.Expiration < curTime {
			expired = append(expired, types.Hash256(offer.Index))
		}
	}
	if len(expired) == 0 {
		return nil
	}

	cancelOfferTx := transactions.NFTokenCancelOffer{
		BaseTx: transactions.BaseTx{
			Account:       minter.account,
			SigningPubKey: minter.pubKey,
		},
		NFTokenOffers: expired,
	}
	// TODO resubmit on failed TX attempts
	// make submit helper function (autofills/signs until success or too many failures)
	minter.client.AutofillTx(minter.account, &cancelOfferTx)
	blob, err := binarycodec.EncodeForSigning(&cancelOfferTx)
	if err != nil {
		return fmt.Errorf("encoding cancel offer tx for signing: %w", err)
	}
	sig, _ := keypairs.Sign(blob, minter.privKey)
	cancelOfferTx.BaseTx.TxnSignature = sig
	fmt.Printf("TX: %+v\n\n", cancelOfferTx)
	tx, err := binarycodec.Encode(&cancelOfferTx)
	if err != nil {
		return fmt.Errorf("encoding cancel offer tx for submitting: %w", err)
	}
	submitReq := txCl.SubmitRequest{
		TxBlob: tx,
	}
	_, _, err = minter.client.Transaction.Submit(&submitReq)
	if err != nil {
		return fmt.Errorf("submiting cancel offer tx: %w", err)
	}
	return nil
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
