// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package marketv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PublisherServiceClient is the client API for PublisherService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublisherServiceClient interface {
	ListPublishedProducts(ctx context.Context, in *ListPublishedProductsRequest, opts ...grpc.CallOption) (*ListPublishedProductsResponse, error)
	GetPublishedProduct(ctx context.Context, in *GetPublishedProductRequest, opts ...grpc.CallOption) (*GetPublishedProductResponse, error)
	CreatePurchaseOrder(ctx context.Context, in *CreatePurchaseOrderRequest, opts ...grpc.CallOption) (*CreatePurchaseOrderResponse, error)
}

type publisherServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPublisherServiceClient(cc grpc.ClientConnInterface) PublisherServiceClient {
	return &publisherServiceClient{cc}
}

func (c *publisherServiceClient) ListPublishedProducts(ctx context.Context, in *ListPublishedProductsRequest, opts ...grpc.CallOption) (*ListPublishedProductsResponse, error) {
	out := new(ListPublishedProductsResponse)
	err := c.cc.Invoke(ctx, "/market.v1.PublisherService/ListPublishedProducts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherServiceClient) GetPublishedProduct(ctx context.Context, in *GetPublishedProductRequest, opts ...grpc.CallOption) (*GetPublishedProductResponse, error) {
	out := new(GetPublishedProductResponse)
	err := c.cc.Invoke(ctx, "/market.v1.PublisherService/GetPublishedProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherServiceClient) CreatePurchaseOrder(ctx context.Context, in *CreatePurchaseOrderRequest, opts ...grpc.CallOption) (*CreatePurchaseOrderResponse, error) {
	out := new(CreatePurchaseOrderResponse)
	err := c.cc.Invoke(ctx, "/market.v1.PublisherService/CreatePurchaseOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublisherServiceServer is the server API for PublisherService service.
// All implementations must embed UnimplementedPublisherServiceServer
// for forward compatibility
type PublisherServiceServer interface {
	ListPublishedProducts(context.Context, *ListPublishedProductsRequest) (*ListPublishedProductsResponse, error)
	GetPublishedProduct(context.Context, *GetPublishedProductRequest) (*GetPublishedProductResponse, error)
	CreatePurchaseOrder(context.Context, *CreatePurchaseOrderRequest) (*CreatePurchaseOrderResponse, error)
	mustEmbedUnimplementedPublisherServiceServer()
}

// UnimplementedPublisherServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPublisherServiceServer struct {
}

func (UnimplementedPublisherServiceServer) ListPublishedProducts(context.Context, *ListPublishedProductsRequest) (*ListPublishedProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPublishedProducts not implemented")
}
func (UnimplementedPublisherServiceServer) GetPublishedProduct(context.Context, *GetPublishedProductRequest) (*GetPublishedProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublishedProduct not implemented")
}
func (UnimplementedPublisherServiceServer) CreatePurchaseOrder(context.Context, *CreatePurchaseOrderRequest) (*CreatePurchaseOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePurchaseOrder not implemented")
}
func (UnimplementedPublisherServiceServer) mustEmbedUnimplementedPublisherServiceServer() {}

// UnsafePublisherServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PublisherServiceServer will
// result in compilation errors.
type UnsafePublisherServiceServer interface {
	mustEmbedUnimplementedPublisherServiceServer()
}

func RegisterPublisherServiceServer(s grpc.ServiceRegistrar, srv PublisherServiceServer) {
	s.RegisterService(&PublisherService_ServiceDesc, srv)
}

func _PublisherService_ListPublishedProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPublishedProductsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServiceServer).ListPublishedProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/market.v1.PublisherService/ListPublishedProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServiceServer).ListPublishedProducts(ctx, req.(*ListPublishedProductsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PublisherService_GetPublishedProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublishedProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServiceServer).GetPublishedProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/market.v1.PublisherService/GetPublishedProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServiceServer).GetPublishedProduct(ctx, req.(*GetPublishedProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PublisherService_CreatePurchaseOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePurchaseOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServiceServer).CreatePurchaseOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/market.v1.PublisherService/CreatePurchaseOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServiceServer).CreatePurchaseOrder(ctx, req.(*CreatePurchaseOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PublisherService_ServiceDesc is the grpc.ServiceDesc for PublisherService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PublisherService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "market.v1.PublisherService",
	HandlerType: (*PublisherServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPublishedProducts",
			Handler:    _PublisherService_ListPublishedProducts_Handler,
		},
		{
			MethodName: "GetPublishedProduct",
			Handler:    _PublisherService_GetPublishedProduct_Handler,
		},
		{
			MethodName: "CreatePurchaseOrder",
			Handler:    _PublisherService_CreatePurchaseOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "market/v1/publisher_service.proto",
}
