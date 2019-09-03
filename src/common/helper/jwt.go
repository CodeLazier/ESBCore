package helper

import (
	"common"
	"github.com/dgrijalva/jwt-go"
)

func ParseJWT(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(common.JWTSecretKey), nil
	})

	if err==nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims, nil
		}
	}
	return nil, err
}
