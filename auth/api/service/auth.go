package service

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"server2/auth/api/gen/v1"
	"server2/auth/dao"
	"time"
)

type Service struct {
	Logger         *zap.Logger
	OpenIdResolver OpenIdResolver
	Mongo          *dao.Mongo
	TokenGen       TokenGen
	TokenExp       time.Duration
	authpb.UnimplementedAuthServiceServer
}

type OpenIdResolver interface {
	Resolve(code string) (string, error)
}

type TokenGen interface {
	GenToken(id string, expIn time.Duration) (string, error)
}

func (s *Service) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	r, err := s.OpenIdResolver.Resolve(request.Code)
	if err != nil {
		s.Logger.Error("登录失败", zap.Error(err))
		return nil, status.Errorf(codes.Unavailable, "找不到openId %v", err)
	}
	s.Logger.Info("接收到code", zap.String("code", request.Code))
	id, err := s.Mongo.ResolveAuthId(ctx, r)
	if err != nil {
		s.Logger.Error("获取id失败", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	token, err := s.TokenGen.GenToken(id, s.TokenExp)
	if err != nil {
		s.Logger.Error("生成token失败", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "")
	}
	return &authpb.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int32(s.TokenExp.Seconds()),
	}, nil
}
