package http

import (
	"errors"
	"log"
	"net/http"
	"simpleMessenger/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	id := c.Query("id")
	login := c.Query("login")
	search := c.Query("search")
	limitStr := c.DefaultQuery("limit", "20")

	if id != "" {
		idVal, err := strconv.Atoi(id)

		if err != nil {
			log.Printf("Invalid ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		user, err := h.userService.GetUserByID(uint(idVal))

		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				log.Printf("User not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Printf("Error getting user by ID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, user)
		return
	}

	if login != "" {
		user, err := h.userService.GetUserByLogin(login)

		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				log.Printf("User not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			log.Printf("Error getting user by login: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, user)
		return
	}

	if search != "" {
		limit, err := strconv.Atoi(limitStr)

		if err != nil {
			log.Printf("Invalid limit parameter: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}

		users, err := h.userService.SearchUsersByLogin(search, limit)

		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				c.JSON(http.StatusOK, gin.H{"users": users})
				return
			}
			log.Printf("Error searching users: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "id, login or search query parameter is required"})
}
