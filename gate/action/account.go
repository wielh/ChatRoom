package action

import (
	"common"
	"context"
	"encoding/json"
	"errorCode"
	"io"

	"net/http"
	pb "proto"
	sc "serviceClient"

	"github.com/gin-gonic/gin"
)

type UserInfoRequest struct {
	UserID string `json:"user_id"`
}

func GetUserInfo(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	if err := c.BindJSON(&userInfoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	ID, userType := common.ToSpecificID(userInfoRequest.UserID)
	switch userType {
	case common.GoogleUser:
		request := &pb.GetGoogleUserInfoRequest{GoogleID: ID}
		response, err := sc.AccountServiceClient.GetGoogleUserInfo(context.Background(), request)
		if err != nil {
			common.ErrorLogger("gate", "sc.AccountServiceClient.GetGoogleUserInfo", "get google user error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
			return
		}
		c.JSON(http.StatusOK, response)
	default:
		c.JSON(http.StatusOK, gin.H{"errcode": errorCode.ErrParameters})
		return
	}
}

func LoginWithGoogleOAuth(c *gin.Context) {
	redirectURL := CreateGoogleOAuthURL()
	c.Redirect(http.StatusSeeOther, redirectURL)
}

func CreateGoogleOAuthURL() string {
	return common.GoogleOauth2Config.AuthCodeURL("state")
}

type UserInfo struct {
	GoogleID  string `json:"sub"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Email     string `json:"email"`
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := common.GoogleOauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errcode": errorCode.ErrGetGoogleToken})
	}

	client := common.GoogleOauth2Config.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		common.ErrorLogger("gate", "client.Get", "get userInfo from google failed", err)
		c.JSON(http.StatusOK, gin.H{"errcode": errorCode.ErrGetGoogleUserInfo})
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	var answer UserInfo
	err = json.Unmarshal(data, &answer)
	if err != nil {
		common.ErrorLogger("gate", "GoogleCallback json.Unmarshal", "Parse json error", err, string(data))
		c.JSON(http.StatusOK, gin.H{"errcode": errorCode.ErrToJsonErr})
	}

	request := &pb.GooogleLoginRequest{
		GoogleID: answer.GoogleID, FirstName: answer.FirstName, LastName: answer.LastName, Email: answer.Email,
	}
	response, err := sc.AccountServiceClient.GoogleLogin(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.AccountServiceClient.GoogleLogin", "micro-service error", err, request)
		c.JSON(http.StatusOK, gin.H{"errcode": response.Errcode})
	} else if response.Errcode != errorCode.ErrSuccess {
		c.JSON(http.StatusOK, gin.H{"errcode": response.Errcode})
	}

	c.SetCookie("token", response.Token, 86400, "/", "localhost:80", false, true)
	c.JSON(http.StatusOK, gin.H{"errcode": errorCode.ErrSuccess})
}

func AuthMiddleware(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrVerifyToken})
		c.Abort()
		return
	}

	claims, err := common.ParseJWT(tokenStr)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrVerifyToken})
		c.Abort()
		return
	}

	c.Set("userID", claims["userID"])
	c.Set("username", claims["username"])
	c.Next()
}
