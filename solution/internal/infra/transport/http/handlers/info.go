package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/lib/mapper"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type InfoHandler struct {
	employeeService *employee.EmployeeService
	transferService *transfer.TransferService
}

func NewInfoHandler(employeeService *employee.EmployeeService, transferService *transfer.TransferService) *InfoHandler {
	return &InfoHandler{
		employeeService: employeeService,
		transferService: transferService,
	}
}

func (h *InfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	emp, err := h.employeeService.GetEmployee(r.Context(), username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, dbErr.ErrEmpNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(utils.MakeErr("not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(utils.MakeErr("failed to get employee info"))
		}
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	ts, err := h.transferService.GetTransfersByEmployee(r.Context(), username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, dbErr.ErrTransferNotFound) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapper.InfoResponse(emp, ts))
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
