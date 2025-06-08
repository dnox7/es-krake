package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const AccessRequirementTableName = "access_requirements"

type AccessRequirement struct {
	ID         int                          `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Code       string                       `gorm:"column:code;type:varchar(25);not null;unique"            json:"code"`
	Operations []AccessRequirementOperation `gorm:"foreignKey:AccessRequirementID"`
	model.BaseModel
}

func (AccessRequirement) TableName() string {
	return AccessOperationsTableName
}
