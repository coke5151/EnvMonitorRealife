package api

import (
	"net/http"

	"github.com/coke5151/EnvMonitorRealife/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type detectEnvironment struct {
	Temperature        float64 `json:"temperature" binding:"required"`
	HumidityPercentage float64 `json:"humidity_percentage" binding:"required"`
	FanStatus          string  `json:"fan_status" binding:"required"`
}

func New(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// get
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// post
	router.POST("/insert", func(ctx *gin.Context) {
		json := &detectEnvironment{}

		// check if there is an binding error (eg. typo)
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "JSON 格式錯誤",
				"details": err.Error(),
			})
			return
		}

		// temperature check
		if json.Temperature < -50 || json.Temperature > 100 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "溫度數值超出合理範圍 (-50°C 到 100°C)",
			})
			return
		}

		// humidity check
		if json.HumidityPercentage < 0 || json.HumidityPercentage > 100 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "濕度百分比必須在 0-100 之間",
			})
			return
		}

		// write into database
		if err := db.Create(&database.DetectEnvironment{
			Temperature:        json.Temperature,
			HumidityPercentage: json.HumidityPercentage,
		}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "資料儲存失敗",
			})
		}

		if err := database.InsertFanStatus(db, json.FanStatus); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "風扇狀態儲存失敗",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "資料新增成功",
			"data":    json,
		})
	})

	// 新增取得風扇狀態的端點
	router.GET("/fan/status", func(ctx *gin.Context) {
		status, err := database.GetLatestFanStatus(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "無法取得風扇狀態",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": status.Status,
		})
	})

	router.PUT("/fan/status", func(ctx *gin.Context) {
		var req struct {
			Status string `json:"status" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON 格式錯誤",
			})
			return
		}

		if req.Status != "off" && req.Status != "clockwise" && req.Status != "counterclockwise" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "無效的風扇狀態",
			})
			return
		}

		if err := database.InsertFanStatus(db, req.Status); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "更新風扇狀態失敗",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "風扇狀態更新成功",
			"status":  req.Status,
		})
	})

	router.Use(cors.Default())
	return router
}
