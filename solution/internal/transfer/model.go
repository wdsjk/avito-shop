package transfer

type Transfer struct {
	ID           int    `db:"id"`
	SenderName   string `db:"sender_name"`
	ReceiverName string `db:"receiver_name"` // "" if transfer to shop
	Amount       int    `db:"amount"`
}

type TransferDto struct {
	SenderName   string `json:"sender_name"`
	ReceiverName string `json:"receiver_name"`
	Amount       int    `json:"amount"`
}

func ToDto(t *Transfer) *TransferDto {
	return &TransferDto{
		SenderName:   t.SenderName,
		ReceiverName: t.ReceiverName,
		Amount:       t.Amount,
	}
}
