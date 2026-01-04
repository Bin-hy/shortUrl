package service

import "context"

func (s *ShortUrlService) GenerateShortUrlV1(ctx context.Context, long_url string) (short_url string, error error) {

	shortUrl, err := s.uc.GenerateShortUrlV1(ctx, long_url)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func (s *ShortUrlService) GetLongUrlV1(ctx context.Context, ShortUrl string) (long_url string, error error) {
	longUrl, err := s.uc.GetLongUrlV1(ctx, ShortUrl)
	if err != nil {
		return "", err
	}

	return longUrl, nil
}
