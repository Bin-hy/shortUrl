package service

import "github.com/Bin-hy/shortUrl/internal/biz"

// 短链接服务
type ShortUrlService struct {
	uc *biz.UrlMapUseCase
}

func NewShortUrlService(uc *biz.UrlMapUseCase) *ShortUrlService {
	return &ShortUrlService{
		uc: uc,
	}
}
