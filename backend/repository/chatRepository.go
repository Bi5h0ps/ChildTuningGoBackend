package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
)

type IChat interface {
	Conn() error
	Select(userName string) (chatHistory []model.ChatHistory, err error)
	Insert(chatHistory *model.ChatHistory) (err error)
}

type ChatRepository struct{}

func (c *ChatRepository) Conn() (err error) {
	err = provider.DatabaseEngine.AutoMigrate(&model.ChatHistory{})
	return
}

func (c *ChatRepository) Select(username string) (chatHistory []model.ChatHistory, err error) {
	if err = c.Conn(); err != nil {
		return
	}
	var chatHistoryList []model.ChatHistory
	// Retrieve records from the "users" table with a specific username
	result := provider.DatabaseEngine.Where("username = ?", username).Order("create_time ASC").Find(&chatHistoryList)
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
	if result := provider.DatabaseEngine.Create(history); result.Error != nil {
		return result.Error
	}
	return nil
}

func NewChatRepository() IChat {
	return &ChatRepository{}
}
