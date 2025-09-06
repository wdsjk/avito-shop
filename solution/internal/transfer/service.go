package transfer

import "context"

type Service interface {
	SaveTransfer(ctx context.Context, senderName, receiverName string, amount int) error
	GetTransfersByEmployee(ctx context.Context, name string) ([]*TransferDto, error)
}

type TransferService struct {
	repo Repository
}

func NewTransferService(repo Repository) *TransferService {
	return &TransferService{repo: repo}
}

func (s *TransferService) SaveTransfer(ctx context.Context, senderName, receiverName string, amount int) error {
	return s.repo.SaveTransfer(ctx, senderName, receiverName, amount)
}

func (s *TransferService) GetTransfersByEmployee(ctx context.Context, name string) ([]*TransferDto, error) {
	transfers, err := s.repo.GetTransfersByEmployee(ctx, name)
	if err != nil {
		return nil, err
	}

	var dtos []*TransferDto
	for _, t := range transfers {
		dtos = append(dtos, ToDto(t))
	}

	return dtos, nil
}
