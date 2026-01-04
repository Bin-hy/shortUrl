package biz

import (
	"context"

	"github.com/BitofferHub/pkg/middlewares/snowflake"
)

// v3 使用雪花算法生成的唯一ID 来生成短链
// 使用Redi 缓存， 加速长短链的查询

func (uc *UrlMapUseCase) GenerateShortUrlV3(ctx context.Context, url string) (string, error) {
	// 1. 查缓存中是否有对应的短链
	shortUrl, err := uc.repo.GetShortUrlFormCache(ctx, url)
	if err != nil {
		return "", err
	}
	if shortUrl != "" {
		return shortUrl, nil
	}

	// 2. 缓存中没有，查数据库
	shortUrl, err = uc.repo.GetShortUrlFormDb(ctx, url)
	if err != nil {
		return "", err
	}
	// 有，同时保存到缓存中
	if shortUrl != "" {
		uc.repo.SaveToCache(ctx, url, shortUrl)
		return shortUrl, nil
	}
	// 3. 数据库也没有， 使用雪花算法生成ID
	id := snowflake.GenID()

	// 4. 生成短链
	shortUrl = generateShortUrl(id)

	// 5. 存储到布隆过滤器
	err = uc.repo.SaveToBloomFilter(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	// 6. 存储到缓存
	err = uc.repo.SaveToCache(ctx, url, shortUrl)
	if err != nil {
		return "", err
	}
	// 7. 存储到数据库
	_, err = uc.repo.CreateToDb(ctx, &UrlMap{
		LongUrl:  url,
		ShortUrl: shortUrl,
	})
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

// 获取长链
func (uc *UrlMapUseCase) GetLongUrlV3(ctx context.Context, shortUrl string) (string, error) {
	// 查布隆过滤器和缓存
	need, longUrl, err := uc.repo.FindShortUrlFormBloomFilterAndCache(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	if longUrl != "" {
		return longUrl, nil
	}

	// 不需要查DB， 直接return
	if need == 0 {
		return "", nil
	}

	// 查数据库
	return uc.repo.GetLongUrlFormDb(ctx, shortUrl)
}
