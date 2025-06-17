package entity

import (
	"github.com/dpe27/es-krake/internal/domain/shared/model"
)

const BuyerInfoTableName = "buyer_infos"

type BuyerInfo struct {
	ID             int     `gorm:"column:id;primaryKey;type:bigint;autoIncrement;not null" json:"id"`
	BuyerAccountID int     `gorm:"column:buyer_account_id;type:bigint;not null" json:"buyer_account_id"`
	FirstName      string  `gorm:"column:first_name;type:varchar(255);not null" json:"first_name"`
	LastName       string  `gorm:"column:last_name;type:varchar(255);not null" json:"last_name"`
	GenderID       int     `gorm:"column:gender_id;type:bigint;not null" json:"gender_id"`
	BirthdayYear   *int    `gorm:"column:birthday_year;type:int" json:"birthday_year"`
	BirthdayMonth  *int    `gorm:"column:birthday_month;type:int" json:"birthday_month"`
	BirthdayDay    *int    `gorm:"column:birthday_day;type:int" json:"birthday_day"`
	PhoneNumber    *string `gorm:"column:phone_number;type:varchar(255)" json:"phone_number"`
	PostalCode     *string `gorm:"column:postal_code;type:varchar(255)" json:"postal_code"`
	Address1       *string `gorm:"column:address1;type:varchar(255)" json:"address1"`
	Address2       *string `gorm:"column:address2;type:varchar(255)" json:"address2"`

	Gender *Gender `gorm:"foreignKey:GenderID" json:"gender"`
	model.BaseModelWithDeleted
}

func (BuyerInfo) TableName() string {
	return BuyerInfoTableName
}
