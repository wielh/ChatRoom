package main

import (
	"common"
	"context"
	dbstructure "dbStructure"
	"fmt"
	"log"
	"net"
	"proto"
	"time"

	"google.golang.org/grpc"
)

type chatServiceServer struct {
	proto.UnimplementedChatServiceServer
}

func (s *chatServiceServer) PushMessage(ctx context.Context, in *proto.PushMessageRequest) (*proto.PushMessageResponse, error) {
	_, err := dbstructure.RoomModel.GetRoomInfo(in.UserID, in.RoomID)
	if err == common.ErrNoRows {
		return &proto.PushMessageResponse{Errcode: common.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomInfo", "Get rooms from DB error", err, in)
		return &proto.PushMessageResponse{Errcode: common.ErrDBOther}, nil
	}

	err = dbstructure.MessageModel.PushMessage(in.UserID, in.Username, in.RoomID, in.Content)
	if err != nil {
		common.ErrorLogger("micro-chat", "PushMessage", "Create message from DB error", err, in)
		return &proto.PushMessageResponse{Errcode: common.ErrDBOther}, nil
	}

	return &proto.PushMessageResponse{Errcode: common.ErrSuccess}, nil
}

func (s *chatServiceServer) GetChatContent(in *proto.GetChatContentRequest, out proto.ChatService_GetChatContentServer) (err error) {
	var timeStamp time.Time
	if in.LastMessageTimeStamp == "" {
		timeStamp = time.Now()
	} else {
		timeStamp, err = time.Parse("2006-01-02 15:04:05", in.LastMessageTimeStamp)
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

		var singalResponse = &proto.GetChatContentResponse{}
		singalResponse.Messages = make([]*proto.ChatMessage, 0)
		if len(messages) == 0 {
			singalResponse.LastMessageTimeStamp = common.TimeStampToString(time.Now())
		} else {
			for _, message := range messages {
				singalResponse.Messages = append(singalResponse.Messages, &proto.ChatMessage{
					ID:        message.ID,
					UserID:    message.UserID,
					RoomID:    message.RoomID,
					Content:   message.Content,
					TimeStamp: common.TimeStampToString(message.Time),
				})
				singalResponse.LastMessageTimeStamp = message.Time.Format("2006-01-02 15:04:05")
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
