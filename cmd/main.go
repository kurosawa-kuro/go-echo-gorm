package main

import (
	"log"
	"net/http"

	"go-echo-gorm/internal/config"
	"go-echo-gorm/internal/database"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Micropost モデルの定義
type Micropost struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// DBの初期化
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// マイグレーション
	db.AutoMigrate(&Micropost{})

	e := echo.New()

	// CRUD ハンドラーの登録
	e.POST("/microposts", func(c echo.Context) error {
		return createMicropost(c, db)
	})
	e.GET("/microposts", func(c echo.Context) error {
		return getMicroposts(c, db)
	})
	e.GET("/microposts/:id", func(c echo.Context) error {
		return getMicropost(c, db)
	})
	e.PUT("/microposts/:id", func(c echo.Context) error {
		return updateMicropost(c, db)
	})
	e.DELETE("/microposts/:id", func(c echo.Context) error {
		return deleteMicropost(c, db)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

// Create
func createMicropost(c echo.Context, db *gorm.DB) error {
	micropost := new(Micropost)
	if err := c.Bind(micropost); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := db.Create(micropost).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create micropost"})
	}

	return c.JSON(http.StatusCreated, micropost)
}

// Read (List)
func getMicroposts(c echo.Context, db *gorm.DB) error {
	var microposts []Micropost
	if err := db.Find(&microposts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch microposts"})
	}
	return c.JSON(http.StatusOK, microposts)
}

// Read (Single)
func getMicropost(c echo.Context, db *gorm.DB) error {
	id := c.Param("id")
	micropost := new(Micropost)

	if err := db.First(micropost, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Micropost not found"})
	}

	return c.JSON(http.StatusOK, micropost)
}

// Update
func updateMicropost(c echo.Context, db *gorm.DB) error {
	id := c.Param("id")
	micropost := new(Micropost)

	if err := db.First(micropost, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Micropost not found"})
	}

	if err := c.Bind(micropost); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := db.Save(micropost).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update micropost"})
	}

	return c.JSON(http.StatusOK, micropost)
}

// Delete
func deleteMicropost(c echo.Context, db *gorm.DB) error {
	id := c.Param("id")

	if err := db.Delete(&Micropost{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete micropost"})
	}

	return c.NoContent(http.StatusNoContent)
}
