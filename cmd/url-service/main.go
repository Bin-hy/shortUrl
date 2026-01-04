package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Bin-hy/shortUrl/internal/biz"
	"github.com/Bin-hy/shortUrl/internal/conf"
	"github.com/Bin-hy/shortUrl/internal/data"
	"github.com/Bin-hy/shortUrl/internal/interfaces"
	"github.com/Bin-hy/shortUrl/internal/service"
	"github.com/BitofferHub/pkg/middlewares/snowflake"
	"github.com/spf13/viper"
)

func main() {
	// 0. Init Snowflake
	// Init might not return error or has different signature. Let's check the library or try to ignore return if it's void.
	// Based on "used as value", it returns nothing.
	snowflake.Init(time.Now(), 1)

	// 1. Load Config
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Init Data
	d, err := data.NewData(&config.Data)
	if err != nil {
		log.Fatalf("failed to init data: %v", err)
	}
	repo := data.NewUrlRepo(d)

	// 3. Init Biz
	uc := biz.NewUrlMapUseCase(repo)

	// 4. Init Service
	svc := service.NewShortUrlService(uc)

	// 5. Init Handler & Router
	h := interfaces.NewHandler(svc)
	router := interfaces.NewRouter(h)

	// 6. Run Server
	addr := config.Server.Http.Addr
	fmt.Printf("Server running at %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func loadConfig() (*conf.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c conf.Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
