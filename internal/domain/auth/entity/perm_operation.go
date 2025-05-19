package entity

type PermissionOperation struct {
	ID               int              `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	PermisssionID    int              `gorm:"column:permission_id;type:bigint;not null" json:"permission_id"`
	AccessOpeationID int              `gorm:"column:access_operation_id;type:bigint;not null" json:"access_operation_id"`
	Permission       *Permission      `gorm:"foreignKey:PermissionID"`
	AccessOperation  *AccessOperation `gorm:"foreignKey:AccessOpeationID"`
}
