package main

import (
	"action"

	"github.com/gin-gonic/gin"
)

func setRouter(router *gin.Engine) {
	account := router.Group("/account")
	account.GET("/google_login", action.LoginWithGoogleOAuth)
	account.GET("/google_callback", action.GoogleCallback)
	account.Use(action.AuthMiddleware)
	account.GET("/user_info", action.GetUserInfo)

	room := router.Group("/room")
	room.Use(action.AuthMiddleware)
	room.PUT("/", action.CreateRoom)
	room.DELETE("/", action.DeleteRoom)
	room.GET("/admin", action.GetRoomsInfoByAdminID)
	room.GET("/user", action.GetRoomsInfoByUserID)
	room.GET("/", action.GetRoomInfo)
	room.PUT("/user", action.AddUserToRoom)
	room.DELETE("/user", action.DeleteUserFromRoom)

	chat := router.Group("/chat")
	chat.Use(action.AuthMiddleware)
	chat.POST("/message", action.PushMessage)
	chat.GET("/message", action.GetChatContext)
}
