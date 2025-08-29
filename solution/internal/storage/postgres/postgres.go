package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(
	user string,
	password string,
	dbName string,
	host string,
	port string,
) (*Storage, error) {
	const op = "storage.postgres.New"

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		coins INT CHECK (coins > -1),
		bought_items VARCHAR(50)[] REFERENCES items(name)
	);
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		price INT CHECK (price > 0) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS transfers (
		id SERIAL PRIMARY KEY,
		sender_id INT REFERENCES employees(id),
		receiver_id INT REFERENCES employees(id),
		amount INT CHECK (amount > 0) NOT NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
