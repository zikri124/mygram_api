package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
)

type SocialMediaRepository interface {
	CreateSocial(ctx context.Context, social *model.SocialMedia) error
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
