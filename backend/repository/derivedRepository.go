package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
)

type IDerivedRepository interface {
	Conn() error
	Select(id int) (exDerived *model.DerivedExercise, err error)
	SelectAll(username string, favoriteId int) (resultList []model.DerivedExercise, err error)
	Insert(exDerived *model.DerivedExercise) (id int, err error)
	Save(exDerived *model.DerivedExercise)
}

type DerivedRepository struct{}

func (d *DerivedRepository) Conn() (err error) {
	err = provider.DatabaseEngine.AutoMigrate(&model.DerivedExercise{})
	return
}

func (d *DerivedRepository) Select(id int) (exDerived *model.DerivedExercise, err error) {
	if err = d.Conn(); err != nil {
		return
	}
	exDerived = &model.DerivedExercise{}
	if result := provider.DatabaseEngine.Where("ID", id).First(exDerived); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (d *DerivedRepository) SelectAll(username string, favoriteId int) (resultList []model.DerivedExercise, err error) {
	if err = d.Conn(); err != nil {
		return
	}
	result := provider.DatabaseEngine.Where(map[string]interface{}{"username": username,
		"favorite_id": favoriteId}).Order("create_time ASC").Find(&resultList)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}

func (d *DerivedRepository) Insert(exDerived *model.DerivedExercise) (id int, err error) {
	if err = d.Conn(); err != nil {
		return
	}
	if result := provider.DatabaseEngine.Create(exDerived); result.Error != nil {
		return 0, result.Error
	}
	return exDerived.ID, nil
}

func (d *DerivedRepository) Save(exDerived *model.DerivedExercise) {
	exDerived.CreateTime = provider.FixTimeFormat(exDerived.CreateTime)
	provider.DatabaseEngine.Save(&exDerived)
}

func NewDerivedRepository() IDerivedRepository {
	return &DerivedRepository{}
}
