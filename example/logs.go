package main

import (
	"go-logger"
)

func main()  {

	logger := go_logger.NewLogger()

	//logger.Attach("console", map[string]interface{}{
	//	"color": false,
	//})
	logger.Attach("file", map[string]interface{}{
		"filename": "test.log",
		"maxSize": 5,
		"maxLine": 77,
	})
	logger.SetLevel(go_logger.LOGGER_LEVEL_DEBUG)
	logger.SetAsync()

	logger.Emergency("this is a emergency log!")
	logger.Alert("this is a alert log!")
	logger.Critical("this is a critical log!")
	logger.Error("this is a error log!")
	logger.Warning("this is a warning log!")
	logger.Notice("this is a notice log!")
	logger.Info("this is a info log!")
	logger.Debug("this is a debug log!")

	logger.Emergency("this is a emergency log!")
	logger.Notice("this is a notice log!")
	logger.Info("this is a info log!")
	logger.Debug("this is a debug log!")

	logger.Flush()
}