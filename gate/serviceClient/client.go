package serviceClient

import (
	"common"
	"fmt"
	"proto"

	"google.golang.org/grpc"
)

var AccountServiceClient proto.AccountServiceClient
var ChatServiceClient proto.ChatServiceClient
var RoomServiceClient proto.RoomServiceClient

func grpcConnect(port int32) (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
}

func MicroServiceClientInit() error {
	accountServiceConn, err := grpcConnect(common.MicroAccountPort)
	if err != nil {
		return err
	}
	AccountServiceClient = proto.NewAccountServiceClient(accountServiceConn)

	chatServiceConn, err := grpcConnect(common.MicroChatPort)
	if err != nil {
		return err
	}
	ChatServiceClient = proto.NewChatServiceClient(chatServiceConn)

	RoomServiceConn, err := grpcConnect(common.MicroRoomPort)
	if err != nil {
		return err
	}
	RoomServiceClient = proto.NewRoomServiceClient(RoomServiceConn)
	return nil
}
