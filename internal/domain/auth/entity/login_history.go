package entity

import (
	"time"

	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const LoginHistoryTableName = "login_histories"

type LoginHistory struct {
	ID              int       `gorm:"column:id;primaryKey;bigint;not null;autoIncrement" json:"id"`
	KcUserID        string    `gorm:"column:kc_user_id;type:varchar(255);not null" json:"kc_user_id"`
	LoggedInAt      time.Time `gorm:"column:logged_in_at;type:timestamp;autoCreateTime" json:"logged_in_at"`
	LoggedDevice    string    `gorm:"column:logged_device;type:varchar(255);not null" json:"logged_device"`
	LoggedIPAddress string    `gorm:"column:logged_ip_address;type:varchar(255);not null" json:"logged_ip_address"`
	model.BaseModel
}

func (LoginHistory) TableName() string {
	return LoginHistoryTableName
}
