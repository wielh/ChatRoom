package dbstructure

import (
	c "common"
)

type Room struct {
	ID      string   `pg:"type:pk,column:id"`
	AdminID string   `pg:"column:adminId"`
	Name    string   `pg:"column:name"`
	UsersID []string `pg:"column:usersId"`
	Deleted bool     `pg:"column:deleted"`
}

type roomModel struct{}

var RoomModel = &roomModel{}

func (r roomModel) RoomCreate(adminId string, name string) error {
	model := &Room{
		AdminID: adminId,
		Name:    name,
		UsersID: []string{adminId},
	}
	_, err := c.DB.Model(model).OnConflict("(name) DO NOTHING").Insert()
	return err
}

func (r roomModel) RoomDelete(adminId string, name string) (err error) {
	model := &Room{
		Deleted: true,
	}
	_, err = c.DB.Model(model).Where("adminId=? AND name=?", adminId, name).Where("Deleted=?", false).Update()
	return
}

func (r roomModel) AddUser(adminId string, roomID string, userId string) error {
	model := &Room{
		UsersID: []string{adminId},
	}
	_, err := c.DB.Model(model).Where("adminId=?", adminId).Where("id=?", roomID).Where(
		"NOT EXISTS (SELECT 1 FROM unnest(usersId) AS elem WHERE elem = ?)", userId).Set(
		"ARRAY_APPEND(st_email, ?)", userId).Update()
	return err
}

func (r roomModel) DeleteUser(adminId string, roomID string, userId string) error {
	model := &Room{
		UsersID: []string{adminId},
	}
	_, err := c.DB.Model(model).Where("adminId=?", adminId).Where("id=?", roomID).Where(
		"EXISTS (SELECT 1 FROM unnest(usersId) AS elem WHERE elem = ?)", userId).Set(
		"ARRAY_REMOVE(st_email, ?)", userId).Update()
	return err
}

func (r roomModel) GetRoomsInfoByAdminID(adminId string) (roomsInfo []*Room, err error) {
	err = c.DB.Model(roomsInfo).Where("adminId=?", adminId).Select()
	return
}

func (r roomModel) GetRoomsInfoByUserID(userId string) (roomsInfo []*Room, err error) {
	err = c.DB.Model(roomsInfo).Where("EXISTS (SELECT 1 FROM unnest(usersId) AS elem WHERE elem = ?)", userId).Select()
	return
}

func (r roomModel) GetRoomInfo(userId string, roomId string) (Room *Room, err error) {
	err = c.DB.Model(Room).Where("EXISTS (SELECT 1 FROM unnest(usersId) AS elem WHERE elem = ?)", userId).Select()
	return
}
