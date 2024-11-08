package main

import (
	"encoding/json"
	"fmt"
	"go-echo-gorm/internal/config"
	"go-echo-gorm/internal/database"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	e  *echo.Echo
	db = setupTestDB()
)

func setupTestDB() *gorm.DB {
	// カレントディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current directory: %v", err))
	}
	fmt.Println("Current directory:", currentDir)

	// テスト用の環境変数を直接設定
	os.Setenv("APP_ENV", "test")
	os.Setenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/web_app_db_integration_test_go?sslmode=disable")

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// テストDB接続
	db, err := database.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	// マイグレーション
	db.AutoMigrate(&Micropost{})

	return db
}

func TestMain(m *testing.M) {
	// テスト前の準備
	e = echo.New()
	code := m.Run()

	// テスト後のクリーンアップ
	db.Exec("DROP TABLE IF EXISTS microposts")

	os.Exit(code)
}

func clearTable() {
	db.Exec("DELETE FROM microposts")
}

func TestGetMicroposts(t *testing.T) {
	clearTable()

	// テストデータの作成
	testPosts := []Micropost{
		{Title: "First Post"},
		{Title: "Second Post"},
	}
	for _, post := range testPosts {
		db.Create(&post)
	}

	// リクエストの作成
	req := httptest.NewRequest(http.MethodGet, "/microposts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テスト実行
	if assert.NoError(t, getMicroposts(c, db)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスのパース
		var response []Micropost
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, 2, len(response))
		assert.Equal(t, "First Post", response[0].Title)
		assert.Equal(t, "Second Post", response[1].Title)
	}
}
