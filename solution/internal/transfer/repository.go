package transfer

import "context"

type Repository interface {
	SaveTransfer(ctx context.Context, senderName, receiverName string, amount int) error
	GetTransfersByEmployee(ctx context.Context, name string) ([]*Transfer, error)
}
