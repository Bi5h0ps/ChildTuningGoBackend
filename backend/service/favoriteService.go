package service

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/repository"
)

type IFavoriteService interface {
	FavoriteAsking(favorite *model.Favorite) error
	FavoriteExercise(favorite *model.Favorite) error
	RemoveAskingFavorite(questionId string) error
	RemoveExerciseFavorite(exerciseId string) error
	GetFavoriteList(username string) (result []model.Favorite, err error)
	GetFavoriteExercise(id int) (*model.Favorite, error)
	FavoriteUpdate(favroite *model.Favorite)
}

type FavoriteService struct {
	FavoriteRepository repository.IFavoriteRepository
}

func (f *FavoriteService) FavoriteAsking(favorite *model.Favorite) error {
	err := f.FavoriteRepository.MarkFavorite(false, favorite)
	return err
}

func (f *FavoriteService) FavoriteExercise(favorite *model.Favorite) error {
	err := f.FavoriteRepository.MarkFavorite(true, favorite)
	return err
}

func (f *FavoriteService) RemoveAskingFavorite(questionId string) error {
	err := f.FavoriteRepository.DeleteById(false, questionId)
	return err
}

func (f *FavoriteService) RemoveExerciseFavorite(exerciseId string) error {
	err := f.FavoriteRepository.DeleteById(true, exerciseId)
	return err
}

func (f *FavoriteService) GetFavoriteList(username string) (result []model.Favorite, err error) {
	return f.FavoriteRepository.SelectAll(username)
}

func (f *FavoriteService) GetFavoriteExercise(id int) (*model.Favorite, error) {
	return f.FavoriteRepository.SelectById(id)
}

func (f *FavoriteService) FavoriteUpdate(favroite *model.Favorite) {
	f.FavoriteRepository.SaveFavoriteUpdate(favroite)
}

func NewFavoriteService(favoriteRepo repository.IFavoriteRepository) IFavoriteService {
	return &FavoriteService{FavoriteRepository: favoriteRepo}
}
