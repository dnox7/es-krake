package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type AccessOperation struct {
	ID         int       `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	ActionID   int       `gorm:"column:action_id;type:smallint;not null" json:"action_id"`
	ResourceID int       `gorm:"column:resource_id;type:smallint;not null" json:"resource_id"`
	Action     *Action   `gorm:"foreignKey:ActionID"`
	Resource   *Resource `gorm:"foreignKey:ResourceID"`
	model.BaseModel
}

func (f AccessOperation) Code() string {
	if f.Action == nil || f.Resource == nil {
		return ""
	}
	return f.Action.Code + ":" + f.Resource.Code
}
