package main

import (
	"go-logger"
)

func main()  {

	logger := go_logger.NewLogger()
	//default attach console, detach console
	logger.Detach("console")

	console := &go_logger.ConsoleConfig{
		Color: true,
		JsonFormat: true,
	}

	logger.Attach("console", go_logger.NewConfigConsole(console))

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