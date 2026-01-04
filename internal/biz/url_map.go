package biz

import (
	"context"
	"time"
)

type UrlMap struct {
	ID       uint      `json:"id"`
	LongUrl  string    `json:"long_url"`
	ShortUrl string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (p *UrlMap) TableName() string {
	return "url_map"
}

// UrlMapRepo 是 URL 映射的仓库接口
type UrlMapRepo interface {
	// 查询长链（从DB）
	GetLongUrlFormDb(context.Context, string) (string, error)

	// 查询短链（从DB）
	GetShortUrlFormDb(context.Context, string) (string, error)
	// 查询短链（从缓存）
	GetShortUrlFormCache(context.Context, string) (string, error)

	CreateToDb(context.Context, *UrlMap) (int64, error)

	SaveToDb(context.Context, *UrlMap) error
	SaveToCache(context.Context, string, string) error

	// 从缓存中获取ID
	GenerateIdFromCache(context.Context) (string, error)
	// 保存短链到布隆过滤器中
	SaveToBloomFilter(context.Context, string) error

	FindShortUrlFormBloomFilter(context.Context, string) (bool, error)

	FindShortUrlFormBloomFilterAndCache(context.Context, string) (int64, string, error)
}

// UrlMapUseCase 是 URL 映射的用例接口
type UrlMapUseCase struct {
	repo UrlMapRepo
	// tm   Transaction
}

func NewUrlMapUseCase(repo UrlMapRepo) *UrlMapUseCase {
	return &UrlMapUseCase{repo: repo}
}
