package telebot

import (
	"fmt"
	"time"

	"github.com/coke5151/EnvMonitorRealife/internal/database"
	tele "gopkg.in/telebot.v4"
	"gorm.io/gorm"
)

func New(token string, db *gorm.DB) *tele.Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	// commands
	bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})
	bot.Handle("/current_temperature", func(c tele.Context) error {
		record := &database.DetectEnvironment{}
		result := db.Select("temperature", "created_at").Last(record)
		if result.Error != nil {
			return c.Send("取得目前溫度失敗。")
		} else {
			return c.Send(
				fmt.Sprintf("目前的溫度是：%v °C\nTimestamp：%v",
					record.Temperature,
					record.CreatedAt.In(time.FixedZone("Asia/Taipei", 8*60*60)).Format("2006-01-02 15:04:05"),
				),
			)
		}
	})
	bot.Handle("/current_humidity", func(c tele.Context) error {
		record := &database.DetectEnvironment{}
		result := db.Select("humidity_percentage", "created_at").Last(record)
		if result.Error != nil {
			return c.Send("取得目前濕度失敗。")
		} else {
			return c.Send(
				fmt.Sprintf("目前的相對濕度是：%v%%\nTimestamp：%v",
					record.HumidityPercentage,
					record.CreatedAt.In(time.FixedZone("Asia/Taipei", 8*60*60)).Format("2006-01-02 15:04:05"),
				),
			)
		}
	})

	return bot
}
