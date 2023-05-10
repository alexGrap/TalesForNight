package main

import (
	"fsm/config"
	"fsm/internal/handlers"
	"fsm/pkg/repository"
	"github.com/joho/godotenv"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM"),
		Poller: &tele.LongPoller{Timeout: 60 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	viperConf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.ParseConfig(viperConf)
	if err != nil {
		log.Fatal(err)
	}
	repository.Connection.Database, err = repository.InitPsqlDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	bot.Use(middleware.AutoRespond())
	bGroup := bot.Group()
	repository.InitTables()
	storage := memory.NewStorage()
	defer storage.Close()
	manager := fsm.NewManager(bGroup, memory.NewStorage())
	handlers.StartHandlers(bGroup, manager)
	bot.Start()
}
