package product

import (
	"github.com/google/uuid"
	"github.com/xyield/xrpl-go/model/transactions/types"
)

type Product struct {
	id        uuid.UUID
	nftTaxon  uint32
	file      string
	publisher types.Address
}

type PurchaseRecord struct {
	nftID string
	user  types.Address
}

func (p *Product) Publisher() types.Address {
	return p.publisher
}

func (p *Product) ProductUUID() uuid.UUID {
	return p.id
}

func (p *PurchaseRecord) User() types.Address {
	return p.user
}
