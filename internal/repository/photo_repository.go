package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
	"gorm.io/gorm"
)

type PhotoRepository interface {
	CreatePhoto(ctx context.Context, photo *model.Photo) error
	GetAllPhotos(ctx context.Context) ([]model.PhotoView, error)
	GetPhotoById(ctx context.Context, photoId uint32) (*model.Photo, error)
	UpdatePhoto(ctx context.Context, photo *model.Photo) error
	DeletePhoto(ctx context.Context, photoId uint32) error
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

func (p *photoRepositoryImpl) GetAllPhotos(ctx context.Context) ([]model.PhotoView, error) {
	db := p.db.GetConnection()
	photos := []model.PhotoView{}

	err := db.
		WithContext(ctx).
		Table("photos").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email, username").Table("users").Where("deleted_at is null")
		}).
		Find(&photos).
		Error

	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (p *photoRepositoryImpl) GetPhotoById(ctx context.Context, photoId uint32) (*model.Photo, error) {
	db := p.db.GetConnection()
	photo := model.Photo{}

	err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", photoId).
		Where("deleted_at IS NULL").
		Find(&photo).
		Error

	if err != nil {
		return nil, err
	}

	return &photo, nil
}

func (p *photoRepositoryImpl) UpdatePhoto(ctx context.Context, photo *model.Photo) error {
	db := p.db.GetConnection()
	err := db.
		WithContext(ctx).
		Updates(&photo).
		Error

	return err
}

func (p *photoRepositoryImpl) DeletePhoto(ctx context.Context, photoId uint32) error {
	db := p.db.GetConnection()
	photo := model.Photo{ID: photoId}

	err := db.
		WithContext(ctx).
		Model(&photo).
		Delete(&photo).
		Error

	return err
}
