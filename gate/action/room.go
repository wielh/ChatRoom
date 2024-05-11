package action

import (
	"common"
	"context"
	"errorCode"
	"net/http"

	pb "proto"
	sc "serviceClient"

	"github.com/gin-gonic/gin"
)

type CreateRoomRequest struct {
	RoomName string `json:"room_name"`
}

func CreateRoom(c *gin.Context) {
	var createRoomRequest CreateRoomRequest
	if err := c.BindJSON(&createRoomRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.CreateRoomRequest{
		UserID:   userID,
		RoomName: createRoomRequest.RoomName,
	}
	response, err := sc.RoomServiceClient.CreateRoom(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.CreateRoom", "create room error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

type DeleteRoomRequest struct {
	RoomID string `json:"room_id"`
}

func DeleteRoom(c *gin.Context) {
	var deleteRoomRequest DeleteRoomRequest
	if err := c.BindJSON(&deleteRoomRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.DeleteRoomRequest{
		AdminID: userID,
		RoomID:  deleteRoomRequest.RoomID,
	}
	response, err := sc.RoomServiceClient.DeleteRoom(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.DeleteRoom", "micro-service error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetRoomsInfoByAdminID(c *gin.Context) {
	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.GetRoomsInfoByAdminIDRequest{
		AdminID: userID,
	}
	response, err := sc.RoomServiceClient.GetRoomsInfoByAdminID(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.GetRoomsInfoByAdminID", "micro-service error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetRoomsInfoByUserID(c *gin.Context) {
	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.GetRoomsInfoByUserIDRequest{
		UserID: userID,
	}
	response, err := sc.RoomServiceClient.GetRoomsInfoByUserID(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.GetRoomsInfoByUserID", "micro-service error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

type GetRoomInfoRequest struct {
	RoomID string `json:"room_id"`
}

func GetRoomInfo(c *gin.Context) {
	var getRoomInfoRequest GetRoomInfoRequest
	if err := c.BindJSON(&getRoomInfoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	userID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.GetRoomInfoRequest{
		UserID: userID,
		RoomID: getRoomInfoRequest.RoomID,
	}
	response, err := sc.RoomServiceClient.GetRoomInfo(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.CreateRoom", "micro-service error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

type AddUserRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func AddUserToRoom(c *gin.Context) {
	var addUserRequest AddUserRequest
	if err := c.BindJSON(&addUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	adminID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.AddUserRequest{
		AdminID: adminID,
		UserID:  addUserRequest.UserID,
		RoomID:  addUserRequest.RoomID,
	}
	response, err := sc.RoomServiceClient.AddUser(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.AddUser", "micro-service error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}

type DeleteUserRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func DeleteUserFromRoom(c *gin.Context) {
	var deleteUserRequest DeleteUserRequest
	if err := c.BindJSON(&deleteUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseJsonFailed})
		return
	}

	adminID, exist := common.GetUserID(c)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"errcode": errorCode.ErrParseToken})
		return
	}

	request := &pb.DeleteUserRequest{
		AdminID: adminID,
		UserID:  deleteUserRequest.UserID,
		RoomID:  deleteUserRequest.RoomID,
	}

	response, err := sc.RoomServiceClient.DeleteUser(context.Background(), request)
	if err != nil {
		common.ErrorLogger("gate", "sc.RoomServiceClient.DeleteUser", "micro-service error", err, request)
		c.JSON(http.StatusInternalServerError, gin.H{"errcode": errorCode.ErrMicroServiceNotResponse})
		return
	}
	c.JSON(http.StatusOK, response)
}
