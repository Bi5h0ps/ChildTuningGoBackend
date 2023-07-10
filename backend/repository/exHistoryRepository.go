package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
)

type IExHistoryRepository interface {
	Conn() error
	Select(exerciseId string) (exHistory *model.ExerciseHistory, err error)
	SelectAll(userId string) (exHistories []model.ExerciseHistory, err error)
	Insert(exHistory *model.ExerciseHistory) (err error)
	Save(exHistory *model.ExerciseHistory)
}

type ExHistoryRepository struct{}

func (u *ExHistoryRepository) Conn() (err error) {
	err = provider.DatabaseEngine.AutoMigrate(&model.ExerciseHistory{})
	return
}

func (u *ExHistoryRepository) Select(exerciseId string) (exHistory *model.ExerciseHistory, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	exHistory = &model.ExerciseHistory{}
	if result := provider.DatabaseEngine.Where("exercise_id", exerciseId).First(exHistory); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (u *ExHistoryRepository) SelectAll(username string) (exHistories []model.ExerciseHistory, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	var exerciseHistoryList []model.ExerciseHistory
	result := provider.DatabaseEngine.Where("username = ?", username).Order("create_time DESC").Find(&exerciseHistoryList)
	if result.Error != nil {
		return nil, result.Error
	}
	return exerciseHistoryList, nil
}

func (u *ExHistoryRepository) Insert(exHistory *model.ExerciseHistory) (err error) {
	if err = u.Conn(); err != nil {
		return
	}
	exHistory.ID = 0
	if result := provider.DatabaseEngine.Create(exHistory); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *ExHistoryRepository) Save(exHistory *model.ExerciseHistory) {
	exHistory.CreateTime = provider.FixTimeFormat(exHistory.CreateTime)
	provider.DatabaseEngine.Save(exHistory)
}

func NewExHistoryRepository() IExHistoryRepository {
	return &ExHistoryRepository{}
}
