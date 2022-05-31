package token

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTTokenGen struct {
	issuer     string
	nowFunc    func() time.Time
	privateKey *rsa.PrivateKey
}

func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

func (J *JWTTokenGen) GenToken(id string, expIn time.Duration) (string, error) {
	nowSec := J.nowFunc().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    J.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expIn.Seconds()),
		Subject:   id,
	})
	signedString, err := token.SignedString(J.privateKey)
	if err != nil {
		return "", err
	}
	return signedString, nil
}
