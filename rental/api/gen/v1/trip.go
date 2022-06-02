package rentalpb

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"server2/rental/trip/dao"
	"server2/shared/auth"
)

type Service struct {
	Logger *zap.Logger
	Mongo  *dao.Mongo
}

func (s *Service) GetTrip(ctx context.Context, request *GetTripRequest) (*GetTripResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetTrips(ctx context.Context, request *GetTripsRequest) (*GetTripResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UpdateTrip(ctx context.Context, request *UpdateTripRequest) (*UpdateTripResponse, error) {
	userId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "身份验证未通过")
	}
	trip, getTripErr := s.Mongo.GetTrip(ctx, request.Id, userId)
	if getTripErr != nil {
		return nil, status.Error(codes.Internal, "")
	}
	if request.Current != nil {
		trip.Trip.Current = s.calcCurrentStatus(trip.Trip, request.Current)
	}
	if request.EndTrip {
		trip.Trip.End = trip.Trip.Current
		trip.Trip.Status = TripStatus_FINISHED
	}
	updateTripErr := s.Mongo.UpdateTrip(ctx, trip.ID.Hex(), userId, trip.UpdateAt, trip.Trip)
	if updateTripErr != nil {
		return nil, status.Error(codes.Internal, updateTripErr.Error())
	}
	return nil, nil
}

func (s *Service) DelTrip(ctx context.Context, request *DelTripRequest) (*DelTripResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CreateTrip(ctx context.Context, request *CreateTripRequest) (*CreateTripResponse, error) {
	userId, err := auth.UserIdFromContext(ctx)
	s.Logger.Info("hello world", zap.String("userId", userId))
	if err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unimplemented, "还没有实现")
}

func (s *Service) mustEmbedUnimplementedTripServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s Service) calcCurrentStatus(trip *Trip, cur *LocationStatus) *LocationStatus {
	return nil
}
