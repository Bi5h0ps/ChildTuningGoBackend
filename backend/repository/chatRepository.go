package repository

import (
	"ChildTuningGoBackend/backend/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IChat interface {
	Conn() error
	Select(userName string) (chatHistory []model.ChatHistory, err error)
	Insert(chatHistory *model.ChatHistory) (err error)
}

const defaultChatHistoryDB = "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/childTuningDB?parseTime=true"

//const defaultChatHistoryDB = "root:Password2023!@tcp(127.0.0.1:3306)/childTuningDB?parseTime=true"

type ChatRepository struct {
	myGormConn *gorm.DB
}

func (c *ChatRepository) Conn() (err error) {
	if c.myGormConn == nil {
		c.myGormConn, err = gorm.Open(mysql.Open(defaultChatHistoryDB), &gorm.Config{})
		if err != nil {
			return
		}
		err = c.myGormConn.AutoMigrate(&model.ChatHistory{})
		if err != nil {
			return
		}
	}
	return nil
}

func (c *ChatRepository) Select(username string) (chatHistory []model.ChatHistory, err error) {
	if err = c.Conn(); err != nil {
		return
	}
	var chatHistoryList []model.ChatHistory
	// Retrieve records from the "users" table with a specific username
	result := c.myGormConn.Where("username = ?", username).Order("create_time ASC").Find(&chatHistoryList)
	if result.Error != nil {
		return nil, result.Error
	}
	return chatHistoryList, nil
}

func (c *ChatRepository) Insert(history *model.ChatHistory) (err error) {
	if err = c.Conn(); err != nil {
		return
	}
	history.ID = 0
	if result := c.myGormConn.Create(history); result.Error != nil {
		return result.Error
	}
	return nil
}

func NewChatRepository(db *gorm.DB) IChat {
	return &ChatRepository{myGormConn: db}
}
