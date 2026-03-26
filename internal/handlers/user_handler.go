package handlers

import (
	"net/http"
	"strings"
	"time"
	"todos/internal/config"
	"todos/internal/models"
	"todos/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
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

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request" + err.Error(),
			})
			return
		}

		user, err := repository.GetUserByEmail(pool, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password" + err.Error(),
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid credentials" + err.Error(),
			})
			return
		}

		// Generate JWT token
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate token" + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"data": LoginResponse{
				Token: tokenString,
			},
		})
	}
}

func TestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found in context",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Test successful",
			"user_id": user_id,
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
