package router

import (
	"net/http"

	"github.com/varopxndx/chat/middleware"
	"github.com/varopxndx/chat/token"

	"github.com/gin-gonic/gin"
)

// Handler methods
type Handler interface {
	Login(g *gin.Context)
	SignUp(g *gin.Context)
	ProcessMessage(g *gin.Context)
	Message(g *gin.Context)
}

// Setup creates the router and prepares the endpoints
func Setup(handler Handler, tokenService *token.TokenService) *gin.Engine {
	router := gin.New()

	// frontend routes
	router.LoadHTMLGlob("public/*.html")
	v1 := router.Group("/v1")
	v1.Static("/assets", "./public/assets")
	v1.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	v1.GET("/signup.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})
	v1.GET("/chat.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	// backend routes
	v1.POST("/login", handler.Login)
	v1.POST("/signup", handler.SignUp)

	// rabbit broker routes
	auth := v1.Group("/")
	auth.Use(middleware.AuthToken(tokenService.ValidateToken))
	auth.GET("/message", handler.Message)
	auth.GET("/ws", handler.ProcessMessage)

	return router
}
