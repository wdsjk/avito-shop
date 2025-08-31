package employee

type Employee struct {
	ID       int      `db:"id"`
	Name     string   `db:"name"`
	Password string   `db:"password"`
	Coins    int      `db:"coins"`
	Items    []string `db:"bought_items"`
}

type EmployeeDto struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Coins    int      `json:"coins"`
	Items    []string `json:"bought_items"`
}

func ToDto(e *Employee) *EmployeeDto {
	return &EmployeeDto{
		Name:  e.Name,
		Coins: e.Coins,
		Items: e.Items,
	}
}
