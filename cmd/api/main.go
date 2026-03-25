package main

import (
	"log"
	"todos/internal/config"
	"todos/internal/database"
	"todos/internal/handlers"

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

	router.POST("/todos", handlers.CreateTodoHandler(db))

	router.GET("/todos", handlers.GetTodosHandler(db))

	router.GET("/todos/:id", handlers.GetTodoHandler(db))

	router.PUT("/todos/:id", handlers.UpdateTodoHandler(db))

	router.DELETE("/todos/:id", handlers.DeleteTodoHandler(db))

	router.POST("/auth/register", handlers.CreateUserHandler(db))

	router.POST("/auth/login", handlers.LoginHandler(db, cfg))

	router.Run(":" + cfg.Port)
}
