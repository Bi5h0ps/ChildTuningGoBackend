package repository

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
)

type IFavoriteRepository interface {
	Conn() (err error)
	Select(isNormal bool, id string) (result *model.Favorite, err error)
	Insert(favorite *model.Favorite) (err error)
	DeleteById(isNormal bool, id string) (err error)
	MarkFavorite(isNormal bool, body *model.Favorite) (err error)
	SelectAll(username string) (result []model.Favorite, err error)
	SelectById(id int) (result *model.Favorite, err error)
	SaveFavoriteUpdate(favorite *model.Favorite)
}

type FavoriteRepository struct{}

func (f *FavoriteRepository) Conn() (err error) {
	err = provider.DatabaseEngine.AutoMigrate(&model.Favorite{})
	return
}

func (f *FavoriteRepository) SaveFavoriteUpdate(favorite *model.Favorite) {
	favorite.CreateTime = provider.FixTimeFormat(favorite.CreateTime)
	provider.DatabaseEngine.Save(favorite)
}

func (f *FavoriteRepository) SelectById(id int) (result *model.Favorite, err error) {
	selection := provider.DatabaseEngine.Where(map[string]interface{}{"ID": id}).First(&result)
	if selection.Error != nil {
		return nil, selection.Error
	}
	return
}

func (f *FavoriteRepository) Select(isNormal bool, id string) (result *model.Favorite, err error) {
	tag := "normal"
	if !isNormal {
		tag = "asking"
	}
	var tuple model.Favorite
	selection := provider.DatabaseEngine.Where(map[string]interface{}{"origin": tag, "origin_id": id}).First(&tuple)
	if selection.Error != nil {
		return nil, selection.Error
	}
	return &tuple, nil
}

func (f *FavoriteRepository) SelectAll(username string) (result []model.Favorite, err error) {
	selection := provider.DatabaseEngine.Where(map[string]interface{}{"username": username, "is_deleted": false}).Order("create_time DESC").Find(&result)
	if selection.Error != nil {
		return nil, selection.Error
	}
	for i, _ := range result {
		result[i].CreateTime = provider.FixTimeFormat(result[i].CreateTime)
	}
	return
}

func (f *FavoriteRepository) Insert(favorite *model.Favorite) (err error) {
	if err = f.Conn(); err != nil {
		return
	}
	if result := provider.DatabaseEngine.Create(favorite); result.Error != nil {
		return result.Error
	}
	return nil
}

func (f *FavoriteRepository) DeleteById(isNormal bool, id string) (err error) {
	if err = f.Conn(); err != nil {
		return
	}
	tuple, selectionError := f.Select(isNormal, id)
	if selectionError != nil {
		return selectionError
	}
	tuple.IsDeleted = true
	tuple.CreateTime = provider.FixTimeFormat(tuple.CreateTime)
	writeBack := provider.DatabaseEngine.Save(&tuple)
	if writeBack.Error != nil {
		return writeBack.Error
	}
	return
}

func (f *FavoriteRepository) MarkFavorite(isNormal bool, body *model.Favorite) (err error) {
	if err = f.Conn(); err != nil {
		return
	}
	tuple, selectionError := f.Select(isNormal, body.OriginId)
	if selectionError != nil {
		insertError := f.Insert(body)
		if insertError != nil {
			return insertError
		}
	} else {
		tuple.IsDeleted = false
		tuple.CreateTime = provider.FixTimeFormat(tuple.CreateTime)
		writeBack := provider.DatabaseEngine.Save(&tuple)
		if writeBack.Error != nil {
			return writeBack.Error
		}
	}
	return
}

func NewFavoriteRepository() IFavoriteRepository {
	return &FavoriteRepository{}
}
