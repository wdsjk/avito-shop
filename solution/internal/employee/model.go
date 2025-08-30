package employee

type Employee struct {
	ID       int      `db:"id"`
	Name     string   `db:"name"`
	Password string   `db:"password"`
	Coins    float64  `db:"coins"`
	Items    []string `db:"bought_items"`
}

type EmployeeDto struct {
	Name  string   `json:"name"`
	Coins float64  `json:"coins"`
	Items []string `json:"bought_items"`
}
