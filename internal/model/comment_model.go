package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint32    `json:"id"`
	UserId    uint32    `json:"user_id"`
	PhotoId   string    `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type CreateComment struct {
	UserId  uint32 `json:"user_id" validate:"required"`
	PhotoId string `json:"photo_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

func (c *Comment) BeforeCreate(db *gorm.DB) (err error) {
	if c.ID == 0 {
		c.ID = uuid.New().ID()
	}
	return
}
