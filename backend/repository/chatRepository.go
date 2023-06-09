package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
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
	pair[0].CreateTime = provider.FixTimeFormat(pair[0].CreateTime)
	pair[1].CreateTime = provider.FixTimeFormat(pair[1].CreateTime)
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
	chatHistory.CreateTime = provider.FixTimeFormat(chatHistory.CreateTime)
	provider.DatabaseEngine.Save(chatHistory)
}

func NewChatRepository() IChat {
	return &ChatRepository{}
}
