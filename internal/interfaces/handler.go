package interfaces

import (
	"github.com/Bin-hy/shortUrl/internal/service"
)

type Handler struct {
	svc *service.ShortUrlService
}

func NewHandler(svc *service.ShortUrlService) *Handler {
	return &Handler{
		svc: svc,
	}
}

type ShortenReq struct {
	LongUrl string `json:"long_url" binding:"required"`
}

type ShortenResp struct {
	ShortUrl string `json:"short_url"`
}
