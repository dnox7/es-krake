package entity

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

type Permission struct {
	ID         int                   `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Name       string                `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Operations []PermissionOperation `gorm:"foreignKey:PermissionID"`
	model.BaseModelWithDeleted
}

func (p Permission) Codes() ([]string, error) {
	if p.Operations == nil || len(p.Operations) == 0 {
		return nil, fmt.Errorf("failed to map permission to code: id=%d", p.ID)
	}

	codes := make([]string, 0, len(p.Operations))
	for _, po := range p.Operations {
		if po.AccessOperation == nil {
			return nil, fmt.Errorf("failed to map permissionOperation to code: id=%d", po.ID)
		}

		c, err := po.AccessOperation.Code()
		if err != nil {
			return nil, err
		}

		codes = append(codes, c)
	}

	return codes, nil
}

func MapPermissionsToCodes(perms []Permission) ([]string, error) {
	codes := []string{}
	for _, p := range perms {
		c, err := p.Codes()
		if err != nil {
			return nil, err
		}
		codes = append(codes, c...)
	}
	return codes, nil
}
