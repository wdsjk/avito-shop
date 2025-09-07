package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type CoinHandler struct {
	employeeService employee.Service
	transferService transfer.Service
	valid           *validator.Validate
	log             *slog.Logger
}

func NewCoinHandler(
	employeeService employee.Service,
	transferService transfer.Service,
	valid *validator.Validate, log *slog.Logger,
) *CoinHandler {
	return &CoinHandler{
		employeeService: employeeService,
		transferService: transferService,
		valid:           valid,
		log:             log,
	}
}

func (h *CoinHandler) Handle(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(utils.MakeErr("unauthorized"))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
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
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}
	defer r.Body.Close()
	if err := h.valid.Struct(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(utils.MakeErr("invalid JSON body"))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	err := h.employeeService.TransferCoins(r.Context(), username, req.ToUser, req.Amount, h.transferService)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case errors.Is(err, dbErr.ErrEmpNotFound):
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(utils.MakeErr("not found"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		case errors.Is(err, dbErr.ErrNoCoins):
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(utils.MakeErr("not enough coins"))
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		default:
			h.log.Error("failed to get transfer info", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(utils.MakeErr("failed to get transfer info"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
