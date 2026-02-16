package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simpleMessenger/internal/service"
	"simpleMessenger/internal/transport/websocket"
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
	wsHub *websocket.Hub,
) {
	public := r.engine.Group("/api")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	protected := r.engine.Group("/api")
	protected.Use(AuthMiddleware(tokenService))
	{
		protected.GET("/chats", chatHandler.GetChats) // query: limit
		protected.POST("/chats", chatHandler.CreateChat)
		protected.DELETE("/chats/:chatId", chatHandler.DeleteChat)

		protected.GET("/chats/:chatId/messages", messageHandler.GetMessages) // query: limit

		protected.POST("/chats/:chatId/messages", messageHandler.SendMessage)

		protected.DELETE("/messages/:messageId", messageHandler.DeleteMessage)

		protected.GET("/ws", func(c *gin.Context) {
			websocket.ServeWs(wsHub, c.Writer, c.Request)
		})

		protected.GET("/me", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(http.StatusOK, gin.H{
				"message": "You are signed in",
				"user_id": userID,
			})
		})
	}

	r.engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Messenger API is running")
	})
}

func (r *Router) Run() {
	fmt.Println("Server is running on port 8080")
	r.engine.Run(":8080")
}
