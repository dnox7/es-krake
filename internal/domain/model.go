package domain

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BaseModelWithDeleted struct {
	BaseModel
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}
