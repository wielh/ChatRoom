package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ToCommonID(ID string, userType UserType) string {
	switch userType {
	case GoogleUser:
		return "googleID::" + ID
	default:
		return ""
	}
}

func CreateToken(userID string, userType int32) string {

	var answer string = ToCommonID(userID, UserType(userType))
	if answer == "" {
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": answer,
		"exp":      time.Now().Add(24 * time.Hour),
	})

	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return ""
	}
	return tokenString
}

func LocalIP(port int32) string {
	return "127.0.0.1:" + string(port)
}
