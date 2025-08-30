package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/wdsjk/avito-shop/internal/employee"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) SaveEmployee(ctx context.Context, name string, password string) (string, error) {
	const op = "infra.storage.postgres.saveEmployee"

	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO employees (name, password, coins, bought_items) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, name, hashedPassword, 1000, []string{})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return name, nil
}

func (r *EmployeeRepository) GetEmployee(ctx context.Context, name string) (*employee.Employee, error) {
	const op = "infra.storage.postgres.getEmployeeInfo"

	stmt, err := r.db.PrepareContext(ctx, `SELECT * FROM employees WHERE name=$1;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var emp employee.Employee
	err = stmt.QueryRowContext(ctx, name).Scan(&emp.ID, &emp.Name, &emp.Password, &emp.Coins, (*pq.StringArray)(&emp.Items))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &emp, nil
}
