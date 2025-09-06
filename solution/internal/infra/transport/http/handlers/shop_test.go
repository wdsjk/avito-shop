package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	dbErr "github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/test"
	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
)

func TestShopHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		item           string
		buyItemErr     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			username:       "username",
			item:           "t-shirt",
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name:           "Unauthorized",
			username:       "",
			item:           "t-shirt",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"unauthorized"}`,
		},
		{
			name:           "Bad Request (item is required)",
			username:       "username",
			item:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"item param is required"}`,
		},
		{
			name:           "Bad Request (item is not found)",
			username:       "username",
			item:           "non-existing_item",
			buyItemErr:     dbErr.ErrItemNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not found"}`,
		},
		{
			name:           "Bad Request (emp is not found)",
			username:       "non-existing_username",
			item:           "item",
			buyItemErr:     dbErr.ErrEmpNotFound,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not found"}`,
		},
		{
			name:           "Bad Request (no coins)",
			username:       "username",
			item:           "item",
			buyItemErr:     dbErr.ErrNoCoins,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"not enough coins"}`,
		},
		{
			name:           "Internal Server Error (failed to buy item)",
			username:       "username",
			item:           "item",
			buyItemErr:     errors.New("db internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"failed to buy item"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			handler := NewShopHandler(
				&test.MockEmployeeService{
					BuyItemFn: func(ctx context.Context, name, item string, shop shop.Shop, t transfer.Service) error {
						return tt.buyItemErr
					},
				}, nil, nil, nil,
			)

			req := httptest.NewRequest(http.MethodGet, "/api/buy/"+tt.item, nil)
			if tt.username != "" {
				req = req.WithContext(context.WithValue(req.Context(), "username", tt.username))
			}
			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("item", tt.item)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

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
