package data

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Bin-hy/shortUrl/internal/biz"
	"gorm.io/gorm"
)

type urlRepo struct {
	data *Data
}

// NewUrlRepo .
func NewUrlRepo(data *Data) biz.UrlMapRepo {
	return &urlRepo{
		data: data,
	}
}

// Data Layer Model
type UrlMap struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LongUrl   string    `gorm:"type:varchar(250);index" json:"long_url"`
	ShortUrl  string    `gorm:"type:varchar(20);uniqueIndex" json:"short_url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (u *UrlMap) TableName() string {
	return "url_map"
}

// Convert helper
func (u *UrlMap) ToBiz() *biz.UrlMap {
	return &biz.UrlMap{
		ID:        u.ID,
		LongUrl:   u.LongUrl,
		ShortUrl:  u.ShortUrl,
		CreatedAt: u.CreatedAt,
	}
}

func FromBiz(u *biz.UrlMap) *UrlMap {
	return &UrlMap{
		ID:        u.ID,
		LongUrl:   u.LongUrl,
		ShortUrl:  u.ShortUrl,
		CreatedAt: u.CreatedAt,
	}
}

func (r *urlRepo) GetLongUrlFormDb(ctx context.Context, shortUrl string) (string, error) {
	var um UrlMap
	if err := r.data.db.WithContext(ctx).Where("short_url = ?", shortUrl).First(&um).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return um.LongUrl, nil
}

func (r *urlRepo) GetShortUrlFormDb(ctx context.Context, longUrl string) (string, error) {
	var um UrlMap
	// Note: LongUrl indexing might be limited by length, but assuming it works for now
	if err := r.data.db.WithContext(ctx).Where("long_url = ?", longUrl).First(&um).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return um.ShortUrl, nil
}

func (r *urlRepo) GetShortUrlFormCache(ctx context.Context, longUrl string) (string, error) {
	val, err := r.data.rdb.Get(ctx, "long:"+longUrl).Result()
	if err != nil {
		return "", nil
	}
	return val, nil
}

func (r *urlRepo) CreateToDb(ctx context.Context, um *biz.UrlMap) (int64, error) {
	dataModel := FromBiz(um)
	dataModel.CreatedAt = time.Now()
	// 使用 FirstOrCreate 避免重复插入错误
	// 注意：这里我们需要根据 LongUrl 来查找
	if err := r.data.db.WithContext(ctx).Where(UrlMap{LongUrl: dataModel.LongUrl}).FirstOrCreate(dataModel).Error; err != nil {
		return 0, err
	}
	return int64(dataModel.ID), nil
}

func (r *urlRepo) SaveToDb(ctx context.Context, um *biz.UrlMap) error {
	dataModel := FromBiz(um)

	// If ID is present, update by ID
	if dataModel.ID != 0 {
		return r.data.db.WithContext(ctx).Model(dataModel).Updates(dataModel).Error
	}

	// If ID is missing, update by LongUrl (needed for V1 logic)
	return r.data.db.WithContext(ctx).Model(&UrlMap{}).Where("long_url = ?", um.LongUrl).Updates(map[string]interface{}{
		"short_url": um.ShortUrl,
	}).Error
}

func (r *urlRepo) SaveToCache(ctx context.Context, longUrl string, shortUrl string) error {
	pipe := r.data.rdb.Pipeline()
	pipe.Set(ctx, "long:"+longUrl, shortUrl, time.Hour*24)
	pipe.Set(ctx, "short:"+shortUrl, longUrl, time.Hour*24)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *urlRepo) GenerateIdFromCache(ctx context.Context) (string, error) {
	id, err := r.data.rdb.Incr(ctx, "global:id_generator").Result()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (r *urlRepo) SaveToBloomFilter(ctx context.Context, shortUrl string) error {
	// Using Set to simulate Bloom Filter
	return r.data.rdb.SAdd(ctx, "bloom:short_urls", shortUrl).Err()
}

func (r *urlRepo) FindShortUrlFormBloomFilter(ctx context.Context, shortUrl string) (bool, error) {
	return r.data.rdb.SIsMember(ctx, "bloom:short_urls", shortUrl).Result()
}

func (r *urlRepo) FindShortUrlFormBloomFilterAndCache(ctx context.Context, shortUrl string) (int64, string, error) {
	// 1. Check Bloom Filter (Simulated with Set)
	exists, err := r.FindShortUrlFormBloomFilter(ctx, shortUrl)
	if err != nil {
		return 0, "", err
	}
	if !exists {
		return 0, "", nil // 0 means not found
	}

	// 2. Check Cache
	longUrl, err := r.data.rdb.Get(ctx, "short:"+shortUrl).Result()
	if err != nil {
		// Cache miss, but Bloom says maybe exists -> need to check DB
		return 1, "", nil
	}

	return 1, longUrl, nil
}
