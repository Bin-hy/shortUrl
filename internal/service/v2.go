package service

import "context"

func (s *ShortUrlService) GenerateShortUrlV2(ctx context.Context, long_url string) (short_url string, error error) {
	return s.uc.GenerateShortUrlV2(ctx, long_url)
}

func (s *ShortUrlService) GetLongUrlV2(ctx context.Context, ShortUrl string) (long_url string, error error) {
	return s.uc.GetLongUrlV2(ctx, ShortUrl)
}
