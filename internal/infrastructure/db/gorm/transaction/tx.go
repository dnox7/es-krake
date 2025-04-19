package tx

import (
	tx "pech/es-krake/internal/domain/shared/transaction"

	"gorm.io/gorm"
)

type gormTx struct {
	tx *gorm.DB
}

func NewGormTx(tx *gorm.DB) tx.Base {
	return &gormTx{tx}
}

func (g *gormTx) GetTx() interface{} {
	return g.tx
}
