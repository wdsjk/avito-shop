package mapper

import (
	handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"
)

func AuthResponse(token string) *handlers_dto.AuthResponse {
	return &handlers_dto.AuthResponse{
		Token: token,
	}
}
