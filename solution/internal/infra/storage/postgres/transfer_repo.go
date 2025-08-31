package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wdsjk/avito-shop/internal/transfer"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) SaveTransfer(ctx context.Context, senderName, receiverName string, amount int) error {
	const op = "infra.storage.postgres.SaveTransfer"

	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO transfers (sender_name, receiver_name, amount) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, senderName, receiverName, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *TransferRepository) GetTransfersByEmployee(ctx context.Context, name string) ([]*transfer.Transfer, error) {
	const op = "infra.storage.postgres.GetTransfersByEmployee"

	stmt, err := r.db.PrepareContext(ctx, "SELECT id, sender_name, receiver_name, amount FROM transfers WHERE sender_name=$1 OR receiver_name=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var transfers []*transfer.Transfer
	rows, err := stmt.QueryContext(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var t transfer.Transfer
		err := rows.Scan(&t.ID, &t.SenderName, &t.ReceiverName, &t.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transfers = append(transfers, &t)
	}

	return transfers, nil
}
