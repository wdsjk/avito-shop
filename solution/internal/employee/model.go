package employee

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (i *Inventory) Scan(src interface{}) error {
	if src == nil {
		*i = make(Inventory)
		return nil
	}

	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("%s", errTypeAssertion)
	}

	var m map[string]int
	if err := json.Unmarshal(bytes, &m); err != nil {
		return fmt.Errorf("failed to unmarshal JSONB: %w", err)
	}

	*i = Inventory(m)
	return nil
}
