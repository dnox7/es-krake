package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

type AccessRequirement struct {
	ID        int            `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Code      string         `gorm:"column:code;type:varchar(25);not null;unique" json:"code"`
	Functions []FunctionCode `gorm:"many2many:access_requirement_function_codes;joinForeignKey:access_requirement_id;joinReferences:function_code_id"`
	model.BaseModel
}
