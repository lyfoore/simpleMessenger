package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	engine := gin.Default()
	return &Router{engine: engine}
}

func (r *Router) SetupRouter() {
	r.engine.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
}

func (r *Router) Run() {
	r.engine.Run(":8080")
	fmt.Println("Server is running on port 8080")
}
