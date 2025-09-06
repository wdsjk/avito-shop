package mapper

import (
	"github.com/wdsjk/avito-shop/internal/employee"
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

func InfoResponse(emp *employee.EmployeeDto, coinHistory []*transfer.TransferDto) *handlers_dto.InfoResponse {
	resp := &handlers_dto.InfoResponse{
		Coins: emp.Coins,
		Inventory: make([]struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		}, 0, len(emp.Inventory)),
		CoinHistory: struct {
			Received []struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			} `json:"received"`
			Sent []struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			} `json:"sent"`
		}{
			Received: make([]struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			}, 0),
			Sent: make([]struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			}, 0),
		},
	}

	for item, count := range emp.Inventory {
		resp.Inventory = append(resp.Inventory, struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		}{
			Type:     item,
			Quantity: count,
		})
	}

	for _, t := range coinHistory {
		switch emp.Name {
		case t.ReceiverName:
			resp.CoinHistory.Received = append(resp.CoinHistory.Received, struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			}{
				FromUser: t.SenderName,
				Amount:   t.Amount,
			})
		case t.SenderName:
			resp.CoinHistory.Sent = append(resp.CoinHistory.Sent, struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			}{
				ToUser: t.ReceiverName,
				Amount: t.Amount,
			})
		}
	}

	return resp
}
