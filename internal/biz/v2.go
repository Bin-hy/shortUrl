package biz

import (
	"context"
	"errors"
	"strconv"
)

// 使用Redis自增 ID 来生成短链，替代 v1的MYSQL 主键
// 使用布隆过滤器，优化DB查询短链的性能

func (uc *UrlMapUseCase) GenerateShortUrlV2(ctx context.Context, longUrl string) (string, error) {
	// 1. 先查询数据库是否有该长链
	shortUrl, err := uc.repo.GetShortUrlFormDb(ctx, longUrl)
	if err != nil {
		return "", err
	}
	// 有短链，直接返回
	if shortUrl != "" {
		return shortUrl, nil
	}

	// 2. 从 redis 缓存中获取 ID
	idStr, err := uc.repo.GenerateIdFromCache(ctx)
	if err != nil {
		return "", err
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)

	// 3. 利用 base62算法，生成短链
	shortUrl = generateShortUrl(id)

	//  4. 保存到布隆过滤器中
	err = uc.repo.SaveToBloomFilter(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	// 5. 保存到数据库中
	_, err = uc.repo.CreateToDb(ctx, &UrlMap{
		LongUrl:  longUrl,
		ShortUrl: shortUrl,
	})
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

// 获取长链
func (uc *UrlMapUseCase) GetLongUrlV2(ctx context.Context, shortUrl string) (string, error) {
	// 1. 从布隆过滤器中查询
	exist, err := uc.repo.FindShortUrlFormBloomFilter(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	// 如果不在布隆过滤器中，一定不在DB中， 提前返回即可
	if !exist {
		return "", errors.New("短链不存在") // TODO!
	}
	return uc.repo.GetLongUrlFormDb(ctx, shortUrl)
}
