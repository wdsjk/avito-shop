package employee

import (
	"context"
	"errors"
	"testing"

	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type mockRepo struct {
	saveEmployeeFn  func(ctx context.Context, name, password string) (string, error)
	getEmployeeFn   func(ctx context.Context, name string) (*Employee, error)
	buyItemFn       func(ctx context.Context, name, item string, shop shop.Shop, t *transfer.TransferService) error
	transferCoinsFn func(ctx context.Context, sender, receiver string, amount int, t *transfer.TransferService) error
}

func (m *mockRepo) SaveEmployee(ctx context.Context, name, password string) (string, error) {
	if m.saveEmployeeFn != nil {
		return m.saveEmployeeFn(ctx, name, password)
	}
	return "", nil
}

func (m *mockRepo) GetEmployee(ctx context.Context, name string) (*Employee, error) {
	if m.getEmployeeFn != nil {
		return m.getEmployeeFn(ctx, name)
	}
	return &Employee{}, nil
}

func (m *mockRepo) BuyItem(ctx context.Context, name, item string, shop shop.Shop, t *transfer.TransferService) error {
	if m.buyItemFn != nil {
		return m.buyItemFn(ctx, name, item, shop, t)
	}
	return nil
}

func (m *mockRepo) TransferCoins(ctx context.Context, sender, receiver string, amount int, t *transfer.TransferService) error {
	if m.transferCoinsFn != nil {
		return m.transferCoinsFn(ctx, sender, receiver, amount, t)
	}
	return nil
}

type mockTransferRepo struct {
	saveTransferFn           func(ctx context.Context, senderName, receiverName string, amount int) error
	getTransfersByEmployeeFn func(ctx context.Context, name string) ([]*transfer.Transfer, error)
}

func (m *mockTransferRepo) SaveTransfer(ctx context.Context, senderName, receiverName string, amount int) error {
	if m.saveTransferFn != nil {
		return m.saveTransferFn(ctx, senderName, receiverName, amount)
	}
	return nil
}

func (m *mockTransferRepo) GetTransfersByEmployee(ctx context.Context, name string) ([]*transfer.Transfer, error) {
	if m.getTransfersByEmployeeFn != nil {
		return m.getTransfersByEmployeeFn(ctx, name)
	}
	return []*transfer.Transfer{}, nil
}

func assertErrorIs(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(got, want) {
		t.Errorf("expected: %v, got: %v", want, got)
	}
}

func TestSaveEmployee_Success(t *testing.T) {
	// Arrange
	expectedName := "name"
	expectedPassword := "password"
	service := NewEmployeeService(&mockRepo{
		saveEmployeeFn: func(ctx context.Context, name, password string) (string, error) {
			if name != expectedName || password != expectedPassword {
				t.Errorf("unexpected params: (%s, %s)", name, password)
			}
			return name, nil
		},
	})

	// Act
	name, err := service.SaveEmployee(context.Background(), expectedName, expectedPassword)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != expectedName {
		t.Errorf("expected: %s, got: %s", expectedName, name)
	}
}

func TestSaveEmployee_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("some err with saving")
	service := NewEmployeeService(&mockRepo{
		saveEmployeeFn: func(ctx context.Context, name, password string) (string, error) {
			return "", expectedErr
		},
	})

	// Act
	_, err := service.SaveEmployee(context.Background(), "", "")

	// Assert
	assertErrorIs(t, err, expectedErr)
}

func TestGetEmployee_Success(t *testing.T) {
	// Arrange
	expectedName := "name"
	service := NewEmployeeService(&mockRepo{
		getEmployeeFn: func(ctx context.Context, name string) (*Employee, error) {
			if name != expectedName {
				t.Errorf("unexpected param: %s", name)
			}
			return &Employee{
				ID:        1,
				Name:      name,
				Password:  "",
				Coins:     1000,
				Inventory: Inventory{},
			}, nil
		},
	})

	// Act
	emp, err := service.GetEmployee(context.Background(), expectedName)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if emp.ID != 1 {
		t.Errorf("expected: %d, got: %d", 1, emp.ID)
	}
	if emp.Name != expectedName {
		t.Errorf("expected: %s, got: %s", expectedName, emp.Name)
	}
	if emp.Password != "" {
		t.Errorf("expected: %s, got: %s", "", emp.Password)
	}
	if emp.Coins != 1000 {
		t.Errorf("expected: %d, got: %d", 1000, emp.Coins)
	}
}

func TestGetEmployee_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("can't save employee")
	service := NewEmployeeService(&mockRepo{
		getEmployeeFn: func(ctx context.Context, name string) (*Employee, error) {
			return nil, expectedErr
		},
	})

	// Act
	_, err := service.GetEmployee(context.Background(), "")

	// Assert
	assertErrorIs(t, err, expectedErr)
}

