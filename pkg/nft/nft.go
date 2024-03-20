package nft

import (
	"fmt"

	binarycodec "github.com/CreatureDev/xrpl-go/binary-codec"
	"github.com/CreatureDev/xrpl-go/client"
	"github.com/CreatureDev/xrpl-go/model/client/account"
	txCl "github.com/CreatureDev/xrpl-go/model/client/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
	"github.com/xyield/xrpl-go/keypairs"
)

type NFTokenMintConfig struct {
	NFTokenTaxon uint
	TransferFee  uint16
	URI          types.NFTokenURI
	Flags        *types.Flag
}

type NFTokenMinter struct {
	client  *client.XRPLClient
	account types.Address
	privKey string
	pubKey  string
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
		// TODO ensure NFT is ready to be sold (no pending Tx)
		if nft.NFTokenTaxon == conf.NFTokenTaxon {
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
		URI:          types.NFTokenURI(conf.URI),
	}
	if err := cl.AutofillTx(minter.account, &mintTx); err != nil {
		return "", fmt.Errorf("autofill mint tx: %w", err)
	}
	bin, err := binarycodec.EncodeForSigning(&mintTx)
	if err != nil {
		return "", fmt.Errorf("encoding mint tx for signing: %w", err)
	}
	sig, _ := keypairs.Sign(string(bin), minter.privKey)
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
