package entity

const SellerAccountTableName = "seller_accounts"

type SellerAccount struct {
}

func (SellerAccount) TableName() string {
	return SellerAccountTableName
}
