package nft_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CreatureDev/market/pkg/nft"
	binarycodec "github.com/CreatureDev/xrpl-go/binary-codec"
	"github.com/CreatureDev/xrpl-go/client"
	jsonrpcclient "github.com/CreatureDev/xrpl-go/client/jsonrpc"
	"github.com/CreatureDev/xrpl-go/keypairs"
	"github.com/CreatureDev/xrpl-go/model/client/account"
	"github.com/CreatureDev/xrpl-go/model/client/faucet"
	txCl "github.com/CreatureDev/xrpl-go/model/client/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
	"github.com/stretchr/testify/assert"
)

var cl *client.XRPLClient

type testAccount struct {
	Account types.Address
	PubKey  string
	PrivKey string
}

var testMinterAcc testAccount = testAccount{
	Account: "rLJhgFU2caBZjYWi3EnbrmZbZVLuyvT5sC",
	PrivKey: "EDE5443AF66200C39D555CCA6B07774749C55F5A19F5F91C782EC018F3715B90E6",
	PubKey:  "EDB3DB8CF9E7392D6848E9923F52D3BCB15750683AE5B2C824EFC2DD963D912C4C",
}

var testBuyerAcc testAccount = testAccount{
	Account: "rHXanFQnJgN4QSKoFp6WRz69a7F66XJdDi",
	PrivKey: "ED48000EFB8477958150F219E48167A5B12CCC2D50163178C006635690AF5CA1D0",
	PubKey:  "ED8B24E0F960835AE706EF4FEA0EEC755021FACDB49950AEB82AC8EDA0A593D78E",
}

func resetNFTs(t *testing.T, minter *nft.NFTokenMinter, conf nft.NFTokenMintConfig) {
	for nftId := minter.GetValidNFT(conf); nftId != ""; nftId = minter.GetValidNFT(conf) {
		err := minter.BurnNFT(nftId)
		assert.Nil(t, err)
		if err != nil {
			break
		}
	}
}

func TestMintToken(t *testing.T) {
	minter := testMinter(testMinterAcc)
	config := nft.NFTokenMintConfig{
		NFTokenTaxon: 1,
		URIGenerator: nft.ConstantURIGenerator(types.NFTokenURI("")),
		Flags:        types.NewFlag().SetFlag(types.FtfBurnable),
	}
	resetNFTs(t, minter, config)
	preMint := minter.GetValidNFT(config)
	assert.Empty(t, preMint)
	id, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
	err = minter.BurnNFT(id)
	assert.Nil(t, err)
}

func TestSellToken(t *testing.T) {
	minter := testMinter(testMinterAcc)
	target := testBuyerAcc.Account
	assert.NotEqual(t, minter.GetAccount(), target)
	config := nft.NFTokenMintConfig{
		NFTokenTaxon: 2,
		URIGenerator: nft.ConstantURIGenerator(types.NFTokenURI("00")),
		Flags:        types.NewFlag().SetFlag(types.FtfBurnable),
		Price:        types.XRPCurrencyAmount(50),
	}
	resetNFTs(t, minter, config)
	preMint := minter.GetValidNFT(config)
	assert.Empty(t, preMint)
	id, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
	id2, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
	hash, err := minter.SellNFT(config, target)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
	sales, err := minter.SellOffersForNFT(hash)
	assert.Nil(t, err)
	assert.NotEmpty(t, sales)
	assert.Equal(t, 1, len(sales))
	assert.Equal(t, types.XRPCurrencyAmount(50), sales[0].Amount)
	err = minter.BurnNFT(id)
	assert.Nil(t, err)
}

