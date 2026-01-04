package service

import "context"

func (s *ShortUrlService) GenerateShortUrlV3(ctx context.Context, long_url string) (short_url string, error error) {
	return s.uc.GenerateShortUrlV3(ctx, long_url)
}

func (s *ShortUrlService) GetLongUrlV3(ctx context.Context, ShortUrl string) (long_url string, error error) {
	return s.uc.GetLongUrlV3(ctx, ShortUrl)
}
