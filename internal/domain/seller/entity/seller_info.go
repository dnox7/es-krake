package entity

const SellerInfoTableName = "seller_infos"

type SellerInfo struct {
}

func (SellerInfo) TableName() string {
	return SellerInfoTableName
}
