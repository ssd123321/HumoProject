package grpc

import (
	"Tasks/Service"
	"Tasks/handlers/grpc/gprc_api"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

type MyGrpcBankingServer struct {
	*gprc_api.UnimplementedGrpcBankingServer
	service *Service.Service
}

func NewGrpcBaningServer(server *gprc_api.UnimplementedGrpcBankingServer, service *Service.Service) *MyGrpcBankingServer {
	return &MyGrpcBankingServer{
		server,
		service,
	}
}
func (m *MyGrpcBankingServer) GetPersonByID(ctx context.Context, income *gprc_api.IncomePersonID) (*gprc_api.ReplyPerson, error) {
	ctx = context.WithValue(ctx, "id", int(income.Id))
	p, err := m.service.GetPersonByID(ctx)
	if err != nil {
		log.Printf("grpc: GetPersonByID - %v", err)
		return nil, grpc.Errorf(codes.NotFound, "Error: %v", err.Error())
	}
	return &gprc_api.ReplyPerson{
		Id:        int64(p.ID),
		Content:   &gprc_api.Content{Dimension: &gprc_api.Dimenstion{Height: float32(p.Content.Dimension.Height), Weight: float32(p.Content.Dimension.Weight)}, Name: p.Content.Name, Age: int64(p.Content.Age)},
		CreatedAt: timestamppb.New(*p.CreatedAt),
		UpdatedAt: timestamppb.New(*p.UpdatedAt),
	}, nil
}
func NewGrpcServer(port string) (net.Listener, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Error - NewGrpcServer: %v", err)
		return nil, err
	}
	return lis, nil
}
