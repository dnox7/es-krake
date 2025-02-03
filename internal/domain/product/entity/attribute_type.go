package entity

import "pech/es-krake/internal/domain"

type AttributeType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	domain.BaseModel
}
