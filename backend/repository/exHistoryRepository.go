package repository

import (
	"ChildTuningGoBackend/backend/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IExHistoryRepository interface {
	Conn() error
	Select(exerciseId string) (exHistory *model.ExerciseHistory, err error)
	SelectAll(userId string) (exHistories []model.ExerciseHistory, err error)
	Insert(exHistory *model.ExerciseHistory) (err error)
}

const defaultExerciseDB = "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/childTuningDB?parseTime=true"

//const defaultUserDB = "root:Password2023!@tcp(127.0.0.1:3306)/childTuningDB?charset=utf8"

type ExHistoryRepository struct {
	myGormConn *gorm.DB
}

func (u *ExHistoryRepository) Conn() (err error) {
	if u.myGormConn == nil {
		u.myGormConn, err = gorm.Open(mysql.Open(defaultExerciseDB), &gorm.Config{})
		if err != nil {
			return
		}
		err = u.myGormConn.AutoMigrate(&model.ExerciseHistory{})
		if err != nil {
			return
		}
	}
	return nil
}

func (u *ExHistoryRepository) Select(exerciseId string) (exHistory *model.ExerciseHistory, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	exHistory = &model.ExerciseHistory{}
	if result := u.myGormConn.Where("exerciseId", exerciseId).First(exHistory); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (u *ExHistoryRepository) SelectAll(username string) (exHistories []model.ExerciseHistory, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	var exerciseHistoryList []model.ExerciseHistory
	result := u.myGormConn.Where("username = ?", username).Find(&exerciseHistoryList)
	if result.Error != nil {
		return nil, result.Error
	}
	return exerciseHistoryList, nil
}

func (u *ExHistoryRepository) Insert(exHistory *model.ExerciseHistory) (err error) {
	if err = u.Conn(); err != nil {
		return
	}
	if result := u.myGormConn.Create(exHistory); result.Error != nil {
		return result.Error
	}
	return nil
}

func NewExHistoryRepository(db *gorm.DB) IExHistoryRepository {
	return &ExHistoryRepository{myGormConn: db}
}
