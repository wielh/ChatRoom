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

type roomServiceServer struct {
	proto.UnimplementedRoomServiceServer
}

func (s *roomServiceServer) CreateRoom(ctx context.Context, in *proto.RoomCreateRequest) (out *proto.RoomCreateResponse, err error) {
	err = dbstructure.RoomModel.RoomCreate(in.UserID, in.RoomName)
	if err != nil {
		out = &proto.RoomCreateResponse{Errcode: common.ErrDBOther}
		return
	}
	out = &proto.RoomCreateResponse{Errcode: common.ErrSuccess}
	return
}

func (s *roomServiceServer) DeleteRoom(ctx context.Context, in *proto.RoomDeleteRequest) (out *proto.RoomDeleteResponse, err error) {
	err = dbstructure.RoomModel.RoomDelete(in.AdminID, in.RoomID)
	if err == common.ErrNoRows {
		out = &proto.RoomDeleteResponse{Errcode: common.ErrDBDataNotFound}
		return
	} else if err != nil {
		out = &proto.RoomDeleteResponse{Errcode: common.ErrDBOther}
		return
	}
	out = &proto.RoomDeleteResponse{Errcode: common.ErrSuccess}
	return
}

func (s *roomServiceServer) GetRoomsInfoByAdminID(ctx context.Context, in *proto.GetRoomsInfoByAdminIDRequest) (out *proto.GetRoomsInfoByAdminIDResponse, err error) {
	rooms, err := dbstructure.RoomModel.GetRoomsInfoByAdminID(in.AdminID)
	if err == common.ErrNoRows || len(rooms) == 0 {
		out = &proto.GetRoomsInfoByAdminIDResponse{Errcode: common.ErrDBDataNotFound}
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
	if err == common.ErrNoRows || len(rooms) == 0 {
		out = &proto.GetRoomsInfoByUserIDResponse{Errcode: common.ErrDBDataNotFound}
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
	lis, err := net.Listen("tcp", common.LocalIP(common.MicroRoomPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterRoomServiceServer(s, &roomServiceServer{})
	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
