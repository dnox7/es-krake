package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const OtpTableName = "otp"

type Otp struct {
	ID        string `gorm:"column:id;primaryKey;bigint;not null;autoIncrement" json:"id"`
	KcUserID  string `gorm:"column:kc_user_id;type:varchar(255);not null"       json:"kc_user_id"`
	Token     string `gorm:"column:token;type:varchar(255);not null"            json:"token"`
	ExpiredAt string `gorm:"column:expired_at;type:datetime;not null"           json:"expired_at"`
	model.BaseModel
}

func (Otp) TableName() string {
	return OtpTableName
}
