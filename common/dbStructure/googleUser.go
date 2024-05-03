package dbstructure

import (
	c "common"
	"time"
)

type GoogleUser struct {
	GoogleId       string `pg:"googleid,pk"`
	FirstName      string
	LastName       string
	Sex            string
	Email          []string `pg:",array"`
	CreateDatetime time.Time
	Birth          time.Time
}

func (u *GoogleUser) TableName() string {
	return "google_users"
}

type googleUserModel struct{}

var GoogleUserModel = &googleUserModel{}

func (model googleUserModel) InsertUser(googleId string, firstName string, lastName string, sex string, email []string, birth time.Time) error {
	product := &GoogleUser{
		GoogleId:  googleId,
		FirstName: firstName,
		LastName:  lastName,
		Sex:       sex,
		Email:     email,
		Birth:     birth,
	}
	_, err := c.DB.Model(product).OnConflict("(googleid) DO NOTHING").Insert()
	return err
}

func (model googleUserModel) SelectUser(googleId string) (*GoogleUser, error) {
	product := &GoogleUser{}
	err := c.DB.Model(product).OnConflict("(googleid) DO NOTHING").Where("googleid=?", googleId).First()
	if err != nil {
		c.WarnLogger("common", "googleUserModel.SelectUser", "no user", err, googleId)
		return nil, err
	}

	return product, nil
}
