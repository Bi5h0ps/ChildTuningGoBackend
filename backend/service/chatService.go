package service

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/repository"
)

type IChatService interface {
	GetChatHistoryByUsername(username string) ([]model.ChatHistory, error)
	GetChatHistoryByQuestionId(questionId string) ([]model.ChatHistory, error)
	WriteChatHistory(history *model.ChatHistory) (err error)
	FavoriteChatHistory(questionId string) (err error)
	UnFavoriteChatHistory(questionId string) (err error)
}

type ChatService struct {
	ChatRepo repository.IChat
}

func (c *ChatService) GetChatHistoryByQuestionId(questionId string) ([]model.ChatHistory, error) {
	return c.ChatRepo.SelectByQuestionId(questionId)
}

func (c *ChatService) FavoriteChatHistory(questionId string) error {
	pair, err := c.ChatRepo.SelectByQuestionId(questionId)
	if err != nil {
		return err
	}
	pair[0].IsFavorite = true
	pair[1].IsFavorite = true
	c.ChatRepo.Save(&pair[0])
	c.ChatRepo.Save(&pair[1])
	return nil
}

func (c *ChatService) UnFavoriteChatHistory(questionId string) (err error) {
	pair, err := c.ChatRepo.SelectByQuestionId(questionId)
	if err != nil {
		return err
	}
	pair[0].IsFavorite = false
	pair[1].IsFavorite = false
	c.ChatRepo.Save(&pair[0])
	c.ChatRepo.Save(&pair[1])
	return nil
}

func (c *ChatService) GetChatHistoryByUsername(username string) ([]model.ChatHistory, error) {
	return c.ChatRepo.SelectByUsername(username)
}

func (c *ChatService) WriteChatHistory(history *model.ChatHistory) (err error) {
	return c.ChatRepo.Insert(history)
}

func NewChatService(repo repository.IChat) IChatService {
	return &ChatService{ChatRepo: repo}
}
