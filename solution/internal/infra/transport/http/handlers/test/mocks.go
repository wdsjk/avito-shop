package test

import (
	"context"

	"github.com/wdsjk/avito-shop/internal/employee"
	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type MockEmployeeService struct {
	SaveEmployeeFn  func(ctx context.Context, name string, password string) (string, error)
	GetEmployeeFn   func(ctx context.Context, name string) (*employee.EmployeeDto, error)
	BuyItemFn       func(ctx context.Context, name, item string, shop shop.Shop, t transfer.Service) error
	TransferCoinsFn func(ctx context.Context, sender, receiver string, amount int, t transfer.Service) error
}

func (m *MockEmployeeService) SaveEmployee(ctx context.Context, name string, password string) (string, error) {
	if m.SaveEmployeeFn != nil {
		return m.SaveEmployeeFn(ctx, name, password)
	}
	return "", nil
}

func (m *MockEmployeeService) GetEmployee(ctx context.Context, name string) (*employee.EmployeeDto, error) {
	if m.GetEmployeeFn != nil {
		return m.GetEmployeeFn(ctx, name)
	}
	return nil, nil
}

func (m *MockEmployeeService) BuyItem(ctx context.Context, name, item string, shop shop.Shop, t transfer.Service) error {
	if m.BuyItemFn != nil {
		return m.BuyItemFn(ctx, name, item, shop, t)
	}
	return nil
}

func (m *MockEmployeeService) TransferCoins(ctx context.Context, sender, receiver string, amount int, t transfer.Service) error {
	if m.TransferCoinsFn != nil {
		return m.TransferCoinsFn(ctx, sender, receiver, amount, t)
	}
	return nil
}
