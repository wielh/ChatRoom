package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func ToCommonID(ID string, userType UserType) string {
	switch userType {
	case GoogleUser:
		return "googleID:" + ID
	default:
		return ""
	}
}

func ToSpecificID(ID string) (string, UserType) {
	if strings.HasPrefix(ID, "googleID:") {
		return ID[len("googleID:"):], GoogleUser
	} else {
		return "", None
	}
}

func CreateToken(userID string, userType int32, username string) string {
	var id string = ToCommonID(userID, UserType(userType))
	if id == "" {
		return ""
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = id
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 过期时间 24 小时

	answer, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return ""
	}
	return answer
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func TimeStampToString(timeStamp time.Time) string {
	return timeStamp.Format("2006-01-02 15:04:05.000")
}

func StringToTimeStamp(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05.000", timeStr)
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

func GetUserID(c *gin.Context) (string, bool) {
	userID, exist := c.Get("userID")
	if !exist {
		return "", false
	}
	return fmt.Sprintf("%v", userID), true
}
