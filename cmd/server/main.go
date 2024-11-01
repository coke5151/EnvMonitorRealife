package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	router.Use(cors.Default())

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	go func() {
		router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()
	b.Start()
}
