package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
	"strings"
)

type IChat interface {
	Conn() error
	SelectByUsername(userName string) (chatHistory []model.ChatHistory, err error)
	SelectByQuestionId(questionId string) (pair []model.ChatHistory, err error)
	Insert(chatHistory *model.ChatHistory) (err error)
	Save(chatHistory *model.ChatHistory)
}

type ChatRepository struct{}

func (c *ChatRepository) Conn() (err error) {
	err = provider.DatabaseEngine.AutoMigrate(&model.ChatHistory{})
	return
}

func (c *ChatRepository) SelectByUsername(username string) (chatHistory []model.ChatHistory, err error) {
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

func (c *ChatRepository) SelectByQuestionId(questionId string) (pair []model.ChatHistory, err error) {
	if err = c.Conn(); err != nil {
		return
	}
	result := provider.DatabaseEngine.Where("question_id = ?", questionId).Find(&pair)
	if result.Error != nil {
		return nil, result.Error
	}
	return
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

func (c *ChatRepository) Save(chatHistory *model.ChatHistory) {
	s := strings.ReplaceAll(chatHistory.CreateTime, "T", " ")
	s = s[:len(s)-1]
	chatHistory.CreateTime = s
	provider.DatabaseEngine.Save(chatHistory)
}

func NewChatRepository() IChat {
	return &ChatRepository{}
}
