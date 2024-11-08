package main

import (
	"log"

	"go-echo-gorm/internal/config"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Config loaded:", cfg)
}
