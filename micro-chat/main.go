package main

import (
	"common"
	dbstructure "dbStructure"
	"log"
	"net"
	"proto"
	"time"

	"google.golang.org/grpc"
)

type chatServiceServer struct {
	proto.UnimplementedChatServiceServer
}

func (s *chatServiceServer) GetChatContent(in *proto.GetChatContentRequest, out proto.ChatService_GetChatContentServer) (err error) {
	timeStamp, err := time.Parse("2006-01-02 15:04:05", in.LastMessageTimeStamp)
	if err != nil {
		log.Println(err)
	}

	for {
		messages, err := dbstructure.MessageModel.GetMessages(in.UserID, timeStamp)
		if err != nil {
			return err
		}

		var singalResponse = &proto.GetChatContentResponse{}
		singalResponse.Messages = make([]*proto.ChatMessage, len(messages))
		for _, message := range messages {
			singalResponse.Messages = append(singalResponse.Messages, &proto.ChatMessage{
				ID:        message.ID,
				UserID:    message.UserID,
				RoomID:    message.RoomID,
				Content:   message.Content,
				TimeStamp: message.Time.Format("2006-01-02 15:04:05"),
			})
		}

		err = out.Send(singalResponse)
		if err != nil {
			log.Fatalf("Send error:%v", err)
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", common.LocalIP(common.MicroChatPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterChatServiceServer(s, &chatServiceServer{})
	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
