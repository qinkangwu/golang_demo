package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JWTVerifier struct {
	PublicKey *rsa.PublicKey
}

func (v *JWTVerifier) Verify(token string) (string, error) {
	parseWithClaims, withClainmsErr := jwt.ParseWithClaims(
		token,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return v.PublicKey, nil
		},
	)
	if withClainmsErr != nil {
		return "", withClainmsErr
	}
	if !parseWithClaims.Valid {
		return "", fmt.Errorf("该token验证不通过")
	}

	c, ok := parseWithClaims.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("未知错误")
	}
	claimsErr := c.Valid()
	if claimsErr != nil {
		return "", claimsErr
	}

	return c.Subject, nil
}
