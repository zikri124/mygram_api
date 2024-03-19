package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialMedia struct {
	ID             uint32    `json:"id"`
	UserId         uint32    `json:"user_id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt
}

type CreateSocialMedia struct {
	UserId         uint32 `json:"user_id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	SocialMediaUrl string `json:"social_media_url" validate:"required"`
}

func (u *SocialMedia) BeforeCreate(db *gorm.DB) (err error) {
	if u.ID == 0 {
		u.ID = uuid.New().ID()
	}
	return
}
