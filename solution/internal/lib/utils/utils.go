package utils

import handlers_dto "github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers/dto"

func MakeErr(msg string) handlers_dto.ErrorResponse {
	return handlers_dto.ErrorResponse{Errors: msg}
}
