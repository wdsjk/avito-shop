package employee

import (
	"context"

	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type Service interface {
	SaveEmployee(ctx context.Context, name string, password string) (string, error)
	GetEmployee(ctx context.Context, name string) (*EmployeeDto, error)
	BuyItem(ctx context.Context, name, item string, shop shop.Shop, t transfer.Service) error
	TransferCoins(ctx context.Context, sender, receiver string, amount int, t transfer.Service) error
}

type EmployeeService struct {
	repo Repository
}

func NewEmployeeService(repo Repository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) SaveEmployee(ctx context.Context, name string, password string) (string, error) {
	return s.repo.SaveEmployee(ctx, name, password)
}

func (s *EmployeeService) GetEmployee(ctx context.Context, name string) (*EmployeeDto, error) {
	employee, err := s.repo.GetEmployee(ctx, name)
	if err != nil {
		return nil, err
	}

	return ToDto(employee), nil
}

func (s *EmployeeService) BuyItem(ctx context.Context, name, item string, shop shop.Shop, t transfer.Service) error {
	return s.repo.BuyItem(ctx, name, item, shop, t)
}

func (s *EmployeeService) TransferCoins(ctx context.Context, sender, receiver string, amount int, t transfer.Service) error {
	return s.repo.TransferCoins(ctx, sender, receiver, amount, t)
}
