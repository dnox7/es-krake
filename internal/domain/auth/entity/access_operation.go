package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const AccessOperationsTableName = "access_operations"

type AccessOperation struct {
	ID          int     `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Name        string  `gorm:"column:name;type:varchar(50);not null"                   json:"name"`
	Code        string  `gorm:"column:code;type:varchar(50);not null;unique"            json:"code"`
	Description *string `gorm:"column:description;type:varchar(255)"                    json:"description"`
	model.BaseModel
}

func (AccessOperation) TableName() string {
	return AccessOperationsTableName
}

func MapOperationsToCodes(ops []AccessOperation) ([]string, error) {
	codes := make([]string, 0, len(ops))
	for _, op := range ops {
		codes = append(codes, op.Code)
	}
	return codes, nil
}
