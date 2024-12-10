package main

import (
	"log"
	"os"

	"github.com/coke5151/EnvMonitorRealife/internal/api"
	"github.com/coke5151/EnvMonitorRealife/internal/database"
	"github.com/coke5151/EnvMonitorRealife/internal/telebot"
	"github.com/joho/godotenv"
)

func main() {
	db := database.New()

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	token := os.Getenv("TELEGRAM_TOKEN")
	bot := telebot.New(token, db)
	router := api.New(db)

	go func() {
		router.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()
	bot.Start()
}
