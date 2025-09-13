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
	"github.com/wdsjk/avito-shop/internal/employee"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/test"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler(t *testing.T) {
	oldGenerateJwt := generateJWT
	defer func() {
		generateJWT = oldGenerateJwt
	}()

	tests := []struct {
		name            string
		username        string
		password        string
		input           string
		existingUser    bool
		expectedSaveErr error
		expectedGetErr  error
		expectedJwtErr  error
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:           "Success (New user)",
			username:       "username",
			password:       "qwerty",
			input:          `{"username":"username","password":"qwerty"}`,
			existingUser:   false,
			expectedGetErr: dbErr.ErrEmpNotFound,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"some token"}`,
		},
		{
			name:           "Success (Old user)",
			username:       "username",
			password:       "qwerty",
			input:          `{"username":"username","password":"qwerty"}`,
			existingUser:   true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"some token"}`,
		},
		{
			name:           "Unauthorized",
			username:       "username",
			password:       "qwerty",
			input:          `{"username":"username","password":"incorrectPassword"}`,
			existingUser:   true,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"invalid username or password"}`,
		},
		{
			name:           "Bad Request",
			username:       "username",
			input:          `{"userName":"username","pasSsword":123}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"invalid JSON body"}`,
		},
		{
			name:            "Internal Server Error (failed to save employee)",
			username:        "username",
			input:           `{"username":"username","password":"qwerty"}`,
			expectedGetErr:  dbErr.ErrEmpNotFound,
			expectedSaveErr: errors.New("db internal error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedBody:    `{"errors":"failed to save employee"}`,
		},
		{
			name:           "Internal Server Error (failed to get employee info)",
			username:       "username",
			input:          `{"username":"username","password":"qwerty"}`,
			expectedGetErr: errors.New("db internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"failed to get employee info"}`,
		},
		{
			name:           "Internal Server Error (failed to generate token)",
			username:       "username",
			input:          `{"username":"username","password":"qwerty"}`,
			expectedGetErr: dbErr.ErrEmpNotFound,
			expectedJwtErr: errors.New("failed to generate token"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"failed to generate token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			generateJWT = func(username string, secret []byte) (string, error) {
				return "some token", tt.expectedJwtErr
			}

			handler := NewAuthHandler(
				&test.MockEmployeeService{
					SaveEmployeeFn: func(ctx context.Context, name, password string) (string, error) {
						return tt.username, tt.expectedSaveErr
					},
					GetEmployeeFn: func(ctx context.Context, name string) (*employee.EmployeeDto, error) {
						if tt.existingUser {
							hash, _ := bcrypt.GenerateFromPassword([]byte(tt.password), bcrypt.DefaultCost)
							return &employee.EmployeeDto{Name: tt.username, Password: string(hash)}, tt.expectedGetErr
						}
						return &employee.EmployeeDto{Name: tt.username}, tt.expectedGetErr
					},
				}, validator.New(), slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
			)

			req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
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
