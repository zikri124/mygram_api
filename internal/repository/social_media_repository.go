package repository

import (
	"context"
	"log"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
	"gorm.io/gorm"
)

type SocialMediaRepository interface {
	CreateSocial(ctx context.Context, social *model.SocialMedia) error
	GetAllSocials(ctx context.Context) ([]model.SocialMediaView, error)
	GetSocialById(ctx context.Context, socialId uint32) (*model.SocialMedia, error)
	UpdateSocial(ctx context.Context, social *model.SocialMedia) error
}

type socialMediaRepositoryImpl struct {
	db infrastructure.GormPostgres
}

func NewSocialMediaRepository(db infrastructure.GormPostgres) SocialMediaRepository {
	return &socialMediaRepositoryImpl{db: db}
}

func (s *socialMediaRepositoryImpl) CreateSocial(ctx context.Context, social *model.SocialMedia) error {
	db := s.db.GetConnection()

	err := db.
		WithContext(ctx).
		Table("social_medias").
		Create(&social).
		Error

	return err
}

func (s *socialMediaRepositoryImpl) GetAllSocials(ctx context.Context) ([]model.SocialMediaView, error) {
	db := s.db.GetConnection()
	socials := []model.SocialMediaView{}

	err := db.
		WithContext(ctx).
		Table("social_medias").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email, username").Table("users").Where("deleted_at is null")
		}).
		Find(&socials).
		Error

	if err != nil {
		return nil, err
	}

	return socials, nil
}

func (s *socialMediaRepositoryImpl) GetSocialById(ctx context.Context, socialId uint32) (*model.SocialMedia, error) {
	db := s.db.GetConnection()
	social := model.SocialMedia{}

	err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("id = ?", socialId).
		Where("deleted_at IS NULL").
		Find(&social).
		Error

	if err != nil {
		return nil, err
	}

	return &social, nil
}

func (s *socialMediaRepositoryImpl) UpdateSocial(ctx context.Context, social *model.SocialMedia) error {
	log.Println(social)

	db := s.db.GetConnection()
	err := db.
		WithContext(ctx).
		Table("social_medias").
		Updates(&social).
		Error

	return err
}
