package verband

import (
	"pkv/api/src/service/captcha"
	verbandService "pkv/api/src/service/verband"
)

type Handler struct {
	service        *verbandService.Service
	captchaService *captcha.Service
}

func NewHandler(service *verbandService.Service, captchaService *captcha.Service) *Handler {
	return &Handler{service: service, captchaService: captchaService}
}
