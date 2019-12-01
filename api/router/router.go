package router

import (
	"gochat/api/handler"
	"gochat/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register ...
func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware())
	initUserRouter(r)
	r.NoRoute(func(c *gin.Context) {
		tools.FailWithMsg(c, "please check request url !")
	})
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/register", handler.Register)
}

// CorsMiddleware ...
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		var openCorsFlag = true
		if openCorsFlag {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}
		c.Next()
	}
}
