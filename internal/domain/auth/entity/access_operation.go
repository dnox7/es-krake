package entity

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

type AccessOperation struct {
	ID         int       `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	ActionID   int       `gorm:"column:action_id;type:smallint;not null" json:"action_id"`
	ResourceID int       `gorm:"column:resource_id;type:smallint;not null" json:"resource_id"`
	Action     *Action   `gorm:"foreignKey:ActionID"`
	Resource   *Resource `gorm:"foreignKey:ResourceID"`
	model.BaseModel
}

func (o AccessOperation) Code() (string, error) {
	if o.Action == nil || o.Resource == nil {
		return "", fmt.Errorf("failed to map access opeation to code, id: %d", o.ID)
	}
	return o.Action.Code + ":" + o.Resource.Code, nil
}

func MapOperationsToCodes(ops []AccessOperation) ([]string, error) {
	codes := make([]string, 0, len(ops))
	for _, op := range ops {
		c, err := op.Code()
		if err != nil {
			return nil, err
		}
		codes = append(codes, c)
	}
	return codes, nil
}
