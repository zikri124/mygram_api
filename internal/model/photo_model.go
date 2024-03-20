package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserId    uint32    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type PhotoCreate struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" validate:"required"`
}

type PhotoView struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserId    uint32    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	User      UserEdit
}

type PhotoRes struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserId    uint32    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePhoto struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" validate:"required"`
	UserId   uint32 `json:"user_id" validate:"required"`
}

func (p *Photo) BeforeCreate(db *gorm.DB) (err error) {
	if p.ID == 0 {
		p.ID = uuid.New().ID()
	}
	return
}
