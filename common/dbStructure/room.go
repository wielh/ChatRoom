package dbstructure

import (
	c "common"
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

type Room struct {
	ID      string `pg:"id,pk"`
	AdminID string
	Name    string
	UsersID []string `pg:",array"`
}

type roomModel struct{}

var RoomModel = &roomModel{}

type RoomHistory struct {
	ID             string `pg:"id,pk"`
	AdminID        string
	Name           string
	UsersID        []string `pg:",array"`
	CreateDatetime time.Time
}

type roomHistoryModel struct{}

var RoomHistoryModel = &roomHistoryModel{}

func (r roomModel) RoomCreate(adminId string, name string) error {
	model := &Room{
		AdminID: adminId,
		Name:    name,
		UsersID: []string{adminId},
	}
	_, err := c.DB.Model(model).OnConflict("(name) DO NOTHING").Insert()
	return err
}

func (r roomModel) roomDelete(tx *pg.Tx, adminId string, room_id string) (err error) {
	_, err = tx.Model(&Room{}).Where("admin_id=? AND id=?", adminId, room_id).Delete()
	return
}

func (r roomHistoryModel) roomHistoryCreate(tx *pg.Tx, adminId string, name string, UsersID []string) error {
	model := &RoomHistory{
		AdminID: adminId,
		Name:    name,
		UsersID: UsersID,
	}
	_, err := tx.Model(model).Insert()
	return err
}

func (r roomModel) RoomDeleteTransection(adminId string, roomId string, context context.Context) (err error) {
	err = c.DB.RunInTransaction(context, func(tx *pg.Tx) error {
		roomInfo, err := r.GetRoomInfoByAdminID(adminId, roomId)
		if err != nil {
			c.WarnLogger("common", " r.GetRoomInfoByAdminID", "Get room by admin_id and id failed", err, adminId, roomId)
			return err
		}

		err = r.roomDelete(tx, adminId, roomId)
		if err != nil {
			c.ErrorLogger("common", "r.roomDelete", "Failed to delete room", err, adminId, roomId)
			return err
		}

		err = RoomHistoryModel.roomHistoryCreate(tx, adminId, roomInfo.Name, roomInfo.UsersID)
		if err != nil {
			c.ErrorLogger("common", "RoomHistoryModel.roomHistoryCreate", "Failed to create room history", err, adminId, roomId, roomInfo)
			return err
		}
		return nil
	})
	return
}

func (r roomModel) AddUser(adminId string, roomID string, userId string) error {
	model := &Room{
		UsersID: []string{adminId},
	}
	_, err := c.DB.Model(model).Where("admin_id=?", adminId).Where("id=?", roomID).Where(
		"NOT EXISTS (SELECT 1 FROM unnest(users_id) AS elem WHERE elem = ?)", userId).Set(
		"ARRAY_APPEND(st_email, ?)", userId).Update()
	return err
}

func (r roomModel) DeleteUser(adminId string, roomID string, userId string) error {
	model := &Room{
		UsersID: []string{adminId},
	}
	_, err := c.DB.Model(model).Where("admin_id=?", adminId).Where("id=?", roomID).Where(
		"EXISTS (SELECT 1 FROM unnest(usersId) AS elem WHERE elem = ?)", userId).Set(
		"ARRAY_REMOVE(st_email, ?)", userId).Update()
	return err
}

func (r roomModel) GetRoomsInfoByAdminID(adminId string) ([]Room, error) {
	var roomsInfo []Room
	err := c.DB.Model(&roomsInfo).Where("admin_id=?", adminId).Select()
	return roomsInfo, err
}

func (r roomModel) GetRoomInfoByAdminID(adminId string, roomId string) (roomInfo *Room, err error) {
	roomInfo = &Room{}
	err = c.DB.Model(roomInfo).Where("admin_id=? and id=?", adminId, roomId).Select()
	return
}

func (r roomModel) GetRoomsInfoByUserID(userId string) ([]Room, error) {
	var roomsInfo []Room
	err := c.DB.Model(&roomsInfo).Where("EXISTS (SELECT 1 FROM unnest(users_id) AS elem WHERE elem = ?)", userId).Select()
	return roomsInfo, err
}

func (r roomModel) GetRoomInfo(userId string, roomId string) (roomInfo *Room, err error) {
	roomInfo = &Room{}
	err = c.DB.Model(roomInfo).Where("id=?", roomId).Where("EXISTS (SELECT 1 FROM unnest(users_id) AS elem WHERE elem = ?)", userId).Select()
	return
}
