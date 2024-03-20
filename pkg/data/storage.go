package data

import (
	"github.com/CreatureDev/market/pkg/product"
	"github.com/CreatureDev/market/pkg/user"
)

type Storage interface {
	AddProduct(product.Product)
	GetProducts() []product.Product
	GetPublisherProducts(user.Publisher) []product.Product
	GetPurchases(user.Consumer) []product.PurchaseRecord
}
