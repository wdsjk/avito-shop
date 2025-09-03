package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wdsjk/avito-shop/internal/employee"
	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmpNotFound  = errors.New("employee not found")
	ErrItemNotFound = errors.New("item not found")
	ErrNoCoins      = errors.New("not enough coins")
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) SaveEmployee(ctx context.Context, name string, password string) (string, error) {
	const op = "infra.storage.postgres.SaveEmployee"

	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO employees (name, password, coins, bought_items) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, name, hashedPassword, 1000, employee.Inventory{})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return name, nil
}

func (r *EmployeeRepository) GetEmployee(ctx context.Context, name string) (*employee.Employee, error) {
	const op = "infra.storage.postgres.GetEmployeeInfo"

	stmt, err := r.db.PrepareContext(ctx, `SELECT * FROM employees WHERE name=$1;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var emp employee.Employee
	err = stmt.QueryRowContext(ctx, name).Scan(&emp.ID, &emp.Name, &emp.Password, &emp.Coins, &emp.Inventory)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrEmpNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &emp, nil
}

func (r *EmployeeRepository) BuyItem(ctx context.Context, name, item string, shop shop.Shop, t *transfer.TransferService) error {
	const op = "infra.storage.postgres.BuyItem"

	emp, err := r.GetEmployee(ctx, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrEmpNotFound)
	}

	price, ok := shop[item]
	if !ok {
		return fmt.Errorf("%s: %w", op, ErrItemNotFound)
	}

	if emp.Coins-price < 0 {
		return fmt.Errorf("%s: %w", op, ErrNoCoins)
	}

	stmt, err := r.db.PrepareContext(ctx, `UPDATE employees SET coins=$1, bought_items=$2 WHERE name=$3;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, emp.Coins-price, emp.Inventory[item]+1, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = t.SaveTransfer(ctx, name, "", price)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *EmployeeRepository) TransferCoins(ctx context.Context, senderName, receiverName string, amount int, t *transfer.TransferService) error {
	const op = "infra.storage.postgres.TransferCoins"

	sender, err := r.GetEmployee(ctx, senderName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: %w", op, ErrEmpNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	receiver, err := r.GetEmployee(ctx, receiverName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: %w", op, ErrEmpNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if sender.Coins-amount < 0 {
		return fmt.Errorf("%s: %w", op, ErrNoCoins)
	}

	stmt, err := r.db.PrepareContext(ctx, `UPDATE employees SET coins=$1 WHERE name=$2;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, sender.Coins-amount, senderName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, receiver.Coins+amount, receiverName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = t.SaveTransfer(ctx, senderName, receiverName, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
