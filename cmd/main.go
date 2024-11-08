package main

import (
	"log"
	"net/http"

	"go-echo-gorm/internal/config"

	"github.com/labstack/echo/v4"
)

func main() {
	// 設定の読み込み
	cfd, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Println("Config loaded:", cfd)

	// Echoインスタンスの作成
	e := echo.New()

	// ルートの設定
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
