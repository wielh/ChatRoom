package common

import (
	"bytes"
	"encoding/json"
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

type LoggerDocument struct {
	FunctionName string `json:"FunctionName"`
	Message      string `json:"message"`
	Level        string `json:"level"`
	ErrMsg       string `json:"errMsg"`
	Data         []any  `json:"data"`
	Time         string `json:"time"`
}

func InfoLogger(username string, functionName string, message string, data ...interface{}) {
	logrus.Info(generateMessage(username, functionName, message, data))
	go func() {
		docStr, err := json.Marshal(LoggerDocument{
			FunctionName: functionName,
			Message:      message,
			Level:        "info",
			Data:         data,
			Time:         time.Now().Format("2006-01-02 15:04:05.111"),
		})

		if err != nil {
			return
		}
		ElasticClient.Index(username, bytes.NewReader(docStr))
	}()
}

func WarnLogger(username string, functionName string, message string, err error, data ...interface{}) {
	logrus.Warn(generateErrorMessage(username, functionName, message, err, data))
	go func() {
		docStr, _ := json.Marshal(LoggerDocument{
			FunctionName: functionName,
			Message:      message,
			Level:        "warn",
			ErrMsg:       err.Error(),
			Data:         data,
			Time:         time.Now().Format("2006-01-02 15:04:05.111"),
		})

		ElasticClient.Index(username, bytes.NewReader(docStr))
	}()
}

func ErrorLogger(username string, functionName string, message string, err error, data ...interface{}) {
	logrus.Error(generateErrorMessage(username, functionName, message, err, data))
	go func() {
		docStr, err := json.Marshal(LoggerDocument{
			FunctionName: functionName,
			Message:      message,
			Level:        "error",
			ErrMsg:       err.Error(),
			Data:         data,
			Time:         time.Now().Format("2006-01-02 15:04:05.111"),
		})

		if err != nil {
			return
		}

		ElasticClient.Index(username, bytes.NewReader(docStr))
	}()
}
