package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint32    `json:"id"`
	UserId    uint32    `json:"user_id"`
	PhotoId   uint32    `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type CreateComment struct {
	PhotoId uint32 `json:"photo_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type CreateCommentRes struct {
	ID        uint32    `json:"id"`
	UserId    uint32    `json:"user_id"`
	Message   string    `json:"message"`
	PhotoId   uint32    `json:"photo_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentView struct {
	ID        uint32    `json:"id"`
	UserId    uint32    `json:"user_id"`
	PhotoId   uint32    `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserItem  `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Photo     PhotoItem `json:"photo" gorm:"foreignKey:PhotoId;references:ID"`
}

func (c *Comment) BeforeCreate(db *gorm.DB) (err error) {
	if c.ID == 0 {
		c.ID = uuid.New().ID()
	}
	return
}
