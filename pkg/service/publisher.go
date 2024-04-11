package service

import (
	"context"

	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/market/pkg/config"
	errors "golang.org/x/xerrors"
)

type PublisherService struct {
	marketv1.UnimplementedPublisherServiceServer
}

func NewPublisherService(_ config.PublisherConfig) *PublisherService {
	return &PublisherService{}
}

func (s *PublisherService) ListPublishedProducts(ctx context.Context, req *marketv1.ListPublishedProductsRequest) (*marketv1.ListPublishedProductsResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *PublisherService) GetPublishedProduct(ctx context.Context, req *marketv1.GetPublishedProductRequest) (*marketv1.GetPublishedProductResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *PublisherService) CreatePurchaseOrder(ctx context.Context, req *marketv1.CreatePurchaseOrderRequest) (*marketv1.CreatePurchaseOrderResponse, error) {
	return nil, errors.Errorf("unimplemented")
}
