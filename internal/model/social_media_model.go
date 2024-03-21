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

type NewSocialMedia struct {
	Name           string `json:"name" validate:"required"`
	SocialMediaUrl string `json:"social_media_url" validate:"required"`
}

type CreateSocialMediaRes struct {
	ID             uint32    `json:"id"`
	UserId         uint32    `json:"user_id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateSocialMediaRes struct {
	ID             uint32    `json:"id"`
	UserId         uint32    `json:"user_id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SocialMediaView struct {
	ID             uint32    `json:"id"`
	UserId         uint32    `json:"user_id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           UserItem  `json:"user" gorm:"foreignKey:UserId;references:ID"`
}

func (u *SocialMedia) BeforeCreate(db *gorm.DB) (err error) {
	if u.ID == 0 {
		u.ID = uuid.New().ID()
	}
	return
}
