package main

import (
	"common"
	"context"
	dbstructure "dbStructure"
	"log"
	"net"
	"proto"

	"google.golang.org/grpc"
)

type accountServer struct {
	proto.UnimplementedAccountServiceServer
}

func (s *accountServer) GoogleLogin(ctx context.Context, in *proto.GooogleLoginRequest) (*proto.GooogleLoginResponse, error) {
	_, err := dbstructure.GoogleUserModel.SelectUser(in.GoogleID)
	if err == common.ErrNoRows {
		err = dbstructure.GoogleUserModel.InsertUser(in.GoogleID, in.FirstName, in.LastName, in.Sex, in.Email, in.Age)
		if err != nil {
			return &proto.GooogleLoginResponse{Errcode: common.ErrDBOther}, err
		}
	} else if err != nil {
		return &proto.GooogleLoginResponse{Errcode: common.ErrDBOther}, err
	}

	return &proto.GooogleLoginResponse{Errcode: common.ErrSuccess, Token: common.CreateToken(in.GoogleID, 0)}, nil
}

func (s *accountServer) GetGoogleUserInfo(ctx context.Context, in *proto.GetGoogleUserInfoRequest) (*proto.GetGoogleUserInfoResponse, error) {
	row, err := dbstructure.GoogleUserModel.SelectUser(in.GoogleID)
	if err == common.ErrNoRows {
		return &proto.GetGoogleUserInfoResponse{Errcode: common.ErrDBDataNotFound}, nil
	} else if err != nil {
		return &proto.GetGoogleUserInfoResponse{Errcode: common.ErrDBOther}, err
	}

	return &proto.GetGoogleUserInfoResponse{
		Errcode:   common.ErrSuccess,
		GoogleID:  row.GoogleId,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Sex:       row.Sex,
		Email:     row.Email,
		Age:       row.Age,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", common.LocalIP(common.MicroAccountPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterAccountServiceServer(s, &accountServer{})
	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
