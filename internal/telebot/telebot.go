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

	// 新增風扇狀態查詢命令
	bot.Handle("/fanstatus", func(c tele.Context) error {
		status, err := database.GetLatestFanStatus(db)
		if err != nil {
			return c.Send("取得風扇狀態失敗。")
		}

		statusText := map[string]string{
			"off":              "關閉",
			"clockwise":        "正轉",
			"counterclockwise": "逆轉",
		}[status.Status]

		return c.Send(fmt.Sprintf(
			"目前風扇狀態：%s\n更新時間：%v",
			statusText,
			status.CreatedAt.In(time.FixedZone("Asia/Taipei", 8*60*60)).Format("2006-01-02 15:04:05"),
		))
	})

	// 新增風扇狀態設定命令
	bot.Handle("/setfan", func(c tele.Context) error {
		args := c.Args()
		if len(args) != 1 {
			return c.Send("使用方式：/setfan [off|clockwise|counterclockwise]")
		}

		status := args[0]
		validStatus := map[string]bool{
			"off":              true,
			"clockwise":        true,
			"counterclockwise": true,
		}

		if !validStatus[status] {
			return c.Send("無效的風扇狀態，請使用 off、clockwise 或 counterclockwise")
		}

		if err := database.InsertFanStatus(db, status); err != nil {
			return c.Send("設定風扇狀態失敗")
		}

		statusText := map[string]string{
			"off":              "關閉",
			"clockwise":        "正轉",
			"counterclockwise": "逆轉",
		}[status]

		return c.Send(fmt.Sprintf("風扇狀態已設定為：%s", statusText))
	})

	return bot
}
