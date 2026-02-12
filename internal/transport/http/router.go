package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"simpleMessenger/internal/service"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	engine := gin.Default()
	return &Router{engine: engine}
}

func (r *Router) SetupRouter(authHandler *AuthHandler, tokenService service.TokenService) {
	r.engine.POST("/api/auth/register", authHandler.Register)
	r.engine.POST("/api/auth/login", authHandler.Login)

	protected := r.engine.Group("/api")
	protected.Use(AuthMiddleware(tokenService))
	{
		protected.GET("/me", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{
				"message": "You are signed in",
				"user_id": userID,
			})
		})
	}

	r.engine.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
}

func (r *Router) Run() {
	r.engine.Run(":8080")
	fmt.Println("Server is running on port 8080")
}
