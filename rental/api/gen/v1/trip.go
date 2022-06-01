package rentalpb

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger *zap.Logger
}

func (s *Service) CreateTrip(ctx context.Context, request *CreateTripRequest) (*CreateTripResponse, error) {
	s.Logger.Info("hello world")
	return nil, status.Error(codes.Unimplemented, "还没有实现")
}

func (s *Service) mustEmbedUnimplementedTripServiceServer() {
	//TODO implement me
	panic("implement me")
}
