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
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		price NUMERIC CHECK (price > 0) NOT NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	INSERT INTO items (name, price)
	SELECT name, price
	FROM (VALUES
		('t-shirt', 80),
		('cup', 20),
		('book', 50),
		('pen', 10),
		('powerbank', 200),
		('hoody', 300),
		('umbrella', 200),
		('socks', 10),
		('wallet', 50),
		('pink-hoody', 500)
	) AS new_items(name, price)
	WHERE NOT EXISTS (SELECT 1 FROM items LIMIT 1);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		coins NUMERIC CHECK (coins > -1),
		bought_items VARCHAR(50)[]
	);`)
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
		amount NUMERIC CHECK (amount > 0) NOT NULL
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
