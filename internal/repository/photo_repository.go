package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
)

type PhotoRepository interface {
	CreatePhoto(ctx context.Context, photo *model.Photo) error
}

type photoRepositoryImpl struct {
	db infrastructure.GormPostgres
}

func NewPhotoRepository(db infrastructure.GormPostgres) PhotoRepository {
	return &photoRepositoryImpl{db: db}
}

func (p *photoRepositoryImpl) CreatePhoto(ctx context.Context, photo *model.Photo) error {
	db := p.db.GetConnection()

	err := db.
		WithContext(ctx).
		Table("photos").
		Create(&photo).
		Error

	return err
}
