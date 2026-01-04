package biz

import "context"

func (uc *UrlMapUseCase) GenerateShortUrlV1(ctx context.Context, longUrl string) (string, error) {
	// 1. 先查询数据库是否有该长链
	shortUrl, err := uc.repo.GetShortUrlFormDb(ctx, longUrl)
	if err != nil {
		return "", err
	}
	// 有短链，直接返回
	if shortUrl != "" {
		return shortUrl, nil
	}

	// 2. 没有，子啊数据库里创建一条记录
	id, err := uc.repo.CreateToDb(ctx, &UrlMap{
		LongUrl: longUrl,
	})
	if err != nil {
		return "", err
	}

	// 3.利用 base62算法，生成短链
	shortUrl = generateShortUrl(id)

	// 4. 更新对应记录, 将短链存储到DB中
	uc.repo.SaveToDb(ctx, &UrlMap{
		LongUrl:  longUrl,
		ShortUrl: shortUrl,
	})
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

// GetLongUrl 获取长链
func (uc *UrlMapUseCase) GetLongUrlV1(ctx context.Context, shortUrl string) (string, error) {
	return uc.repo.GetLongUrlFormDb(ctx, shortUrl)
}
