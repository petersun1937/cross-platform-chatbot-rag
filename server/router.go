package server

import (
	"crossplatform_chatbot/handlers"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func (s *Server) InitRoutes(r *gin.Engine, conf *config.Config, srv *service.Service) {
func (s *Server) InitRoutes(handler *handlers.Handler) {
	// Set up logging to a file (bot.log)
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = file

	// Enable CORS
	s.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Use "*" to allow all origins
		//AllowOrigins:     []string{"http://localhost:3000"}, // localhost needs to be specified directly
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware to log requests
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// Define routes
	s.router.POST("/line/webhook", handler.HandleLineWebhook)
	s.router.POST("/telegram/webhook", handler.HandleTelegramWebhook)
	s.router.GET("/messenger/webhook", handler.VerifyMessengerWebhook) // For webhook verification
	s.router.POST("/messenger/webhook", handler.HandleMessengerWebhook)
	s.router.GET("/instagram/webhook", handler.VerifyMessengerWebhook) // For webhook verification
	s.router.POST("/instagram/webhook", handler.HandleMessengerWebhook)
	s.router.POST("/api/message", handler.HandlerGeneralBot)
	s.router.POST("/api/document/upload", handler.HandlerDocumentUpload)

	//r.POST("/login", handlers.Login)

	// Protected routes
	/*authorized := r.Group("/api")
	authorized.Use(middleware.JWTMiddleware())
	{
		authorized.POST("/message", handlers.HandleCustomMessage)
		// Add other protected routes here
	}*/

	fmt.Println("Server routes initialized")
	//fmt.Println("Server started")
	//r.Run(":8080")
}
