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

	"github.com/go-playground/validator/v10"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/test"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

func TestCoinHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		input          string
		expectedErr    error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			username:       "username",
			input:          `{"toUser":"receiver","amount":100}`,
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name:           "Unauthorized",
			username:       "",
			input:          `{"toUser":"receiver","amount":100}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"unauthorized"}`,
		},
		{
			name:           "Bad Request (invalid json body - invalid input)",
			username:       "username",
			input:          `{"toUser":1,"amount":100}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"invalid JSON body"}`,
		},
		{
			name:           "Bad Request (invalid json body - validation error)",
			username:       "username",
			input:          `{"toUser":"","amount":100}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"invalid JSON body"}`,
		},
		{
			name:           "Bad Request (emp is not found)",
			username:       "non-existing_username",
			input:          `{"toUser":"receiver","amount":100}`,
			expectedErr:    dbErr.ErrEmpNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not found"}`,
		},
		{
			name:           "Bad Request (no coins)",
			username:       "username",
			input:          `{"toUser":"receiver","amount":100000}`,
			expectedErr:    dbErr.ErrNoCoins,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not enough coins"}`,
		},
		{
			name:           "Internal Server Error (failed to get transfer info)",
			username:       "username",
			input:          `{"toUser":"receiver","amount":100}`,
			expectedErr:    errors.New("db internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"failed to get transfer info"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			handler := NewCoinHandler(
				&test.MockEmployeeService{
					TransferCoinsFn: func(ctx context.Context, sender, receiver string, amount int, t transfer.Service) error {
						return tt.expectedErr
					},
				},
				&test.MockTransferService{
					SaveTransferFn: func(ctx context.Context, senderName, receiverName string, amount int) error {
						return tt.expectedErr
					},
				}, validator.New(), slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
			)

			req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
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
