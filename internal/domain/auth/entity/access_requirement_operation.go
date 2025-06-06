package entity

const AccessRequirementOperationTableName = "access_requirement_operations"

type AccessRequirementOperation struct {
	ID                  int                `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	AccessRequirementID int                `gorm:"column:access_requirement_id;type:bigint;not null"       json:"access_requirement_id"`
	AccessOpeationID    int                `gorm:"column:access_operation_id;type:bigint;not null"         json:"access_operation_id"`
	AccessRequirement   *AccessRequirement `gorm:"foreignKey:AccessRequirementID"`
	AccessOperation     *AccessOperation   `gorm:"foreignKey:AccessOpeationID"`
}

func (AccessRequirementOperation) TableName() string {
	return AccessRequirementOperationTableName
}
