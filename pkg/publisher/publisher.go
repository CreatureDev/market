package publisher

import (
	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/market/pkg/product"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
)

type Publisher struct {
	Account  types.Address
	Client   marketv1.PublisherServiceClient
	Products []*product.Product
}
