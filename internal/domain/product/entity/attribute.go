package entity

import "pech/es-krake/internal/domain/shared/model"

type Attribute struct {
	ID              int            `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" json:"id"`
	Name            string         `gorm:"column:name;type:varchar(10);not null"                          json:"name"`
	Description     *string        `gorm:"column:description;type:text"                                   json:"description"`
	AttributeTypeID int            `gorm:"column:attribute_type_id;not null;type:smallint"                json:"attribute_type_id"`
	AttributeType   *AttributeType `gorm:"<-:false;foreignKey:AttributeTypeID"`
	model.BaseModel
}
