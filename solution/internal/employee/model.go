package employee

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	errTypeAssertion = errors.New("type assertion to []byte failed")
)

type Inventory map[string]int

type Employee struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Password  string `db:"password"`
	Coins     int    `db:"coins"`
	Inventory `db:"bought_items"`
}

func (i Inventory) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Inventory) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errTypeAssertion
	}

	m := make(map[string]int)
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}

	*i = m
	return nil
}
