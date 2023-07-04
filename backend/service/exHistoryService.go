package service

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/repository"
)

type IExHistoryService interface {
	SaveExHistory(history *model.ExerciseHistory) error
	GetExeHistoriesByUsername(username string) ([]model.ExerciseHistory, error)
	GetExHistoryById(exerciseId string) (*model.ExerciseHistory, error)
}

type ExHistoryService struct {
	ExerciseRepository repository.IExHistoryRepository
}

func (e *ExHistoryService) SaveExHistory(history *model.ExerciseHistory) error {
	return e.ExerciseRepository.Insert(history)
}

func (e *ExHistoryService) GetExeHistoriesByUsername(username string) (histories []model.ExerciseHistory, err error) {
	return e.ExerciseRepository.SelectAll(username)
}

func (e *ExHistoryService) GetExHistoryById(exerciseId string) (*model.ExerciseHistory, error) {
	return e.ExerciseRepository.Select(exerciseId)
}

func NewExHistoryService(exHistoryRepo repository.IExHistoryRepository) IExHistoryService {
	return &ExHistoryService{ExerciseRepository: exHistoryRepo}
}
