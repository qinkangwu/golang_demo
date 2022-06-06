package service

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"server2/rental/api/gen/v1"
	"server2/rental/trip/dao"
	"server2/shared/auth"
	"time"
)

type Service struct {
	Logger         *zap.Logger
	Mongo          *dao.Mongo
	ProfileManager ProfileManager
	CarManager     CarManager
	PoiManager     PoiManager
	rentalpb.UnimplementedTripServiceServer
}

func (s *Service) GetTrip(ctx context.Context, request *rentalpb.GetTripRequest) (*rentalpb.GetTripResponse, error) {
	userId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "身份验证未通过")
	}
	tripId := request.Id
	getTrip, getTripErr := s.Mongo.GetTrip(ctx, tripId, userId)
	if getTripErr != nil {
		return nil, status.Error(codes.InvalidArgument, "trip id错误")
	}
	return &rentalpb.GetTripResponse{
		Trip: &rentalpb.TripEntity{
			Id:   getTrip.ID.Hex(),
			Trip: getTrip.Trip,
		},
	}, nil
}

func (s *Service) GetTrips(ctx context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponses, error) {
	userId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "身份验证未通过")
	}
	getTrips, err := s.Mongo.GetTrips(ctx, request.Status, userId)
	if err != nil {
		s.Logger.Error("获取不到trips", zap.Error(err))
		return nil, status.Error(codes.Internal, "服务器端错误")
	}

	res := &rentalpb.GetTripsResponses{}

	for _, tr := range getTrips {
		res.Trips = append(res.Trips, &rentalpb.TripEntity{
			Id:   tr.ID.Hex(),
			Trip: tr.Trip,
		})
	}
	return res, err
}

func (s *Service) UpdateTrip(ctx context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.UpdateTripResponse, error) {
	userId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "身份验证未通过")
	}
	trip, getTripErr := s.Mongo.GetTrip(ctx, request.Id, userId)
	if getTripErr != nil {
		return nil, status.Error(codes.Internal, "")
	}
	if trip.Trip.Current == nil {
		s.Logger.Error("getTripErr")
		return nil, status.Error(codes.Internal, "错误")
	}
	cur := trip.Trip.Current.Location
	if request.Current != nil {
		cur = request.Current.Location
	}
	trip.Trip.Current = s.calcCurrentStatus(trip.Trip.Current, cur)
	if request.EndTrip {
		trip.Trip.End = trip.Trip.Current
		trip.Trip.Status = rentalpb.TripStatus_FINISHED
	}
	updateTripErr := s.Mongo.UpdateTrip(ctx, trip.ID.Hex(), userId, trip.UpdateAt, trip.Trip)
	if updateTripErr != nil {
		return nil, status.Error(codes.Internal, updateTripErr.Error())
	}
	return nil, nil
}

func (s *Service) DelTrip(ctx context.Context, request *rentalpb.DelTripRequest) (*rentalpb.DelTripResponse, error) {
	//TODO implement me
	panic("implement me")
}

type ProfileManager interface {
	Verify(ctx context.Context, userId string) (string, error)
}

type CarManager interface {
	Verify(ctx context.Context, carId string, location *rentalpb.Location) error
	Unlock(ctx context.Context, carId string) error
}

type PoiManager interface {
	GetPoiName(l *rentalpb.Location) (string, error)
}

func (s *Service) CreateTrip(ctx context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	userId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("hello world", zap.String("userId", userId))
	if request.CarId == "" || request.Start == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
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

	createTrip, createTripErr := s.Mongo.CreateTrip(ctx, &rentalpb.Trip{
		UserId: userId,
		CarId:  request.CarId,
		Start: &rentalpb.LocationStatus{
			Location:     request.Start,
			FeeCent:      0,
			KmDriven:     0,
			LocationDesc: poiName,
		},
		Current: &rentalpb.LocationStatus{
			Location:     request.Start,
			FeeCent:      0,
			KmDriven:     0,
			LocationDesc: poiName,
		},
		End:        nil,
		Status:     rentalpb.TripStatus_IN_PROGRESS,
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
	return &rentalpb.TripEntity{
		Id:   createTrip.ID.Hex(),
		Trip: createTrip.Trip,
	}, nil
}

const centsPerSec = 0.7

func nowFunc() int64 {
	return time.Now().Unix()
}
func (s Service) calcCurrentStatus(last *rentalpb.LocationStatus, cur *rentalpb.Location) *rentalpb.LocationStatus {
	esc := float64(nowFunc() - last.TimestampSec)
	return &rentalpb.LocationStatus{
		Location:     cur,
		FeeCent:      last.FeeCent + int32(esc*centsPerSec),
		KmDriven:     12,
		LocationDesc: "秦家坝",
	}
}
