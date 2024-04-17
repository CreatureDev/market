package nft

import (
	"fmt"
	"strconv"
	"time"

	binarycodec "github.com/CreatureDev/xrpl-go/binary-codec"
	"github.com/CreatureDev/xrpl-go/client"
	"github.com/CreatureDev/xrpl-go/keypairs"
	txCl "github.com/CreatureDev/xrpl-go/model/client/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
	"go.uber.org/atomic"
	"golang.org/x/exp/slices"
)

type NFTokenBroker struct {
	client  *client.XRPLClient
	account types.Address
	privKey string
	pubKey  string
	fee     uint
}

func (b *NFTokenBroker) BrokeredPurchaseOrder(id types.NFTokenID, target types.Address) {

}

func (b *NFTokenBroker) DirectPurchaseOrder(id types.NFTokenID, target types.Address) {

}

type NFTokenBrokerRunner struct {
	broker       *NFTokenBroker
	frequency    time.Duration
	tokenChannel chan types.NFTokenID
	cancel       chan bool
	isRunning    *atomic.Bool
}

func (b *NFTokenBrokerRunner) BrokerNFT(id types.NFTokenID) {
	b.tokenChannel <- id
}

func NewNFTokenBrokerRunner(broker *NFTokenBroker, frequency time.Duration) *NFTokenBrokerRunner {
	return &NFTokenBrokerRunner{
		broker:       broker,
		frequency:    frequency,
		tokenChannel: make(chan types.NFTokenID),
		cancel:       make(chan bool, 1),
		isRunning:    atomic.NewBool(false),
	}
}

func (b *NFTokenBrokerRunner) Cancel() {
	b.cancel <- true
}

func (b *NFTokenBrokerRunner) brokerNFT(id types.NFTokenID) bool {
	sales, err := SellOffersForNFT(b.broker.client, id)
	if err != nil || len(sales) == 0 {
		// return TRUE to not retry NFT in case of error/missing sales
		return true
	}
	buys, err := BuyOffersForNFT(b.broker.client, id)
	if err != nil || len(buys) == 0 {
		// return FALSE to retry if no buy offers exist
		return false
	}
	isXrp := sales[0].Amount.Kind() == types.XRP
	for _, buy := range buys {
		var fee types.CurrencyAmount
		if buy.Amount.Kind() != sales[0].Amount.Kind() {
			continue
		}
		if !isXrp {
			sellCur := sales[0].Amount.(types.IssuedCurrencyAmount)
			buyCur := buy.Amount.(types.IssuedCurrencyAmount)
			if sellCur.Issuer != buyCur.Issuer || sellCur.Currency != buyCur.Currency {
				continue
			}
			buyVal, err := strconv.ParseFloat(buyCur.Value, 64)
			if err != nil {
				continue
			}
			sellVal, err := strconv.ParseFloat(sellCur.Value, 64)
			if err != nil {
				continue
			}
			if buyVal <= sellVal {
				continue
			}
			feeVal := strconv.FormatFloat(buyVal-sellVal, 'g', 6, 64)
			fee = types.IssuedCurrencyAmount{
				Issuer:   sellCur.Issuer,
				Currency: sellCur.Currency,
				Value:    feeVal,
			}
		} else {
			sellVal := sales[0].Amount.(types.XRPCurrencyAmount)
			buyVal := buy.Amount.(types.XRPCurrencyAmount)
			if buyVal <= sellVal {
				continue
			}
			fee = buyVal - sellVal
		}
		offerTx := transactions.NFTokenAcceptOffer{
			BaseTx: transactions.BaseTx{
				Account:       b.broker.account,
				SigningPubKey: b.broker.pubKey,
			},
			NFTokenBrokerFee: fee,
			NFTokenSellOffer: types.Hash256(sales[0].NFTokenOfferIndex),
			NFTokenBuyOffer:  types.Hash256(buy.NFTokenOfferIndex),
		}
		b.broker.client.AutofillTx(b.broker.account, &offerTx)
		blob, _ := binarycodec.EncodeForSigning(&offerTx)
		sig, _ := keypairs.Sign(blob, b.broker.privKey)
		offerTx.BaseTx.TxnSignature = sig
		tx, _ := binarycodec.Encode(&offerTx)
		submitReq := txCl.SubmitRequest{
			TxBlob: tx,
		}
		_, _, err = b.broker.client.Transaction.Submit(&submitReq)
		if err != nil {
			continue
		}
		return true
	}
	return false
}

func (b *NFTokenBrokerRunner) brokerNFTs() {
	var nfts []types.NFTokenID
R:
	for {
		select {
		case id := <-b.tokenChannel:
			if slices.Contains(nfts, id) {
				continue R
			}
			nfts = append(nfts, id)
		default:
			break R
		}
	}
	for _, id := range nfts {
		if !b.brokerNFT(id) {
			b.tokenChannel <- id
		}
	}
}

func (b *NFTokenBrokerRunner) run() {
	timer := time.NewTimer(b.frequency)
	cancel := false
	for !cancel {
		select {
		case <-b.cancel:
			cancel = true
		case <-timer.C:
			b.brokerNFTs()
			timer.Reset(b.frequency)
		}
	}

	b.isRunning.Store(false)
}

func (b *NFTokenBrokerRunner) Run() error {
	if b.tokenChannel == nil ||
		b.cancel == nil ||
		b.broker == nil ||
		b.frequency == 0 ||
		b.isRunning == nil {
		return fmt.Errorf("token broker was initialized incorrectly")
	}
	if !b.isRunning.CAS(false, true) {
		return fmt.Errorf("provided broker is already running")
	}
	go b.run()
	return nil
}
