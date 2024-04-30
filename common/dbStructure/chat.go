package dbstructure

import (
	c "common"
	"time"
)

type message struct {
	ID       string    `pg:"type:pk,column:id"`
	RoomID   string    `pg:"column:roomId"`
	UserID   string    `pg:"column:userId"`
	Username string    `pg:"column:username"`
	Time     time.Time `pg:"column:time"`
	Deleted  bool      `pg:"column:deleted"`
	Content  string    `pg:"column:content"`
}

type messageModel struct {
}

var MessageModel = &messageModel{}

func (m messageModel) GetMessages(userId string, timeStamp time.Time) (messages []*message, err error) {
	err = c.DB.Model(messages).Where("deleted = ?", false).Where("time >=?", timeStamp).Where(
		"EXISTS (SELECT 1 FROM unnest(usersId) AS elem WHERE elem = ?)", userId).Order("time DESC").Limit(1000).Select()
	return
}

func (m messageModel) DeleteMessage(id string, userId string, roomId string) (err error) {
	model := &message{
		Deleted: true,
	}
	_, err = c.DB.Model(model).Where("deleted = ?", false).Where("id=? and userId=? and roomId=?", id, userId, roomId).Update()
	return
}
