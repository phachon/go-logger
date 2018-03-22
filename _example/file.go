package main

import (
	"github.com/phachon/go-logger"
)

func main()  {

	logger := go_logger.NewLogger()

	fileConfig := &go_logger.FileConfig{
		Filename : "./test.log",
		LevelFileName : map[int]string{
			go_logger.LOGGER_LEVEL_ERROR: "./error.log",
			go_logger.LOGGER_LEVEL_INFO: "./info.log",
			go_logger.LOGGER_LEVEL_DEBUG: "./debug.log",
		},
		MaxSize : 1024 * 1024,
		MaxLine : 10000,
		DateSlice : "d",
		JsonFormat: true,
	}
	logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, go_logger.NewConfigFile(fileConfig))
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
		if i == 21000 {
			break
		}
	}

	logger.Flush()
}