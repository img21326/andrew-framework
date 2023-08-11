package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var JWTHelperInstance *JWTHelper

type JWTHelper struct {
	serverKey string
}

func newJwtHelper() *JWTHelper {
	serverKey := viper.GetViper().GetString("SERVER_KEY")
	return &JWTHelper{
		serverKey: serverKey,
	}
}

func GetJWTHelper() *JWTHelper {
	if JWTHelperInstance == nil {
		JWTHelperInstance = newJwtHelper()
	}
	return JWTHelperInstance
}

func (s *JWTHelper) CreateToken(userModel interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userModel,
		"exp":  time.Now().Add(time.Hour * 24 * 5).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(s.serverKey))
	return tokenString
}

func (s *JWTHelper) VerifyToken(token string) (interface{}, error) {
	verifyToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.serverKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	expireTime := claims["exp"].(float64)
	if time.Now().Unix() > int64(expireTime) {
		return nil, jwt.ErrTokenExpired
	}
	return claims["user"], nil
}
