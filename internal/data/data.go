// package data 是数据访问层，负责与数据库交互——具体的逻辑
package data

import (
	"context"
	"fmt"
	"time"

	"github.com/Bin-hy/shortUrl/internal/conf"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Data struct {
	db   *gorm.DB
	rdb  *redis.Client
	conf *conf.Data
}

func NewData(c *conf.Data) (*Data, error) {
	// MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Addr,
		c.Database.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(c.Database.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.Database.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(c.Database.MaxIdleTime) * time.Second)

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           c.Redis.Db,
		PoolSize:     c.Redis.PoolSize,
		ReadTimeout:  parseDuration(c.Redis.ReadTimeout),
		WriteTimeout: parseDuration(c.Redis.WriteTimeout),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Data{
		db:   db,
		rdb:  rdb,
		conf: c,
	}, nil
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
