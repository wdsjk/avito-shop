package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/test"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

func TestInfoHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		infoErr        error
		transferErr    error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			username:       "username",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"coins":0,"inventory":[],"coinHistory":{"received":[],"sent":[]}}`,
		},
		{
			name:           "Unauthorized",
			username:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"unauthorized"}`,
		},
		{
			name:           "Bad Request (emp is not found)",
			username:       "non-existing_username",
			infoErr:        dbErr.ErrEmpNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not found"}`,
		},
		{
			name:           "Bad Request (transfers are not found)",
			username:       "username",
			transferErr:    dbErr.ErrTransferNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not found"}`,
		},
		{
			name:           "Internal Server Error (failed to get employee info)",
			username:       "username",
			infoErr:        errors.New("db internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"failed to get employee info"}`,
		},
		{
			name:           "Internal Server Error (failed to get transfer info)",
			username:       "username",
			transferErr:    errors.New("db internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"failed to get transfer info"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			handler := NewInfoHandler(
				&test.MockEmployeeService{
					GetEmployeeFn: func(ctx context.Context, name string) (*employee.EmployeeDto, error) {
						return &employee.EmployeeDto{
							Name: tt.username,
						}, tt.infoErr
					},
				},
				&test.MockTransferService{
					GetTransfersByEmployeeFn: func(ctx context.Context, name string) ([]*transfer.TransferDto, error) {
						return []*transfer.TransferDto{}, tt.transferErr
					},
				}, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
			)

			req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
			if tt.username != "" {
				req = req.WithContext(context.WithValue(req.Context(), "username", tt.username))
			}
			rr := httptest.NewRecorder()

			// Act
			handler.Handle(rr, req)

			// Assert
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected: %d, got: %d", tt.expectedStatus, rr.Code)
			}
			body := strings.TrimSpace(rr.Body.String())
			if body != tt.expectedBody {
				t.Errorf("expected: %s, got: %s", tt.expectedBody, body)
			}
		})
	}
}
