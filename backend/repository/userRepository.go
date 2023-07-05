package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
	"errors"
)

type IUser interface {
	Conn() error
	Select(userName string) (user *model.User, err error)
	Insert(user *model.User) (userId int64, err error)
}

type UserRepository struct{}

func (u *UserRepository) Conn() (err error) {
	err = provider.DatabaseEngine.AutoMigrate(&model.User{})
	return
}

func (u *UserRepository) Select(username string) (user *model.User, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	user = &model.User{}
	if result := provider.DatabaseEngine.Where("username", username).First(user); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (u *UserRepository) Insert(user *model.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	//not allowed to set user id
	user.ID = 0
	//check if user already exists
	checkUser, _ := u.Select(user.Username)
	if checkUser != nil {
		//user already exist
		return 0, errors.New("Username already exists!")
	}
	if result := provider.DatabaseEngine.Create(user); result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func NewUserRepository() IUser {
	return &UserRepository{}
}
