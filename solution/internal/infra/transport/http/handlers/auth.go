package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
	"github.com/wdsjk/avito-shop/internal/lib/mapper"
	"github.com/wdsjk/avito-shop/internal/lib/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	employeeService *employee.EmployeeService
}

func NewAuthHandler(employeeService *employee.EmployeeService) *AuthHandler {
	return &AuthHandler{
		employeeService: employeeService,
	}
}

func (h *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req handlers_dto.AuthRequest
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

	emp, err := h.employeeService.GetEmployee(r.Context(), req.Username)
	if errors.Is(err, dbErr.ErrEmpNotFound) {
		name, err := h.employeeService.SaveEmployee(r.Context(), req.Username, req.Password)
		if name == "" && err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(utils.MakeErr("failed to save employee"))
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		token, err := utils.GenerateJWT(req.Username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(utils.MakeErr("failed to generate token"))
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(mapper.AuthResponse(token))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	} else if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(utils.MakeErr("failed to get employee info"))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(emp.Password), []byte(req.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(utils.MakeErr("invalid username or password"))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	token, err := utils.GenerateJWT(req.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(utils.MakeErr("failed to generate token"))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapper.AuthResponse(token))
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
