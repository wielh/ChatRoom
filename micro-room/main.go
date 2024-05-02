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

type roomServiceServer struct {
	proto.UnimplementedRoomServiceServer
}

func (s *roomServiceServer) CreateRoom(ctx context.Context, in *proto.CreateRoomRequest) (out *proto.CreateRoomResponse, err error) {
	err = dbstructure.RoomModel.RoomCreate(in.UserID, in.RoomName)
	if err == common.ErrNoRows {
		out = &proto.CreateRoomResponse{Errcode: common.ErrDBDataAlreadyExist}
		return
	} else if err != nil {
		common.ErrorLogger("micro-room", "CreateRoom", "Create room from DB error", err, in)
		out = &proto.CreateRoomResponse{Errcode: common.ErrDBOther}
		return
	}
	out = &proto.CreateRoomResponse{Errcode: common.ErrSuccess}
	return
}

func (s *roomServiceServer) DeleteRoom(ctx context.Context, in *proto.DeleteRoomRequest) (out *proto.DeleteRoomResponse, err error) {
	err = dbstructure.RoomModel.RoomDeleteTransection(in.AdminID, in.RoomID, ctx)
	if err == common.ErrNoRows {
		out = &proto.DeleteRoomResponse{Errcode: common.ErrDBDataNotFound}
		return
	} else if err != nil {
		common.ErrorLogger("micro-room", "DeleteRoom", "Delete room from DB error", err, in)
		out = &proto.DeleteRoomResponse{Errcode: common.ErrDBOther}
		return
	}
	out = &proto.DeleteRoomResponse{Errcode: common.ErrSuccess}
	return
}

func (s *roomServiceServer) GetRoomsInfoByAdminID(ctx context.Context, in *proto.GetRoomsInfoByAdminIDRequest) (out *proto.GetRoomsInfoByAdminIDResponse, err error) {
	rooms, err := dbstructure.RoomModel.GetRoomsInfoByAdminID(in.AdminID)
	if err == common.ErrNoRows {
		out = &proto.GetRoomsInfoByAdminIDResponse{Errcode: common.ErrDBDataNotFound}
		return
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomsInfoByAdminID", "Get rooms from DB error", err, in)
		out = &proto.GetRoomsInfoByAdminIDResponse{Errcode: common.ErrDBOther}
		return
	}

	out = &proto.GetRoomsInfoByAdminIDResponse{}
	out.Errcode = common.ErrSuccess
	out.RoomsInfo = []*proto.RoomInfo{}
	for _, room := range rooms {
		out.RoomsInfo = append(out.RoomsInfo, &proto.RoomInfo{
			ID:      room.ID,
			Name:    room.Name,
			AdminID: room.AdminID,
			UsersID: room.UsersID,
		})
	}
	return
}

func (s *roomServiceServer) GetRoomsInfoByUserID(ctx context.Context, in *proto.GetRoomsInfoByUserIDRequest) (out *proto.GetRoomsInfoByUserIDResponse, err error) {
	rooms, err := dbstructure.RoomModel.GetRoomsInfoByUserID(in.UserID)
	//|| len(rooms) == 0
	if err == common.ErrNoRows {
		out = &proto.GetRoomsInfoByUserIDResponse{Errcode: common.ErrDBDataNotFound}
		return
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomsInfoByUserID", "Get rooms from DB error", err, in)
		out = &proto.GetRoomsInfoByUserIDResponse{Errcode: common.ErrDBOther}
		return
	}

	out = &proto.GetRoomsInfoByUserIDResponse{}
	out.Errcode = common.ErrSuccess
	out.RoomsInfo = []*proto.RoomInfo{}
	for _, room := range rooms {
		out.RoomsInfo = append(out.RoomsInfo, &proto.RoomInfo{
			ID:      room.ID,
			Name:    room.Name,
			AdminID: room.AdminID,
			UsersID: room.UsersID,
		})
	}
	return
}

func (s *roomServiceServer) GetRoomInfo(ctx context.Context, in *proto.GetRoomInfoRequest) (out *proto.GetRoomInfoResponse, err error) {
	room, err := dbstructure.RoomModel.GetRoomInfo(in.UserID, in.RoomID)
	if err == common.ErrNoRows {
		out = &proto.GetRoomInfoResponse{Errcode: common.ErrDBDataNotFound}
		return
	} else if err != nil {
		common.ErrorLogger("micro-room", "GetRoomInfo", "Get rooms from DB error", err, in)
		out = &proto.GetRoomInfoResponse{Errcode: common.ErrDBOther}
		return
	}

	out = &proto.GetRoomInfoResponse{
		Errcode: common.ErrSuccess,
		RoomInfo: &proto.RoomInfo{
			ID:      room.ID,
			Name:    room.Name,
			AdminID: room.AdminID,
			UsersID: room.UsersID,
		},
	}
	return
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
