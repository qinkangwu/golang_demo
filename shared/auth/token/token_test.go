package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

const token = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoic2VydmVyMiIsInN1YiI6IjEyMzQ1Njc4OTAifQ.ba5upL0CO4kxa7pIX_7hHOpgAQwu2LxhHdmzog6OsWDiQZGAvBLd4MWMYI4e4dONuNhKAOxLDowLHKui4pm3aS6OxLcifvv3zKJmURF0wWzuzp5pdkmjCoDSzXTyy-ske0s5Pg4Xfhv_LSe71ddinnC8YNmRXtJRAONNjceySmcNmHwEhH09Az3fHKSjQkWkMk1MgRu19p3XAaNBpbmY6RTzgG1CA_GU-8umNwW94fdqE3JaeH2WLWTxVVUzE4baJhlSSUSI2YXodaWjNC5kDF9ECTSigogd6NV0wOWQJ9moFBVBfJRdL_601ZpJKHyrXDSN918wA-zcuJHF5vFz8g`

func TestJWTVerifier_Verify(t *testing.T) {
	publicKeyFromPEM, parseErr := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if parseErr != nil {
		t.Errorf("ParseRSAPublicKeyFromPEM - 错误 %v", parseErr)
	}
	v := &JWTVerifier{
		PublicKey: publicKeyFromPEM,
	}

	_, vError := v.Verify(token)
	if vError != nil {
		t.Errorf("Verify - 错误 %v", vError)
	}
}
