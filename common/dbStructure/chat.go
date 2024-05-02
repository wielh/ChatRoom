package dbstructure

import (
	c "common"
	"time"
)

type message struct {
	ID       string `pg:"id,pk"`
	RoomID   string
	UserID   string
	Username string
	Time     time.Time
	Deleted  bool
	Content  string
}

type messageModel struct {
}

var MessageModel = &messageModel{}

func (m messageModel) PushMessage(userId string, username string, roomId string, content string) (err error) {
	model := &message{
		UserID:   userId,
		Username: username,
		RoomID:   roomId,
		Content:  content,
		Deleted:  false,
	}
	_, err = c.DB.Model(model).Insert()
	return
}

func (m messageModel) GetMessages(userId string, timeStamp time.Time) ([]message, error) {
	var messages []message

	err := c.DB.Model(&messages).Where("deleted = ?", false).
		Where("user_id = ?", userId).Where("time >=? ", timeStamp).Order("time ASC").Limit(1000).Select()
	return messages, err
}

func (m messageModel) DeleteMessage(id string, userId string, roomId string) (err error) {
	model := message{
		Deleted: true,
	}
	_, err = c.DB.Model(&model).Where("deleted = ?", false).Where("id=? and userId=? and roomId=?", id, userId, roomId).Update()
	return
}
