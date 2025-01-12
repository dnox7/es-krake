package utils

import "time"

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseModelWithDeleted struct {
	BaseModel
	DeletedAt interface{}
}
