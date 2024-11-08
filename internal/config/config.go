package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// プロジェクトルートからの相対パスを構築
	envFile := fmt.Sprintf(".env.%s", env)
	paths := []string{
		filepath.Join("configs", envFile),             // 通常のパス
		filepath.Join("..", "configs", envFile),       // 1つ上のディレクトリから
		filepath.Join("..", "..", "configs", envFile), // 2つ上のディレクトリから
	}

	// いずれかのパスで.envファイルを読み込む
	var loaded bool
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				loaded = true
				break
			}
		}
	}

	// テスト環境の場合は、環境変数が設定されていれば.envファイルは必須としない
	if !loaded && env != "test" {
		return nil, fmt.Errorf("error loading .env.%s file", env)
	}

	return &Config{
		Environment: env,
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}, nil
}
