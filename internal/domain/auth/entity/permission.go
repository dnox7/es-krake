package entity

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

type Permission struct {
	ID         int       `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(50);not null" json:"name"`
	ActionID   int       `gorm:"column:action_id;type:smallint;not null" json:"action_id"`
	ResourceID int       `gorm:"column:resource_id;type:smallint;not null" json:"resource_id"`
	Action     *Action   `gorm:"foreignKey:ActionID"`
	Resource   *Resource `gorm:"foreignKey:ResourceID"`
	model.BaseModelWithDeleted
}

func (p Permission) Code() string {
	if p.Action == nil || p.Resource == nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", p.Action.Code, p.Resource.Code)
}
