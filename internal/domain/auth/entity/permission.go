package entity

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

type Permission struct {
	ID         int         `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	ParentID   int         `gorm:"column:parent_id;type:bigint" json:"parent_id"`
	Name       string      `gorm:"column:name;type:varchar(50);not null" json:"name"`
	ActionID   int         `gorm:"column:action_id;type:smallint;not null" json:"action_id"`
	ResourceID int         `gorm:"column:resource_id;type:smallint;not null" json:"resource_id"`
	Action     *Action     `gorm:"foreignKey:ActionID"`
	Resource   *Resource   `gorm:"foreignKey:ResourceID"`
	ParentPerm *Permission `gorm:"foreignKey:ParentID"`
	model.BaseModelWithDeleted
}

func (p Permission) Code() (string, error) {
	if p.Action == nil || p.Resource == nil {
		return "", fmt.Errorf("failed to map permission to code, id: %d", p.ID)
	}
	return p.Action.Code + ":" + p.Resource.Code, nil
}

func MapPermissionToCodes(perms []Permission) ([]string, error) {
	codes := make([]string, 0, len(perms))
	for _, p := range perms {
		c, err := p.Code()
		if err != nil {
			return nil, err
		}
		codes = append(codes, c)
	}
	return codes, nil
}
