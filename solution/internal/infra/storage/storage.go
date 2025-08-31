package storage

import (
	"database/sql"
	"fmt"

	"github.com/wdsjk/avito-shop/internal/config"
)

func NewStorage(config *config.Config) (*sql.DB, error) {
	const op = "infra.storage.storage.NewStorage"

	dns := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName,
	)
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: migrations
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		coins INT CHECK (coins > -1),
		bought_items VARCHAR(50)[]
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// basically a shop reference into the employees table
	stmt, err = db.Prepare(`
	INSERT INTO employees (name, password, coins, bought_items)
	SELECT name, password, coins, bought_items
	FROM (VALUES
		('', '', 0, ARRAY[]::VARCHAR[]) 
	) AS new_employee(name, password, coins, bought_items)
	WHERE NOT EXISTS (SELECT 1 FROM employees LIMIT 1);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS transfers (
		id SERIAL PRIMARY KEY,
		sender_name VARCHAR(50) REFERENCES employees(name),
		receiver_name VARCHAR(50) REFERENCES employees(name),
		amount INT CHECK (amount > 0) NOT NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
