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
	Logger         *zap.Logger
	Mongo          *dao.Mongo
	ProfileManager ProfileManager
	CarManager     CarManager
	PoiManager     PoiManager
}

func (s *Service) GetTrip(ctx context.Context, request *GetTripRequest) (*GetTripResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetTrips(ctx context.Context, request *GetTripsRequest) (*GetTripsResponses, error) {
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

type ProfileManager interface {
	Verify(ctx context.Context, userId string) (string, error)
}

type CarManager interface {
	Verify(ctx context.Context, carId string, location *Location) error
	Unlock(ctx context.Context, carId string) error
}

type PoiManager interface {
	GetPoiName(l *Location) (string, error)
}

func (s *Service) CreateTrip(ctx context.Context, request *CreateTripRequest) (*TripEntity, error) {
	userId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("hello world", zap.String("userId", userId))
	identityId, verifyErr := s.ProfileManager.Verify(ctx, userId)
	if verifyErr != nil {
		return nil, status.Error(codes.FailedPrecondition, verifyErr.Error())
	}
	carVerifyErr := s.CarManager.Verify(ctx, request.CarId, request.Start)
	if carVerifyErr != nil {
		return nil, status.Error(codes.FailedPrecondition, carVerifyErr.Error())
	}

	poiName, getPoiNameErr := s.PoiManager.GetPoiName(request.Start)
	if getPoiNameErr != nil {
		s.Logger.Info("没有找到行程locationDesc", zap.Stringer("location", request.Start))
	}

	createTrip, createTripErr := s.Mongo.CreateTrip(ctx, &Trip{
		UserId: userId,
		CarId:  request.CarId,
		Start: &LocationStatus{
			Location:     request.Start,
			FeeCent:      0,
			KmDriven:     0,
			LocationDesc: poiName,
		},
		Current: &LocationStatus{
			Location:     request.Start,
			FeeCent:      0,
			KmDriven:     0,
			LocationDesc: poiName,
		},
		End:        nil,
		Status:     TripStatus_IN_PROGRESS,
		IdentityId: identityId,
	})
	if createTripErr != nil {
		return nil, status.Error(codes.AlreadyExists, createTripErr.Error())
	}
	go func(c context.Context, carId string) {
		err := s.CarManager.Unlock(c, carId)
		if err != nil {
			s.Logger.Error("车辆开锁失败", zap.Error(err))
		}
	}(context.Background(), createTrip.Trip.CarId)
	return &TripEntity{
		Id:   createTrip.ID.Hex(),
		Trip: createTrip.Trip,
	}, nil
}

func (s *Service) mustEmbedUnimplementedTripServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s Service) calcCurrentStatus(trip *Trip, cur *LocationStatus) *LocationStatus {
	return nil
}
