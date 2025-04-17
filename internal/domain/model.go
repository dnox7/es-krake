package domain

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
}

type BaseModelWithDeleted struct {
	BaseModel
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:int(11);default:0"`
}
