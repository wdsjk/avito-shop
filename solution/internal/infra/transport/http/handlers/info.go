package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wdsjk/avito-shop/internal/employee"
	"github.com/wdsjk/avito-shop/internal/lib/mapper"
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	emp, err := h.employeeService.GetEmployee(r.Context(), username)
	if err != nil {
		http.Error(w, "failed to get employee info", http.StatusInternalServerError)
		return
	}

	ts, err := h.transferService.GetTransfersByEmployee(r.Context(), username)
	if err != nil {
		http.Error(w, "failed to get transfer history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapper.MapResponse(emp, ts))
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
