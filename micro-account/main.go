package main

import (
	"common"
	"context"
	dbstructure "dbStructure"
	"fmt"
	"net"
	"proto"

	"google.golang.org/grpc"
)

type accountServer struct {
	proto.UnimplementedAccountServiceServer
}

func (s *accountServer) GoogleLogin(ctx context.Context, in *proto.GooogleLoginRequest) (*proto.GooogleLoginResponse, error) {
	birthTimeStamp, err := common.StringToTimeStamp(in.Birth)
	if err != nil {
		common.WarnLogger("micro-account", "GoogleLogin", "convert birth time stamp error", err, in)
		return &proto.GooogleLoginResponse{Errcode: common.ErrParameters}, nil
	}

	_, err = dbstructure.GoogleUserModel.SelectUser(in.GoogleID)
	if err == common.ErrNoRows {
		err = dbstructure.GoogleUserModel.InsertUser(in.GoogleID, in.FirstName, in.LastName, in.Sex, in.Email, birthTimeStamp)
		if err != nil {
			common.ErrorLogger("micro-account", "GoogleLogin", "Insert user to DB error", err, in)
			return &proto.GooogleLoginResponse{Errcode: common.ErrDBOther}, nil
		}
	} else if err != nil {
		common.ErrorLogger("micro-account", "GoogleLogin", "Select user from DB error", err, in)
		return &proto.GooogleLoginResponse{Errcode: common.ErrDBOther}, nil
	}
	return &proto.GooogleLoginResponse{Errcode: common.ErrSuccess, Token: common.CreateToken(in.GoogleID, 0)}, nil
}

func (s *accountServer) GetGoogleUserInfo(ctx context.Context, in *proto.GetGoogleUserInfoRequest) (*proto.GetGoogleUserInfoResponse, error) {
	row, err := dbstructure.GoogleUserModel.SelectUser(in.GoogleID)
	if err == common.ErrNoRows {
		return &proto.GetGoogleUserInfoResponse{Errcode: common.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-account", "GetGoogleUserInfo", "Select user from DB error", err, in)
		return &proto.GetGoogleUserInfoResponse{Errcode: common.ErrDBOther}, nil
	}

	return &proto.GetGoogleUserInfoResponse{
		Errcode:        common.ErrSuccess,
		GoogleID:       row.GoogleId,
		FirstName:      row.FirstName,
		LastName:       row.LastName,
		Sex:            row.Sex,
		Email:          row.Email,
		CreateDateTime: row.CreateDatetime.String(),
		Birth:          common.TimeStampToString(row.Birth),
	}, nil
}

func main() {
	common.ConfigInit()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", common.MicroAccountPort))
	if err != nil {
		common.ErrorLogger("micro-account", "main", fmt.Sprintf("Failed to listen port %v", common.MicroAccountPort), err)
		return
	}

	s := grpc.NewServer()
	proto.RegisterAccountServiceServer(s, &accountServer{})
	common.InfoLogger("micro-account", "main", "Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		common.ErrorLogger("micro-account", "main", "Starting gRPC server failed", err)
		return
	}
}
