package main

import (
	"log"
	"net/http"

	"go-echo-gorm/internal/config"
	"go-echo-gorm/internal/database"

	_ "go-echo-gorm/docs" // swagger docs

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// @title Micropost API
// @version 1.0
// @description This is a sample micropost server.
// @host localhost:8080
// @BasePath /
type Micropost struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}

// @Summary Create a new micropost
// @Description Create a new micropost with the provided title
// @Tags microposts
// @Accept json
// @Produce json
// @Param micropost body Micropost true "Micropost object"
// @Success 201 {object} Micropost
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microposts [post]
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

// @Summary Get all microposts
// @Description Get a list of all microposts
// @Tags microposts
// @Produce json
// @Success 200 {array} Micropost
// @Failure 500 {object} map[string]string
// @Router /microposts [get]
func getMicroposts(c echo.Context, db *gorm.DB) error {
	var microposts []Micropost
	if err := db.Find(&microposts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch microposts"})
	}
	return c.JSON(http.StatusOK, microposts)
}

// @Summary Get a micropost by ID
// @Description Get a micropost by its ID
// @Tags microposts
// @Produce json
// @Param id path int true "Micropost ID"
// @Success 200 {object} Micropost
// @Failure 404 {object} map[string]string
// @Router /microposts/{id} [get]
func getMicropost(c echo.Context, db *gorm.DB) error {
	id := c.Param("id")
	micropost := new(Micropost)

	if err := db.First(micropost, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Micropost not found"})
	}

	return c.JSON(http.StatusOK, micropost)
}

// @Summary Update a micropost
// @Description Update a micropost with the provided title
// @Tags microposts
// @Accept json
// @Produce json
// @Param id path int true "Micropost ID"
// @Param micropost body Micropost true "Micropost object"
// @Success 200 {object} Micropost
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microposts/{id} [put]
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

// @Summary Delete a micropost
// @Description Delete a micropost by its ID
// @Tags microposts
// @Param id path int true "Micropost ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Router /microposts/{id} [delete]
func deleteMicropost(c echo.Context, db *gorm.DB) error {
	id := c.Param("id")

	if err := db.Delete(&Micropost{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete micropost"})
	}

	return c.NoContent(http.StatusNoContent)
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

	// マイグレ��ション
	db.AutoMigrate(&Micropost{})

	e := echo.New()

	// Swagger UIのエンドポイントを追加
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
