package main

import (
	"log"
	"todos/internal/config"
	"todos/internal/database"
	"todos/internal/handlers"
	"todos/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var err error
	var db *pgxpool.Pool
	var cfg *config.Config

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.Close()

	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "pong",
			"status":   "success",
			"database": "connected",
			"data":     nil,
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(db))

	router.POST("/auth/login", handlers.LoginHandler(db, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", handlers.CreateTodoHandler(db))

		protected.GET("", handlers.GetTodosHandler(db))

		protected.GET("/:id", handlers.GetTodoHandler(db))

		protected.PUT("/:id", handlers.UpdateTodoHandler(db))

		protected.DELETE("/:id", handlers.DeleteTodoHandler(db))
	}

	router.Run(":" + cfg.Port)
}
