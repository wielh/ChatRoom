package action

import (
	"common"
	"context"
	"encoding/json"
	"errorCode"
	"net/http"

	pb "proto"
	sc "serviceClient"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PushMessageRequest struct {
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Content  string `json:"content"`
}

func PushMessage(c *gin.Context) {
	var pushMessageRequest PushMessageRequest
	if err := c.BindJSON(&pushMessageRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.PushMessageRequest{
		UserID:   userID,
		Username: pushMessageRequest.Username,
		RoomID:   pushMessageRequest.RoomID,
		Content:  pushMessageRequest.Content,
	}
	response, err := sc.ChatServiceClient.PushMessage(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.ChatServiceClient.PushMessage", "push message error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

type GetChatContentRequest struct {
	TimeStamp string `json:"timeStamp"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetChatContext(c *gin.Context) {
	// connection between client and gate
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.WarnLogger("gate", "upgrader.Upgrade", "connection upgrade failed", err)
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrWebSocketUpgradeFailed})
		return
	}
	defer conn.Close()

	// get parameters
	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	_, timeStampByte, err := conn.ReadMessage()
	if err != nil {
		common.ErrorLogger("gate", "conn.ReadMessage", "read timestamp failed", err)
		return
	}

	// connection between gate and micro-service
	request := &pb.GetChatContentRequest{
		UserID:               userID,
		LastMessageTimeStamp: string(timeStampByte),
	}

	response, err := sc.ChatServiceClient.GetChatContent(context.Background(), request)
	defer response.CloseSend()
	if err != nil {
		common.ErrorLogger("gate", "sc.ChatServiceClient.GetChatContent", "push message error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}

	for {
		/*
			signalType, _, err := conn.ReadMessage()
			if err != nil || signalType == websocket.CloseMessage {
				common.InfoLogger("gate", " conn.ReadMessage", "websocket connection stop", err, signalType)
				return
			}*/

		msg, err := response.Recv()
		if err != nil {
			common.ErrorLogger("gate", "sc.ChatServiceClient.GetChatContent", "receive stream messages from micro-service failed", err)
			return
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			common.ErrorLogger("gate", "GetChatContext json.Marshal", "Marshal data failed", err, msg)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			common.ErrorLogger("gate", "GetChatContext conn.WriteMessage", "Websocket WriteMessage failed", err, msg)
			return
		}
	}
}
