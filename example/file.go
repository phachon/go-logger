package main

import (
	"go-logger"
)

func main()  {

	logger := go_logger.NewLogger()

	fileConfig := &go_logger.FileConfig{
		Filename : "./test.log",
		MaxSize : 1024 * 1024,
		MaxLine : 100000,
		DateSlice : "d",
	}
	logger.Attach("file", go_logger.NewConfigFile(fileConfig))
	logger.SetLevel(go_logger.LOGGER_LEVEL_DEBUG)
	logger.SetAsync()

	i := 0
	for  {
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

		logger.Emergency("this is a emergency log!")
		logger.Alert("this is a alert log!")
		logger.Critical("this is a critical log!")

		i += 1
		if i == 100 {
			break
		}
	}

	logger.Flush()
}