package shopify

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/Shridhar2104/logilo/shopify/pb"
)
type grpcServer struct {
	pb.UnimplementedShopifyServiceServer
	service Service
}

func NewGRPCServer(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterShopifyServiceServer(server, &grpcServer{
		UnimplementedShopifyServiceServer: pb.UnimplementedShopifyServiceServer{}, // Add this
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) SyncOrders(ctx context.Context, r *pb.SyncOrdersRequest) (*pb.SyncOrdersResponse, error) {

	err := s.service.SyncOrders(ctx, r.ShopName, r.SinceId, int(r.Limit), r.Token)
	if err != nil {
		return nil, err
	}
	return &pb.SyncOrdersResponse{}, nil


}
func (s *grpcServer) GetOrdersForShopAndAccount(ctx context.Context, r *pb.GetOrdersForShopAndAccountRequest) (*pb.GetOrdersForShopAndAccountResponse, error){
	orders, err := s.service.GetOrdersForShopAndAccount(ctx, r.ShopName, r.AccountId)
	if err != nil {
		return nil, err
	}
	ordersPb := make([]*pb.Order, len(orders))
	for i, order := range orders {
		ordersPb[i] = &pb.Order{
			Id: order.ID,
			AccountId: order.AccountId,
			ShopId: order.ShopName,
			TotalPrice: float32(order.TotalPrice),
			OrderId: order.OrderId,
		}
	}
	return &pb.GetOrdersForShopAndAccountResponse{
		Orders: ordersPb,
	}, nil

}
func (s *grpcServer) UpdateOrder(ctx context.Context, r *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error){

	order := &Order{
		ID: r.Order.Id,
		ShopName: r.ShopName,
		AccountId: r.AccountId,
		TotalPrice: float64(r.Order.TotalPrice),
	}
	err := s.service.UpdateOrder(ctx, *order, r.AccountId, r.ShopName)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{}, nil

}
func (s *grpcServer) StoreToken(ctx context.Context, r *pb.StoreTokenRequest) (*pb.StoreTokenResponse, error){
	err := s.service.StoreToken(ctx, r.ShopName, r.AccountId, r.Token)
	if err != nil {
		return nil, err
	}
	return &pb.StoreTokenResponse{}, nil
}	
