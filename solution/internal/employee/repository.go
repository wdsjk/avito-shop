package employee

import (
	"context"

	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type Repository interface {
	SaveEmployee(ctx context.Context, name, password string) (string, error)
	GetEmployee(ctx context.Context, name string) (*Employee, error)
	BuyItem(ctx context.Context, name, item string, shop shop.Shop, t transfer.Service) error
	TransferCoins(ctx context.Context, sender, receiver string, amount int, t transfer.Service) error
}
