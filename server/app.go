package server

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/handlers"
)

type App struct { //TODO: app or GetConfig for config?
	Config config.Config // *
	Server *Server
	// Service *service.Service // for database operations
	// Router  *gin.Engine
	// Bots    map[string]bot.Bot
	// Server  *Server
	//LineBot    bot.LineBot
	//TgBot      bot.TgBot
	//FbBot      bot.FbBot
	//GeneralBot bot.GeneralBot // custom frontend
}

func (a App) Run() error {
	// initialize http server routes from app struct
	a.Server.Start()

	// for _, bot := range a.Bots {
	// 	if err := bot.Run(); err != nil {
	// 		// log.Fatal("running bot failed:", err)
	// 		fmt.Printf("running bot failed: %s", err.Error())
	// 		return err
	// 	}
	// }

	// // initialize http server routes from app struct
	// go a.RunRoutes(a.Config, a.Service, *a.Server)

	return nil
}

func NewApp(conf config.Config, handler *handlers.Handler) *App {

	// svrcfg := config.ServerConfig{
	// 	Host: os.Getenv("HOST"),
	// 	//Port:    8080, // Default port, can be overridden
	// 	Port:    conf.Port,
	// 	Timeout: 30 * time.Second,
	// 	MaxConn: 100,
	// }

	// svr := New(svrcfg, svc, conf)
	// //svr := New(cfg)

	return &App{
		Config: conf,
		Server: New(conf.ServerConfig, handler),
		// Config:  conf,
		// Service: svc,
		// Router:  gin.Default(),
		// Bots:    bots,
		// Server:  svr,
		// // LineBot:    lineBot,    // Store the initialized LineBot
		// // TgBot:      tgBot,      // Store the initialized TgBot
		// // FbBot:      fbBot,      // Store the initialized FbBot
		// // GeneralBot: generalBot, // Store the GeneralBot for /api/message route
	}
}