func TestBuyToken(t *testing.T) {
	minter := testMinter(testMinterAcc)
	target := testBuyerAcc.Account
	assert.NotEqual(t, minter.GetAccount(), target)
	config := nft.NFTokenMintConfig{
		NFTokenTaxon: 3,
		URIGenerator: nft.ConstantURIGenerator(types.NFTokenURI("00")),
		Flags:        types.NewFlag().SetFlag(types.FtfBurnable),
		Price:        types.XRPCurrencyAmount(50),
	}
	resetNFTs(t, minter, config)
	preMint := minter.GetValidNFT(config)
	assert.Empty(t, preMint)
	id, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
	id2, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
	hash, err := minter.SellNFT(config, target)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
	sales, err := minter.SellOffersForNFT(hash)
	assert.Nil(t, err)
	assert.NotEmpty(t, sales)
	assert.Equal(t, 1, len(sales))
	assert.Equal(t, types.XRPCurrencyAmount(50), sales[0].Amount)
	buyTx := &transactions.NFTokenAcceptOffer{
		BaseTx: transactions.BaseTx{
			Account:       testBuyerAcc.Account,
			SigningPubKey: testBuyerAcc.PubKey,
		},
		NFTokenSellOffer: types.Hash256(sales[0].NFTokenOfferIndex),
	}
	err = cl.AutofillTx(testBuyerAcc.Account, buyTx)
	assert.Nil(t, err)
	blob, _ := binarycodec.EncodeForSigning(buyTx)
	sig, _ := keypairs.Sign(blob, testBuyerAcc.PrivKey)
	buyTx.BaseTx.TxnSignature = sig
	tx, _ := binarycodec.Encode(buyTx)
	submitReq := txCl.SubmitRequest{
		TxBlob: tx,
	}
	sub, _, err := cl.Transaction.Submit(&submitReq)
	fmt.Printf("%+v\n", sub)
	assert.Nil(t, err)
	buyerMinter := testMinter(testBuyerAcc)
	err = buyerMinter.BurnNFT(id)
	assert.Nil(t, err)
}

func TestExpireSale(t *testing.T) {
	minter := testMinter(testMinterAcc)
	target := testBuyerAcc.Account
	assert.NotEqual(t, minter.GetAccount(), target)
	config := nft.NFTokenMintConfig{
		NFTokenTaxon: 4,
		URIGenerator: nft.ConstantURIGenerator(types.NFTokenURI("00")),
		Flags:        types.NewFlag().SetFlag(types.FtfBurnable),
		Price:        types.XRPCurrencyAmount(50),
		Timeout:      5,
	}
	resetNFTs(t, minter, config)
	preMint := minter.GetValidNFT(config)
	assert.Empty(t, preMint)
	id, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
	hash, err := minter.SellNFT(config, target)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
	time.Sleep(5 * time.Second)
	err = minter.CancelExpiredNFTSales()
	assert.Nil(t, err)
	sales, err := minter.SellOffersForNFT(hash)
	assert.Nil(t, err)
	assert.Empty(t, sales)
	id = minter.GetValidNFT(config)
	assert.NotEmpty(t, id)
	assert.Nil(t, err)
	err = minter.BurnNFT(hash)
	assert.Nil(t, err)
}

func TestGetValidNFT(t *testing.T) {
	minter := testMinter(testMinterAcc)
	target := testBuyerAcc.Account
	assert.NotEqual(t, minter.GetAccount(), target)
	config := nft.NFTokenMintConfig{
		NFTokenTaxon: 5,
		URIGenerator: nft.ConstantURIGenerator(types.NFTokenURI("00")),
		Flags:        types.NewFlag().SetFlag(types.FtfBurnable),
		Price:        types.XRPCurrencyAmount(50),
	}
	resetNFTs(t, minter, config)
	preMint := minter.GetValidNFT(config)
	assert.Empty(t, preMint)
	id, err := minter.GetOrMintNFT(config)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
	id2 := minter.GetValidNFT(config)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
	err = minter.BurnNFT(id)
	assert.Nil(t, err)
}

func testMinter(acc testAccount) *nft.NFTokenMinter {
	return nft.NewNFTokenMinter(cl, acc.Account, acc.PrivKey, acc.PubKey)
}

func init() {
	conf, err := client.NewJsonRpcConfig("https://s.altnet.rippletest.net:51234/")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	rpccl := jsonrpcclient.NewJsonRpcClient(conf)
	cl = client.NewXRPLClient(rpccl)
	initAccount(testMinterAcc)
	initAccount(testBuyerAcc)
}

func initAccount(acc testAccount) {
	accReq := &account.AccountInfoRequest{
		Account: acc.Account,
	}
	_, _, err := cl.Account.AccountInfo(accReq)
	if err != nil && strings.Contains(err.Error(), "NotFound") {
		_, _, err := cl.Faucet.FundAccount(&faucet.FundAccountRequest{Destination: acc.Account})
		if err != nil {
			panic("Failed to fund account")
		}
	}
}
