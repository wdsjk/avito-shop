package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type CoinHandler struct {
	employeeService *employee.EmployeeService
	transferService *transfer.TransferService
}

func NewCoinHandler(employeeService *employee.EmployeeService, transferService *transfer.TransferService) *CoinHandler {
	return &CoinHandler{
		employeeService: employeeService,
		transferService: transferService,
	}
}

func (h *CoinHandler) Handle(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(utils.MakeErr("unauthorized"))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	var req handlers_dto.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(utils.MakeErr("invalid JSON body"))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}
	defer r.Body.Close()

	err := h.employeeService.TransferCoins(r.Context(), username, req.ToUser, req.Amount, h.transferService)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, dbErr.ErrEmpNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(utils.MakeErr("not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(utils.MakeErr("failed to get transfer info"))
		}
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
