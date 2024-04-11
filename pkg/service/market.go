package service

import (
	"context"

	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/market/pkg/config"
	"github.com/CreatureDev/market/pkg/product"
	errors "golang.org/x/xerrors"
	"google.golang.org/grpc"
)

type MarketService struct {
	marketv1.UnimplementedMarketServiceServer
}

func NewMarketService(_ config.Config) *MarketService {
	return &MarketService{}
}

func (s *MarketService) PurchaseProduct(ctx context.Context, req *marketv1.PurchaseProductRequest) (*marketv1.PurchaseProductResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *MarketService) GetPoduct(ctx context.Context, req *marketv1.GetProductRequest) (*marketv1.GetProductResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *MarketService) ListProducts(ctx context.Context, req *marketv1.ListProductsRequest) (*marketv1.ListProductsResponse, error) {
	return nil, errors.Errorf("unimplemented")
}

func (s *MarketService) RegisterPublisher(ctx context.Context, req *marketv1.RegisterPublisherRequest) (*marketv1.RegisterPublisherResponse, error) {
	conn, err := grpc.Dial(req.GetConnection(), grpc.EmptyDialOption{})
	if err != nil {
		return nil, err
	}
	publisher := marketv1.NewPublisherServiceClient(conn)
	res, err := publisher.ListPublishedProducts(ctx, &marketv1.ListPublishedProductsRequest{}, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}
	prodsProto := res.GetProducts()
	var prods []*product.Product
	for _, p := range prodsProto {
		prod, err := product.ProductFromProto(p)
		if err != nil {
			return nil, err
		}
		prods = append(prods, prod)
	}

	market.Add

	return &marketv1.RegisterPublisherResponse{}, nil
}
