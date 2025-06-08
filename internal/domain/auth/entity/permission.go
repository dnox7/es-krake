package entity

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
	"github.com/dpe27/es-krake/pkg/utils"
	"golang.org/x/sync/errgroup"
)

const PermissionTableName = "permissions"

type Permission struct {
	ID         int                   `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Name       string                `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	Operations []PermissionOperation `gorm:"foreignKey:PermissionID"`
	model.BaseModel
}

func (Permission) TableName() string {
	return PermissionTableName
}

func (p Permission) Codes() ([]string, error) {
	if len(p.Operations) == 0 {
		return nil, fmt.Errorf("failed to map permission to code: id=%d", p.ID)
	}

	codes := make([]string, 0, len(p.Operations))
	for _, po := range p.Operations {
		if po.AccessOperation == nil {
			return nil, fmt.Errorf("failed to map permissionOperation to code: id=%d", po.ID)
		}

		codes = append(codes, po.AccessOperation.Code)
	}

	return codes, nil
}

type PermissionSlice []Permission

func (perms PermissionSlice) MapPermissionsToCodes() ([]string, error) {
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

func (perms PermissionSlice) HasRequiredOperations(requiredOps []AccessOperation) (bool, error) {
	var (
		g         errgroup.Group
		permCodes []string
		opCodes   []string
	)

	g.Go(func() error {
		var err error
		permCodes, err = perms.MapPermissionsToCodes()
		return err
	})

	g.Go(func() error {
		var err error
		opCodes, err = MapOperationsToCodes(requiredOps)
		return err
	})

	if err := g.Wait(); err != nil {
		return false, err
	}

	return utils.IsSubSet(opCodes, permCodes), nil
}
