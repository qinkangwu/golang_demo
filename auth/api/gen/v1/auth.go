package authpb

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger         *zap.Logger
	OpenIdResolver OpenIdResolver
}

type OpenIdResolver interface {
	Resolve(code string) (string, error)
}

func (s *Service) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	r, err := s.OpenIdResolver.Resolve(request.Code)
	if err != nil {
		s.Logger.Error("登录失败", zap.Error(err))
		return nil, status.Errorf(codes.Unavailable, "找不到openId %v", err)
	}
	s.Logger.Info("接收到code", zap.String("code", request.Code))
	return &LoginResponse{
		AccessToken: r,
		ExpiresIn:   7200,
	}, nil
}

func (s *Service) mustEmbedUnimplementedAuthServiceServer() {
	return
}
