package employee

import "context"

type Repository interface {
	SaveEmployee(ctx context.Context, name string, password string) (string, error)
	GetEmployee(ctx context.Context, name string) (*Employee, error)
}
