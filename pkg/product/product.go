package product

import (
	"fmt"
	"strconv"

	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
	"github.com/google/uuid"
)

type Product struct {
	id            uuid.UUID
	nftTaxon      uint32
	file          string
	publisher     types.Address
	standardPrice types.CurrencyAmount
	currentPrice  types.CurrencyAmount
}

func basePriceFromProto(proto *marketv1.Product) (types.CurrencyAmount, error) {
	switch proto.GetDenomination() {
	case marketv1.Product_DROPS:
		xrp, err := strconv.ParseUint(proto.GetStandardPrice(), 10, 64)
		if err != nil {
			return nil, err
		}
		return types.XRPCurrencyAmount(xrp), nil

	case marketv1.Product_XRP:
		xrp, err := strconv.ParseFloat(proto.GetStandardPrice(), 32)
		if err != nil {
			return nil, err
		}
		return types.XRPDropsFromFloat(float32(xrp)), nil
	case marketv1.Product_ISSUED:
		return types.IssuedCurrencyAmount{
			Issuer:   types.Address(proto.GetIssuer()),
			Currency: proto.GetCurrency(),
			Value:    proto.GetStandardPrice(),
		}, nil
	}
	return nil, fmt.Errorf("unknown currency format")
}

func (p *Product) Publisher() types.Address {
	return p.publisher
}

func (p *Product) ProductUUID() uuid.UUID {
	return p.id
}

func ProductFromProto(proto *marketv1.Product) (*Product, error) {
	uid, err := uuid.Parse(proto.GetUid())
	if err != nil {
		return nil, err
	}
	ret := &Product{
		id:        uid,
		publisher: types.Address(proto.GetPublisher()),
		nftTaxon:  proto.GetTaxon(),
	}
	return ret
}
