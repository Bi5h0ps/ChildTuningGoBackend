package service

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/repository"
)

type IChatService interface {
	GetChatHistoryByUsername(username string) ([]model.ChatHistory, error)
	WriteChatHistory(history *model.ChatHistory) (err error)
}

type ChatService struct {
	ChatRepo repository.IChat
}

func (c *ChatService) GetChatHistoryByUsername(username string) ([]model.ChatHistory, error) {
	return c.ChatRepo.Select(username)
}

func (c *ChatService) WriteChatHistory(history *model.ChatHistory) (err error) {
	return c.ChatRepo.Insert(history)
}

func NewChatService(repo repository.IChat) IChatService {
	return &ChatService{ChatRepo: repo}
}
