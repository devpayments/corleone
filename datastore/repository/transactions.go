package repository

import (
	"github.com/devpayments/core/models"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	*BaseRepository[models.TransactionModel, models.Transaction]
}

func NewTransactionRepository(dbCon *sqlx.DB, tableName string) *TransactionRepository {
	baseRepo := NewBaseRepository[models.TransactionModel, models.Transaction](dbCon, tableName)
	return &TransactionRepository{BaseRepository: baseRepo}
}