func TestBuyItem_Success(te *testing.T) {
	// Arrange
	employee := &Employee{
		ID:        1,
		Name:      "name",
		Password:  "",
		Coins:     1000,
		Inventory: Inventory{},
	}
	expectedItem := "t-shirt"
	s := shop.NewShop()
	expectedAmount := 80
	tr := transfer.NewTransferService(&mockTransferRepo{
		saveTransferFn: func(ctx context.Context, senderName, receiverName string, amount int) error {
			if senderName != employee.Name || receiverName != "" || amount != expectedAmount {
				te.Errorf("unexpected params: (%s, %s, %d)", senderName, receiverName, amount)
			}
			employee.Coins = employee.Coins - amount
			return nil
		},
	})
	service := NewEmployeeService(&mockRepo{
		buyItemFn: func(ctx context.Context, name, item string, s shop.Shop, t *transfer.TransferService) error {
			if _, ok := s[item]; !ok {
				te.Errorf("no item in shop: %s", item)
			}

			if employee.Coins-s[item] < 0 {
				te.Errorf("no coins: %d", employee.Coins)
			}

			employee.Inventory[item] = employee.Inventory[item] + 1

			err := t.SaveTransfer(ctx, name, "", expectedAmount)
			if err != nil {
				te.Errorf("unexpected error: %v", err)
			}

			return nil
		},
	})

	// Act
	err := service.BuyItem(context.Background(), employee.Name, expectedItem, s, tr)

	// Assert
	if err != nil {
		te.Fatalf("unexpected error: %v", err)
	}
	if employee.Coins != 1000-expectedAmount {
		te.Errorf("expected: %d, got: %d", 1000-expectedAmount, employee.Coins)
	}
	if len(employee.Inventory) < 1 {
		te.Errorf("expected: %v, got: %v", Inventory{
			"t-shirt": 1,
		}, employee.Inventory)
	}
}

func TestBuyItem_Fail(t *testing.T) {
	// Arrange
	expectedNoCoinsErr := errors.New("no coins")
	expectedNotFoundErr := errors.New("not found")
	tests := []struct {
		name string
		want error
	}{
		{
			name: "NotFound",
			want: expectedNotFoundErr,
		},
		{
			name: "NoCoins",
			want: expectedNoCoinsErr,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transfer.NewTransferService(&mockTransferRepo{})
			service := NewEmployeeService(&mockRepo{
				buyItemFn: func(ctx context.Context, name, item string, s shop.Shop, t *transfer.TransferService) error {
					return tt.want
				},
			})

			err := service.BuyItem(context.Background(), "", "", shop.NewShop(), tr)

			assertErrorIs(t, err, tt.want)
		})
	}
}

func TestTransferCoins_Success(te *testing.T) {
	// Arrange
	senderEmp := &Employee{
		ID:        1,
		Name:      "sender",
		Password:  "",
		Coins:     1000,
		Inventory: Inventory{},
	}
	receiverEmp := &Employee{
		ID:        1,
		Name:      "receiver",
		Password:  "",
		Coins:     0,
		Inventory: Inventory{},
	}
	expectedAmount := 100
	tr := transfer.NewTransferService(&mockTransferRepo{
		saveTransferFn: func(ctx context.Context, senderName, receiverName string, amount int) error {
			if senderName != senderEmp.Name || receiverName != receiverEmp.Name || amount != expectedAmount {
				te.Errorf("unexpected params: (%s, %s, %d)", senderName, receiverName, amount)
			}
			senderEmp.Coins = senderEmp.Coins - amount
			receiverEmp.Coins = receiverEmp.Coins + amount
			return nil
		},
	})
	service := NewEmployeeService(&mockRepo{
		transferCoinsFn: func(ctx context.Context, sender, receiver string, amount int, t *transfer.TransferService) error {
			if senderEmp.Coins-amount < 0 {
				te.Errorf("no coins: %d", senderEmp.Coins)
			}

			err := t.SaveTransfer(ctx, sender, receiver, amount)
			if err != nil {
				te.Errorf("unexpected error: %v", err)
			}

			return nil
		},
	})

	// Act
	err := service.TransferCoins(context.Background(), senderEmp.Name, receiverEmp.Name, expectedAmount, tr)

	// Assert
	if err != nil {
		te.Fatalf("unexpected error: %v", err)
	}
	if senderEmp.Coins != 1000-expectedAmount {
		te.Errorf("expected: %d, got: %d", 1000-expectedAmount, senderEmp.Coins)
	}
	if receiverEmp.Coins != expectedAmount {
		te.Errorf("expected: %d, got: %d", expectedAmount, receiverEmp.Coins)
	}
}

func TestTransferCoins_Fail(t *testing.T) {
	// Arrange
	expectedNoCoinsErr := errors.New("no coins")
	expectedNotFoundErr := errors.New("not found")
	tests := []struct {
		name string
		want error
	}{
		{
			name: "NotFound",
			want: expectedNotFoundErr,
		},
		{
			name: "NoCoins",
			want: expectedNoCoinsErr,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := transfer.NewTransferService(&mockTransferRepo{})
			service := NewEmployeeService(&mockRepo{
				transferCoinsFn: func(ctx context.Context, sender, receiver string, amount int, t *transfer.TransferService) error {
					return tt.want
				},
			})

			err := service.TransferCoins(context.Background(), "", "", 100, tr)

			assertErrorIs(t, err, tt.want)
		})
	}
}
