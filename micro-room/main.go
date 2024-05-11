package main

import (
	"common"
	"context"
	dbstructure "dbStructure"
	"errorCode"
	"fmt"
	"net"
	"proto"

	"google.golang.org/grpc"
)

type roomServiceServer struct {
	proto.UnimplementedRoomServiceServer
}

func (s *roomServiceServer) CreateRoom(ctx context.Context, in *proto.CreateRoomRequest) (*proto.CreateRoomResponse, error) {
	err := dbstructure.RoomModel.RoomCreate(in.UserID, in.RoomName)
	if err == common.ErrNoRows {
		return &proto.CreateRoomResponse{Errcode: errorCode.ErrDBDataAlreadyExist}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "CreateRoom", "Create room from DB error", err, in)
		return &proto.CreateRoomResponse{Errcode: errorCode.ErrDBOther}, nil
	}
	return &proto.CreateRoomResponse{Errcode: errorCode.ErrSuccess}, nil
}

func (s *roomServiceServer) DeleteRoom(ctx context.Context, in *proto.DeleteRoomRequest) (*proto.DeleteRoomResponse, error) {
	err := dbstructure.RoomModel.RoomDeleteTransection(in.AdminID, in.RoomID, ctx)
	if err == common.ErrNoRows {
		return &proto.DeleteRoomResponse{Errcode: errorCode.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "DeleteRoom", "Delete room from DB error", err, in)
		return &proto.DeleteRoomResponse{Errcode: errorCode.ErrDBOther}, nil

	}
	return &proto.DeleteRoomResponse{Errcode: errorCode.ErrSuccess}, nil
}

func (s *roomServiceServer) GetRoomsInfoByAdminID(ctx context.Context, in *proto.GetRoomsInfoByAdminIDRequest) (*proto.GetRoomsInfoByAdminIDResponse, error) {
	rooms, err := dbstructure.RoomModel.GetRoomsInfoByAdminID(in.AdminID)
	if err == common.ErrNoRows {
		return &proto.GetRoomsInfoByAdminIDResponse{Errcode: errorCode.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomsInfoByAdminID", "Get rooms from DB error", err, in)
		return &proto.GetRoomsInfoByAdminIDResponse{Errcode: errorCode.ErrDBOther}, nil
	}

	out := &proto.GetRoomsInfoByAdminIDResponse{}
	out.Errcode = errorCode.ErrSuccess
	out.RoomsInfo = []*proto.RoomInfo{}
	for _, room := range rooms {
		out.RoomsInfo = append(out.RoomsInfo, &proto.RoomInfo{
			ID:      room.ID,
			Name:    room.Name,
			AdminID: room.AdminID,
			UsersID: room.UsersID,
		})
	}
	return out, nil
}

func (s *roomServiceServer) GetRoomsInfoByUserID(ctx context.Context, in *proto.GetRoomsInfoByUserIDRequest) (*proto.GetRoomsInfoByUserIDResponse, error) {
	rooms, err := dbstructure.RoomModel.GetRoomsInfoByUserID(in.UserID)
	if err == common.ErrNoRows {
		return &proto.GetRoomsInfoByUserIDResponse{Errcode: errorCode.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomsInfoByUserID", "Get rooms from DB error", err, in)
		return &proto.GetRoomsInfoByUserIDResponse{Errcode: errorCode.ErrDBOther}, nil
	}

	out := &proto.GetRoomsInfoByUserIDResponse{}
	out.Errcode = errorCode.ErrSuccess
	out.RoomsInfo = []*proto.RoomInfo{}
	for _, room := range rooms {
		out.RoomsInfo = append(out.RoomsInfo, &proto.RoomInfo{
			ID:      room.ID,
			Name:    room.Name,
			AdminID: room.AdminID,
			UsersID: room.UsersID,
		})
	}
	return out, nil
}

func (s *roomServiceServer) GetRoomInfo(ctx context.Context, in *proto.GetRoomInfoRequest) (*proto.GetRoomInfoResponse, error) {
	room, err := dbstructure.RoomModel.GetRoomInfo(in.UserID, in.RoomID)
	if err == common.ErrNoRows {
		return &proto.GetRoomInfoResponse{Errcode: errorCode.ErrDBDataNotFound}, nil
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomInfo", "Get rooms from DB error", err, in)
		return &proto.GetRoomInfoResponse{Errcode: errorCode.ErrDBOther}, nil
	}

	out := &proto.GetRoomInfoResponse{
		Errcode: errorCode.ErrSuccess,
		RoomInfo: &proto.RoomInfo{
			ID:      room.ID,
			Name:    room.Name,
			AdminID: room.AdminID,
			UsersID: room.UsersID,
		},
	}
	return out, err
}

func (s *roomServiceServer) AddUser(ctx context.Context, in *proto.AddUserRequest) (*proto.AddUserResponse, error) {
	var out *proto.AddUserResponse
	exist, err := dbstructure.UserModel.UserExist(in.UserID)
	if err != nil {
		common.ErrorLogger("micro-room", "dbstructure.UserModel.UserExist", "check user exist error", err, in)
		out = &proto.AddUserResponse{Errcode: errorCode.ErrDBOther}
		return out, nil
	} else if !exist {
		out = &proto.AddUserResponse{Errcode: errorCode.ErrUserNotExist}
		return out, nil
	}

	err = dbstructure.RoomModel.AddUser(in.AdminID, in.RoomID, in.UserID)
	if err != nil {
		common.ErrorLogger("micro-room", "dbstructure.RoomModel.AddUser", "add user to room error", err, in)
		out = &proto.AddUserResponse{Errcode: errorCode.ErrDBOther}
		return out, nil
	}
	out = &proto.AddUserResponse{Errcode: errorCode.ErrSuccess}
	return out, nil
}

func (s *roomServiceServer) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	var out *proto.DeleteUserResponse
	err := dbstructure.RoomModel.DeleteUser(in.AdminID, in.RoomID, in.UserID)
	if err != nil {
		common.ErrorLogger("micro-room", "dbstructure.RoomModel.DeleteUser", "delete user from room error", err, in)
		out = &proto.DeleteUserResponse{Errcode: errorCode.ErrDBOther}
		return out, nil
	}
	out = &proto.DeleteUserResponse{Errcode: errorCode.ErrSuccess}
	return out, nil
}

func main() {
	common.ConfigInit()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", common.MicroRoomPort))
	if err != nil {
		common.ErrorLogger("micro-room", "main", fmt.Sprintf("Failed to listen port %v", common.MicroRoomPort), err)
		return
	}

	s := grpc.NewServer()
	proto.RegisterRoomServiceServer(s, &roomServiceServer{})
	common.InfoLogger("micro-room", "main", "Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		common.ErrorLogger("micro-room", "main", "Starting gRPC server failed", err)
		return
	}
}
