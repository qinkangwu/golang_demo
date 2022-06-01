package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"os"
	"server2/shared/auth/token"
	"strings"
)

func Interceptor(publicFile string) (grpc.UnaryServerInterceptor, error) {
	file, osOpenErr := os.Open(publicFile)
	if osOpenErr != nil {
		return nil, osOpenErr
	}
	readAll, readAllErr := ioutil.ReadAll(file)
	if readAllErr != nil {
		return nil, readAllErr
	}
	pem, parseErr := jwt.ParseRSAPublicKeyFromPEM(readAll)
	if parseErr != nil {
		return nil, parseErr
	}
	i := &interceptor{
		verifier: &token.JWTVerifier{
			PublicKey: pem,
		},
	}
	return i.HandleRequest, nil
}

type interceptor struct {
	verifier tokenVerifier
}

type tokenVerifier interface {
	Verify(token string) (string, error)
}

func (i *interceptor) HandleRequest(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	fromContext, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "身份验证不通过")
	}
	userId, verifyErr := i.verifier.Verify(fromContext)
	if verifyErr != nil {
		return nil, status.Error(codes.Unauthenticated, "身份验证不通过")
	}

	return handler(ContextWithUserId(ctx, userId), req)
}

func tokenFromContext(c context.Context) (string, error) {
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "身份验证不通过")
	}
	tkn := ""
	for _, v := range m["authorization"] {
		if strings.HasPrefix(v, "Bearer ") {
			tkn = v[len("Bearer "):]
		}
	}
	if tkn == "" {
		return "", status.Error(codes.Unauthenticated, "身份验证不通过")
	}
	return tkn, nil
}

type userIdKey struct{}

func ContextWithUserId(c context.Context, userId string) context.Context {
	return context.WithValue(c, userIdKey{}, userId)
}

func UserIdFromContext(c context.Context) (string, error) {
	userId := c.Value(userIdKey{})
	if userId == "" {
		return "", fmt.Errorf("没有找到")
	}
	uId, ok := userId.(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "身份验证不通过")
	}
	return uId, nil
}
