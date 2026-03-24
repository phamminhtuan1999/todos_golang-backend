package handlers

import (
	"net/http"
	"strings"
	"todos/internal/models"
	"todos/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request" + err.Error(),
			})
			return
		}

		// Validate password length
		if len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password must be at least 8 characters long",
			})
			return
		}

		var hashedPassword string = HashPassword(req.Password)

		user := &models.User{
			Email:        req.Email,
			PasswordHash: hashedPassword,
		}

		createdUser, err := repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				c.JSON(http.StatusConflict, gin.H{
					"message": "User with this email already exists" + err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create user" + err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"data":    createdUser,
		})
	}

}

func HashPassword(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
