package transfer

import (
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	saveTransferFn           func(ctx context.Context, senderName, receiverName string, amount int) error
	getTransfersByEmployeeFn func(ctx context.Context, name string) ([]*Transfer, error)
}

func (m *mockRepo) SaveTransfer(ctx context.Context, senderName, receiverName string, amount int) error {
	if m.saveTransferFn != nil {
		return m.saveTransferFn(ctx, senderName, receiverName, amount)
	}
	return nil
}

func (m *mockRepo) GetTransfersByEmployee(ctx context.Context, name string) ([]*Transfer, error) {
	if m.getTransfersByEmployeeFn != nil {
		return m.getTransfersByEmployeeFn(ctx, name)
	}
	return []*Transfer{}, nil
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

func TestSaveTransfer_Success(t *testing.T) {
	// Arrange
	senderName := "sender"
	receiverName := "receiver"
	amount := 100
	service := NewTransferService(&mockRepo{
		saveTransferFn: func(ctx context.Context, senderName, receiverName string, amount int) error {
			return nil
		},
	})

	// Act
	err := service.SaveTransfer(context.Background(), senderName, receiverName, amount)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveTransfer_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("something went wrong")
	service := NewTransferService(&mockRepo{
		saveTransferFn: func(ctx context.Context, senderName, receiverName string, amount int) error {
			return expectedErr
		},
	})

	// Act
	err := service.SaveTransfer(context.Background(), "", "", 0)

	// Assert
	assertErrorIs(t, err, expectedErr)
}

func TestGetTransfersByEmployee_Success(t *testing.T) {
	// Arrange
	senderName := "name"
	service := NewTransferService(&mockRepo{
		getTransfersByEmployeeFn: func(ctx context.Context, name string) ([]*Transfer, error) {
			return []*Transfer{
				{
					ID:           1,
					SenderName:   senderName,
					ReceiverName: "",
					Amount:       100,
				},
			}, nil
		},
	})

	// Act
	trArr, err := service.GetTransfersByEmployee(context.Background(), senderName)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trArr) != 1 {
		t.Errorf("expected: %d, got: %d", 1, len(trArr))
	}

	tr := trArr[0]
	if tr.SenderName != senderName {
		t.Errorf("expected: %s, got: %s", senderName, tr.SenderName)
	}
	if tr.ReceiverName != "" {
		t.Errorf("expected: %s, got: %s", "", tr.ReceiverName)
	}
	if tr.Amount != 100 {
		t.Errorf("expected: %d, got: %d", 100, tr.Amount)
	}
}

func TestGetTransfersByEmployee_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("transfer not found")
	service := NewTransferService(&mockRepo{
		getTransfersByEmployeeFn: func(ctx context.Context, name string) ([]*Transfer, error) {
			return []*Transfer{}, expectedErr
		},
	})

	// Act
	_, err := service.GetTransfersByEmployee(context.Background(), "")

	// Assert
	assertErrorIs(t, err, expectedErr)
}
