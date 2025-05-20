package entity

type Enterprise struct {
	ID int `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
}
