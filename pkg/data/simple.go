package data

import (
	"fmt"

	"github.com/CreatureDev/market/pkg/product"
	"github.com/CreatureDev/market/pkg/user"
	"github.com/google/uuid"
)

type SimpleStorage struct {
	products  []product.Product
	ownership map[uuid.UUID
	][]product.PurchaseRecord
}

func (s *SimpleStorage) AddProduct(prod product.Product) error {
	if _, e := s.ownership[prod.ProductUUID()]; e {
		return fmt.Errorf("Product %v already exists", prod)
	}

	s.products = append(s.products, prod)
	rec := make([]product.PurchaseRecord, 0)
	s.ownership[prod] = rec

	return nil
}

func (s *SimpleStorage) GetProducts() []product.Product {
	return s.products
}

func (s *SimpleStorage) GetPublisherProducts(pub user.Publisher) []product.Product {
	ret := make([]product.Product, 0)
	for _, v := range s.products {
		if v.Publisher() == pub.Uid() {
			ret = append(ret, v)
		}
	}
	return ret
}

func (s *SimpleStorage) GetPurchases(u user.Consumer) []product.PurchaseRecord {
	ret := make([]product.PurchaseRecord, 0)
	for _, prs := range s.ownership {
		for _, pr := range prs {
			if pr.User() == u.Uid() {
				ret = append(ret, pr)
			}
		}
	}
	return ret
}
