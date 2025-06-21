package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const BuyerAccountTableName = "buyer_accounts"

type BuyerAccount struct {
	ID           int     `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Nickname     *string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	KcUserID     string  `gorm:"column:kc_user_id;type:varchar(255);not null;unique" json:"kc_user_id"`
	LoginEnabled bool    `gorm:"column:login_enabled;type:boolean;not null;default:true" json:"login_enabled"`
	MailAddress  string  `gorm:"column:mail_address;type:varchar(255);not null" json:"mail_address"`
	MailVerified bool    `gorm:"column:mail_verified;type:boolean;not null;default:false" json:"mail_verified"`

	BuyerInfo *BuyerInfo `gorm:"foreignKey:BuyerAccountID" json:"buyer_info"`
	model.BaseModelWithDeleted
}
