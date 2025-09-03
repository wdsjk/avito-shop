package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

type ShopHandler struct {
	employeeService *employee.EmployeeService
	transferService *transfer.TransferService
	shop            shop.Shop
	log             *slog.Logger
}

func NewShopHandler(
	employeeService *employee.EmployeeService, transferService *transfer.TransferService,
	shop shop.Shop, log *slog.Logger,
) *ShopHandler {
	return &ShopHandler{
		employeeService: employeeService,
		transferService: transferService,
		shop:            shop,
		log:             log,
	}
}

func (h *ShopHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	item := chi.URLParam(r, "item")
	if item == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(utils.MakeErr("item param is required"))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	err := h.employeeService.BuyItem(r.Context(), username, item, h.shop, h.transferService)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case errors.Is(err, dbErr.ErrEmpNotFound) || errors.Is(err, dbErr.ErrItemNotFound):
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(utils.MakeErr("not found"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		case errors.Is(err, dbErr.ErrNoCoins):
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(utils.MakeErr("not enough coins"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(utils.MakeErr("failed to buy item"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
