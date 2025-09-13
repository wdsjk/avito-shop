package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
	"github.com/wdsjk/avito-shop/internal/lib/mapper"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	secretKey   = []byte(os.Getenv("jwt_secret"))
	generateJWT = utils.GenerateJWT
)

type AuthHandler struct {
	employeeService employee.Service
	valid           *validator.Validate
	log             *slog.Logger
}

func NewAuthHandler(
	employeeService employee.Service,
	valid *validator.Validate,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		employeeService: employeeService,
		valid:           valid,
		log:             logger,
	}
}

func (h *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req handlers_dto.AuthRequest
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

	emp, err := h.employeeService.GetEmployee(r.Context(), req.Username)
	if errors.Is(err, dbErr.ErrEmpNotFound) {
		name, err := h.employeeService.SaveEmployee(r.Context(), req.Username, req.Password)
		if name != req.Username || err != nil {
			h.log.Error("failed to save employee", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(utils.MakeErr("failed to save employee"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		token, err := generateJWT(req.Username, secretKey)
		if err != nil {
			h.log.Error("failed to generate token", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(utils.MakeErr("failed to generate token"))
			if err != nil {
				h.log.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(mapper.AuthResponse(token))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	} else if err != nil {
		h.log.Error("failed to get employee info", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(utils.MakeErr("failed to get employee info"))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(emp.Password), []byte(req.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(utils.MakeErr("invalid username or password"))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	token, err := generateJWT(req.Username, secretKey)
	if err != nil {
		h.log.Error("failed to generate token", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(utils.MakeErr("failed to generate token"))
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapper.AuthResponse(token))
	if err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
