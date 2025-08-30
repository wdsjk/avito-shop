package employee

import "context"

type Service struct {
	repo Repository
}

func NewEmployeeService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveEmployee(ctx context.Context, name string, password string) (string, error) {
	return s.repo.SaveEmployee(ctx, name, password)
}

func (s *Service) GetEmployee(ctx context.Context, name string) (*EmployeeDto, error) {
	employee, err := s.repo.GetEmployee(ctx, name)
	if err != nil {
		return nil, err
	}

	return &EmployeeDto{
		Name:  employee.Name,
		Coins: employee.Coins,
		Items: employee.Items,
	}, nil
}
