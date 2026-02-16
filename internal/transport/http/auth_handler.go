package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleMessenger/internal/model"
	"simpleMessenger/internal/service"
)

type LoginRequest struct {
	Username string `json:"username"`
}

type RegisterRequest struct {
	Username string `json:"username"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	req := &RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("failed to unmarshal register request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad register request"})
		return
	}

	err := h.authService.Register(&model.User{Login: req.Username})
	if err != nil {
		log.Printf("failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to register user"})
		return
	}

	response, err := h.authService.Login(req.Username)
	if err != nil {
		log.Printf("failed to login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to login"})
		return
	}
	c.JSON(http.StatusOK, response)
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	req := &LoginRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("failed to unmarshal login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad login request"})
		return
	}

	if req.Username == "" {
		log.Printf("failed to unmarshal login request: empty username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is empty"})
		return
	}

	response, err := h.authService.Login(req.Username)
	if err != nil {
		log.Printf("failed to login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to login"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func AuthMiddleware(tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			log.Printf("no token provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		userID, err := tokenService.VerifyToken(token)
		if err != nil {
			log.Printf("failed to verify token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		ctx := context.WithValue(c.Request.Context(), "user_id", userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserIdFromContext(c *gin.Context) (uint, error) {
	value, exists := c.Get("user_id")

	if !exists {
		return 0, fmt.Errorf("user is not authenticated")
	}

	userID, ok := value.(uint)

	if !ok {
		return 0, fmt.Errorf("invalid user id type")
	}

	return userID, nil
}
