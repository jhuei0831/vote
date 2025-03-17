package middleware

import (
	"github.com/gin-gonic/gin"
	"time"

	"vote/app/utils"
)

// LoggerToFile 是一個 Gin 中介函數，用於將請求日誌記錄到文件中。
// 它會記錄請求的開始時間和結束時間，計算請求的延遲時間，
// 並記錄請求的方法、URI、狀態碼和客戶端 IP 地址。
func LoggerToFile() gin.HandlerFunc {
	logger := utils.Logger()
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.Infof("| %3d | %13v | %15s | %s | %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}