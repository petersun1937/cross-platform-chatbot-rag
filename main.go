package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/handlers"
	"crossplatform_chatbot/server"
	"crossplatform_chatbot/service"
)

func main() {

	// Initialize config (only once)
	conf := config.GetConfig()

	// Initialize database
	db := database.NewDatabase(conf)
	if err := db.Init(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	// Initialize service
	svc := service.NewService(conf.BotConfig, &conf.EmbeddingConfig, db)
	if err := svc.RunBots(); err != nil {
		log.Fatal("Failed to run the Bot services:", err)
	}
	//svc := service.NewService(db)

	// Initialize handler
	handler := handlers.NewHandler(svc)

	// initialize bots
	// bots := createBots(conf, svc)

	// Initialize the app (app acts as the central hub for the application, holds different initialized values)
	//app := server.NewApp(conf, svc, bots)
	app := server.NewApp(*conf, handler)
	if err := app.Run(); err != nil {
		log.Fatal("Failed to run the app:", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1) // creates a channel named quit

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // tells the program to listen for specific signals (SIGINT and SIGTERM) and send them to the quit channel.
	<-quit                                               // channel receive operation; blocking/waiting until a signal is received in quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // create context with timeout
	defer cancel()                                                          // ensure the context is canceled when the function exists

	if err := app.Server.Shutdown(ctx); err != nil { // graceful shutdown
		log.Fatal("Server Shutdown: ", err)
	}

	fmt.Println("Server exiting")

}
