package dbstructure

import (
	c "common"
)

type GoogleUser struct {
	GoogleId  string   `pg:"type:pk,column:googleId"`
	FirstName string   `pg:"column:firstName"`
	LastName  string   `pg:"column:lastName"`
	Sex       string   `pg:"column:sex"`
	Email     []string `pg:"column:email"`
	Age       int32    `pg:"column:age"`
}

type googleUserModel struct{}

var GoogleUserModel = &googleUserModel{}

func (model googleUserModel) InsertUser(googleId string, firstName string, lastName string, sex string, email []string, age int32) error {
	product := &GoogleUser{
		GoogleId:  googleId,
		FirstName: firstName,
		LastName:  lastName,
		Sex:       sex,
		Email:     email,
		Age:       age,
	}
	_, err := c.DB.Model(product).OnConflict("(googleId) DO NOTHING").Insert()
	return err
}

func (model googleUserModel) SelectUser(googleId string) (*GoogleUser, error) {
	product := &GoogleUser{}
	err := c.DB.Model(product).OnConflict("(googleId) DO NOTHING").Where("googleId=?", googleId).First()
	if err != nil {
		return nil, err
	}

	return product, nil
}
