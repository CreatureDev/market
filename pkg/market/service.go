package market

import (
	"context"

	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/market/pkg/config"
	errors "golang.org/x/xerrors"
)

type Service struct {
	marketv1.UnimplementedMarketServiceServer
}

func NewService(_ config.Config) *Service {
	return &Service{}
}

func (s *Service) PostProduct(ctx context.Context, req *marketv1.PostProductRequest) (*marketv1.PostProductResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *Service) GetPoduct(ctx context.Context, req *marketv1.GetProductRequest) (*marketv1.GetProductResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *Service) ListProducts(ctx context.Context, req *marketv1.ListProductsRequest) (*marketv1.ListProductsResponse, error) {
	return nil, errors.Errorf("unimplemented")
}
