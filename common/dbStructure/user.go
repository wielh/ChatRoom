package dbstructure

import (
	c "common"
	"time"
)

type GoogleUser struct {
	GoogleId       string `pg:"googleid,pk"`
	FirstName      string
	LastName       string
	Email          string
	CreateDatetime time.Time
}

func (u *GoogleUser) TableName() string {
	return "google_users"
}

type googleUserModel struct{}

var GoogleUserModel = &googleUserModel{}

func (model googleUserModel) InsertUser(googleId string, firstName string, lastName string, email string) error {
	product := &GoogleUser{
		GoogleId:  googleId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	_, err := c.DB.Model(product).OnConflict("(googleid) DO NOTHING").Insert()
	return err
}

func (model googleUserModel) SelectUser(googleId string) (*GoogleUser, error) {
	product := &GoogleUser{}
	err := c.DB.Model(product).OnConflict("(googleid) DO NOTHING").Where("googleid=?", googleId).First()
	if err != nil {
		c.WarnLogger("common", "googleUserModel.SelectUser", "db error", err, googleId)
		return nil, err
	}

	return product, nil
}

func (model googleUserModel) UserExist(googleId string) (bool, error) {
	product := &GoogleUser{}
	exist, err := c.DB.Model(product).Where("googleid=?", googleId).Exists()
	if err != nil {
		c.WarnLogger("common", "googleUserModel.SelectUser", "db error", err, googleId)
		return false, err
	}

	return exist, nil
}

type userModel struct{}

var UserModel = &userModel{}

func (model userModel) UserExist(userId string) (bool, error) {
	id, userType := c.ToSpecificID(userId)

	switch userType {
	case c.GoogleUser:
		exist, err := GoogleUserModel.UserExist(id)
		return exist, err
	default:
		return false, nil
	}
}
