package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/devpayments/core/models"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type PaymentRepository struct {
	*BaseRepository[models.PaymentModel, models.Payment]
}

func NewPaymentRepository(dbCon *sqlx.DB, tableName string) *PaymentRepository {
	baseRepo := NewBaseRepository[models.PaymentModel, models.Payment](dbCon, tableName)
	return &PaymentRepository{BaseRepository: baseRepo}
}

func (r *PaymentRepository) Update(ctx context.Context, p *models.Payment) (int64, error) {
	p.UpdatedAt = time.Now()

	var m models.PaymentModel
	model := m.FromEntity(*p).(models.PaymentModel)

	ds := goqu.Update(r.tableName).
		Set(model).
		Where(goqu.Ex{
			"id": model.ID,
		})

	updateSQL, _, err := ds.ToSQL()
	if err != nil {
		panic(err)
	}
	fmt.Println(updateSQL)

	res, err := r.db.ExecContext(ctx, updateSQL)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *PaymentRepository) GetPaymentTransaction(ctx context.Context, txnType string, paymentId uuid.UUID) (*models.Transaction, error) {
	var transactionField string
	if txnType == "destination" {
		transactionField = "destination_transaction_id"
	} else if txnType == "source" {
		transactionField = "source_transaction_id"
	} else {
		return nil, errors.New("invalid transaction type")
	}

	query, _, err := goqu.
		From(r.tableName).
		Join(
			goqu.T("transactions"),
			goqu.On(
				goqu.I(transactionField).Eq(goqu.I("transactions.id")),
			),
		).
		Select("transactions.*").
		Where(goqu.Ex{"payments.id": paymentId.String()}).
		ToSQL()

	row := r.db.QueryRowxContext(ctx, query)

	var m models.TransactionModel
	err = row.StructScan(&m)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	t := m.ToEntity()
	return &t, nil
}

func (r *PaymentRepository) SetPaymentTransaction(ctx context.Context, paymentId uuid.UUID, transaction *models.Transaction) error {
	var transactionField string
	if transaction.Type == "destination" {
		transactionField = "destination_transaction_id"
	} else if transaction.Type == "source" {
		transactionField = "source_transaction_id"
	} else {
		return errors.New("invalid transaction type")
	}

	var m models.TransactionModel
	model := m.FromEntity(*transaction)

	// Create New Transaction
	insertTransactionSQL, _, err := goqu.Insert("transactions").Rows(
		model,
	).ToSQL()
	if err != nil {
		panic(err)
	}
	fmt.Println(transaction)

	_, err = r.db.ExecContext(ctx, insertTransactionSQL)
	if err != nil {
		return err
	}

	// Update Payment Transaction
	updatePaymentSQL, _, err := goqu.
		Update("payments").
		Set(map[string]interface{}{
			transactionField: transaction.ID,
		}).
		Where(goqu.Ex{"id": paymentId.String()}).
		ToSQL()
	_, err = r.db.ExecContext(ctx, updatePaymentSQL)
	if err != nil {
		panic(err)
	}

	return nil
}
