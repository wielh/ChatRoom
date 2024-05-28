package main

import (
	"common"
	"context"
	dbstructure "dbStructure"
	"errorCode"
	"fmt"
	"log"
	"net"
	"proto"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type chatServiceServer struct {
	proto.UnimplementedChatServiceServer
}

func (s *chatServiceServer) PushMessage(ctx context.Context, in *proto.PushMessageRequest) (*proto.PushMessageResponse, error) {
	_, err := dbstructure.RoomModel.GetRoomInfo(in.UserID, in.RoomID)
	if err == common.ErrNoRows {
		return &proto.PushMessageResponse{Errcode: errorCode.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomInfo", "Get rooms from DB error", err, in)
		return &proto.PushMessageResponse{Errcode: errorCode.ErrDBOther}, nil
	}

	err = dbstructure.MessageModel.PushMessage(in.UserID, in.Username, in.RoomID, in.Content)
	if err != nil {
		common.ErrorLogger("micro-chat", "PushMessage", "Create message from DB error", err, in)
		return &proto.PushMessageResponse{Errcode: errorCode.ErrDBOther}, nil
	}

	return &proto.PushMessageResponse{Errcode: errorCode.ErrSuccess}, nil
}

func (s *chatServiceServer) GetChatContent(in *proto.GetChatContentRequest, out proto.ChatService_GetChatContentServer) (err error) {
	var timeStamp time.Time
	if in.LastMessageTimeStamp == "" {
		timeStamp = time.Now()
	} else {
		t := strings.Replace(in.LastMessageTimeStamp, "\n", "", -1)
		t = strings.Replace(t, "\r", "", -1)
		timeStamp, err = common.StringToTimeStamp(t)
		if err != nil {
			common.WarnLogger("micro-chat", "GetChatContent", "Parse timestamp from string failed", err, in)
			return
		}
	}

	for {
		messages, err := dbstructure.MessageModel.GetMessages(in.UserID, timeStamp)
		if err != nil {
			common.ErrorLogger("micro-chat", "dbstructure.MessageModel.GetMessages", "Get message from DB error", err, in)
			return err
		}
		fmt.Println("length:", len(messages))

		var singalResponse = &proto.GetChatContentResponse{}
		if len(messages) == 0 {
			timeStamp = time.Now()
			singalResponse.LastMessageTimeStamp = common.TimeStampToString(timeStamp)
		} else {
			singalResponse.Messages = []*proto.ChatMessage{}
			for _, message := range messages {
				timeStamp = message.Time
				singalResponse.Messages = append(singalResponse.Messages, &proto.ChatMessage{
					ID:        message.ID,
					UserID:    message.UserID,
					RoomID:    message.RoomID,
					Content:   message.Content,
					TimeStamp: common.TimeStampToString(timeStamp),
				})
				singalResponse.LastMessageTimeStamp = common.TimeStampToString(timeStamp)
			}
		}

		err = out.Send(singalResponse)
		if err != nil {
			common.InfoLogger("micro-chat", "out.Send", "Send error:", err)
			return err
		}
		time.Sleep(3 * time.Second)
	}
}

func main() {
	common.ConfigInit()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", common.MicroChatPort))
	if err != nil {
		common.ErrorLogger("micro-chat", "main", fmt.Sprintf("Failed to listen port %v", common.MicroAccountPort), err)
		return
	}
	s := grpc.NewServer()

	proto.RegisterChatServiceServer(s, &chatServiceServer{})
	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	common.InfoLogger("micro-chat", "main", "Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		common.ErrorLogger("micro-chat", "main", "Starting gRPC server failed", err)
		return
	}
}
