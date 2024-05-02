package common

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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

func TimeStampToString(timeStamp time.Time) string {
	return timeStamp.Format("2006-01-02 15:04:05")
}

func StringToTimeStamp(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}

func generateMessage(username string, functionName string, message string, data any) string {
	return fmt.Sprintf("FunctionName: %s, username: %s, message: %s, data:%+v", functionName, username, message, data)
}

func generateErrorMessage(username string, functionName string, message string, err error, data any) string {
	return fmt.Sprintf("FunctionName: %s, username: %s, message: %s, error:%v, data:%+v", functionName, username, message, err, data)
}

func InfoLogger(username string, functionName string, message string, data ...interface{}) {
	logrus.Info(generateMessage(username, functionName, message, data))
}

func WarnLogger(username string, functionName string, message string, err error, data ...interface{}) {
	logrus.Warn(generateErrorMessage(username, functionName, message, err, data))
}

func ErrorLogger(username string, functionName string, message string, err error, data ...interface{}) {
	logrus.Error(generateErrorMessage(username, functionName, message, err, data))
}
