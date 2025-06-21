package entity

import "github.com/dpe27/es-krake/internal/domain/shared/model"

const EnterpriseTableName = "enterprises"

type Enterprise struct {
	ID            int     `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	Name          string  `gorm:"column:name;type:varchar(255);not null"                  json:"name"`
	KcRealmName   string  `gorm:"column:kc_realm_name;type:varchar(255);not null"         json:"kc_realm_name"`
	KcClientID    string  `gorm:"column:kc_client_id;type:varchar(255);not null"          json:"kc_client_id"`
	IsActive      bool    `gorm:"column:is_active;type:boolean;not null;default:true"     json:"is_active"`
	Address       string  `gorm:"column:address;type:varchar(255);not null"               json:"address"`
	Phone         string  `gorm:"column:phone;type:varchar(25);not null"                  json:"phone"`
	Fax           string  `gorm:"column:fax;type:varchar(50);not null"                    json:"fax"`
	MailAddress   string  `gorm:"column:mail_address;type:varchar(255);not null"          json:"mail_address"`
	MailSignature string  `gorm:"column:mail_signature;type:varchar(255);not null"        json:"mail_signature"`
	ThumbnailPath *string `gorm:"column:thumbnail_path;type:varchar(255)"                 json:"thumbnail_path"`
	WebsiteURL    *string `gorm:"column:website_url;type:varchar(255)"                    json:"website_url"`
	Notes         string  `gorm:"column:notes;type:varchar(255)"                          json:"notes"`
	model.BaseModel
}

func (Enterprise) TableName() string {
	return EnterpriseTableName
}
