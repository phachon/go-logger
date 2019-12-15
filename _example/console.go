package main

import (
	"github.com/phachon/go-logger"
)

func main() {

	logger := go_logger.NewLogger()
	//default attach console, detach console
	logger.Detach("console")

	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}

	logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)

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

	logger.Emergencyf("this is a emergency %d log!", 10)
	logger.Alertf("this is a alert %s log!", "format")
	logger.Criticalf("this is a critical %s log!", "format")
	logger.Errorf("this is a error %s log!", "format")
	logger.Warningf("this is a warning %s log!", "format")
	logger.Noticef("this is a notice %s log!", "format")
	logger.Infof("this is a info %s log!", "format")
	logger.Debugf("this is a debug %s log!", "format")

	logger.Flush()
}
