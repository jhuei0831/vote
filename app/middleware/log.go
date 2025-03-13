package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// Logger 建立並返回一個配置好的 logrus.Logger 實例。
// 它會在當前工作目錄下創建一個 logs 目錄，並在其中生成一個以當前日期命名的日誌文件。
// 日誌文件會以追加模式打開，並設置日誌級別為 Debug。
// 日誌的時間戳格式為 "2006-01-02 15:04:05"。
func Logger() *logrus.Logger {
	// 取得當前時間
	now := time.Now()
	// 設定日誌文件路徑
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	// 創建 logs 目錄
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	// 設定日誌文件名稱，格式為 "2006-01-02.log"
	logFileName := now.Format("2006-01-02") + ".log"

	// 組合日誌文件完整路徑
	fileName := path.Join(logFilePath, logFileName)
	// 檢查日誌文件是否存在，不存在則創建
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}

	// 以追加模式打開日誌文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	// 創建一個新的 logrus.Logger 實例
	logger := logrus.New()
	// 設置日誌輸出到文件
	logger.Out = src
	// 設置日誌級別為 Debug
	logger.SetLevel(logrus.DebugLevel)
	// 設置日誌格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp: true,
	})

	// 返回配置好的 logger 實例
	return logger
}

// LoggerToFile 是一個 Gin 中介函數，用於將請求日誌記錄到文件中。
// 它會記錄請求的開始時間和結束時間，計算請求的延遲時間，
// 並記錄請求的方法、URI、狀態碼和客戶端 IP 地址。
func LoggerToFile() gin.HandlerFunc {
	logger := Logger()
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