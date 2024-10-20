package server

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type App struct { //TODO: app or GetConfig for config?
	Config  *config.Config
	Service *service.Service // for database operations
	Router  *gin.Engine
	Bots    map[string]bot.Bot
	//LineBot    bot.LineBot
	//TgBot      bot.TgBot
	//FbBot      bot.FbBot
	//GeneralBot bot.GeneralBot // custom frontend
}

func (a App) Run() error {
	// running bots
	for _, bot := range a.Bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}

	// initialize http server routes from app struct
	go a.RunRoutes(a.Config, a.Service)

	return nil
}

func NewApp(conf *config.Config, srv *service.Service, bots map[string]bot.Bot) *App {

	return &App{
		Config:  conf,
		Service: srv,
		Router:  gin.Default(),
		Bots:    bots,
		// LineBot:    lineBot,    // Store the initialized LineBot
		// TgBot:      tgBot,      // Store the initialized TgBot
		// FbBot:      fbBot,      // Store the initialized FbBot
		// GeneralBot: generalBot, // Store the GeneralBot for /api/message route
	}
}
