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

func (r *Router) SetupRouter(
	authHandler *AuthHandler,
	chatHandler *ChatHandler,
	messageHandler *MessageHandler,
	tokenService service.TokenService,
) {
	public := r.engine.Group("/api")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	protected := r.engine.Group("/api")
	protected.Use(AuthMiddleware(tokenService))
	{
		//protected.GET("/chats", chatHandler.GetChats) // query: limit
		protected.POST("/chats", chatHandler.CreateChat)
		protected.DELETE("/chats/:chatId", chatHandler.DeleteChat)

		//protected.GET("/chats/:chatId/messages", messageHandler.GetMessages) // query: limit

		protected.POST("/chats/:chatId/messages", messageHandler.SendMessage)

		protected.DELETE("/messages/:messageId", messageHandler.DeleteMessage)

		protected.GET("/me", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{
				"message": "You are signed in",
				"user_id": userID,
			})
		})
	}

	r.engine.GET("/", func(c *gin.Context) {
		c.String(200, "Messenger API is running")
	})
}

func (r *Router) Run() {
	r.engine.Run(":8080")
	fmt.Println("Server is running on port 8080")
}
