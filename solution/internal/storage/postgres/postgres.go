package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
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

	return &Storage{db: db}, nil
}

// Employee operations

func (s *Storage) SaveEmployee(name, password string) (string, error) {
	const op = "storage.postgres.saveEmployee"

	stmt, err := s.db.Prepare(`INSERT INTO employees (name, password, coins, bought_items) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(name, hashedPassword, 1000, []string{})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return name, nil
}
