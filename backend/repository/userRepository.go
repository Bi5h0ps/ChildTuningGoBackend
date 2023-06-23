package repository

import (
	"ChildTuningGoBackend/backend/model"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IUser interface {
	Conn() error
	Select(userName string) (user *model.User, err error)
	Insert(user *model.User) (userId int64, err error)
}

//const defaultUserDB = "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/childTuningDB?charset=utf8"

const defaultUserDB = "root:Password2023!@tcp(127.0.0.1:3306)/childTuningDB?charset=utf8"

type UserRepository struct {
	myGormConn *gorm.DB
}

func (u *UserRepository) Conn() (err error) {
	if u.myGormConn == nil {
		u.myGormConn, err = gorm.Open(mysql.Open(defaultUserDB), &gorm.Config{})
		if err != nil {
			return
		}
		err = u.myGormConn.AutoMigrate(&model.User{})
		if err != nil {
			return
		}
	}
	return nil
}

func (u *UserRepository) Select(username string) (user *model.User, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	user = &model.User{}
	if result := u.myGormConn.Where("username", username).First(user); result.Error != nil {
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
	if result := u.myGormConn.Create(user); result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func NewUserRepository(db *gorm.DB) IUser {
	return &UserRepository{myGormConn: db}
}
