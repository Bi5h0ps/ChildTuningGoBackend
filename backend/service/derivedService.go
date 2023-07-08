package service

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/repository"
)

type IDerivedService interface {
	SaveNewDerived(exDerievd *model.DerivedExercise) (int, error)
	SaveDerivedUpdate(exDerievd *model.DerivedExercise)
	GetDerivedById(id int) (*model.DerivedExercise, error)
	GetAllDerived(username string, favoriteId int) ([]model.DerivedExercise, error)
}

type DerivedService struct {
	DerivedRepository repository.IDerivedRepository
}

func (d *DerivedService) SaveNewDerived(exDerievd *model.DerivedExercise) (int, error) {
	return d.DerivedRepository.Insert(exDerievd)
}

func (d *DerivedService) SaveDerivedUpdate(exDerievd *model.DerivedExercise) {
	d.DerivedRepository.Save(exDerievd)
}

func (d *DerivedService) GetDerivedById(id int) (*model.DerivedExercise, error) {
	return d.DerivedRepository.Select(id)
}

func (d *DerivedService) GetAllDerived(username string, favoriteId int) ([]model.DerivedExercise, error) {
	return d.DerivedRepository.SelectAll(username, favoriteId)
}

func NewDerivedService(repo repository.IDerivedRepository) IDerivedService {
	return &DerivedService{DerivedRepository: repo}
}
