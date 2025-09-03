package handlers

import (
	"github.com/wdsjk/avito-shop/internal/employee"
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

// func (h *CoinHandler) Handle(w http.ResponseWriter, r *http.Request) {
// 	username, ok := r.Context().Value("username").(string)
// 	if !ok || username == "" {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		err := json.NewEncoder(w).Encode(makeErr("unauthorized"))
// 		if err != nil {
// 			http.Error(w, "failed to encode response", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	emp, err := h.employeeService.GetEmployee(r.Context(), username)
// 	if err != nil {
// 		http.Error(w, "failed to get employee info", http.StatusInternalServerError)
// 		return
// 	}

// 	ts, err := h.transferService.GetTransfersByEmployee(r.Context(), username)
// 	if err != nil {
// 		http.Error(w, "failed to get transfer history", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	err = json.NewEncoder(w).Encode(mapper.MapResponse(emp, ts))
// 	if err != nil {
// 		http.Error(w, "failed to encode response", http.StatusInternalServerError)
// 		return
// 	}
// }
