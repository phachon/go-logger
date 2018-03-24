package go_logger

import "testing"

func TestNewLogger(t *testing.T) {
	NewLogger()
}

func TestLogger_Attach(t *testing.T) {

	logger := NewLogger()
	fileConfig := &FileConfig{}
	logger.Attach("file", LOGGER_LEVEL_DEBUG, NewConfigFile(fileConfig))
	outputs := logger.outputs
	for _, outputLogger := range outputs {
		if outputLogger.Name != "file" {
			t.Error("file attach failed")
		}
	}
}

func TestLogger_Detach(t *testing.T) {

	logger := NewLogger()
	logger.Detach("console")

	outputs := logger.outputs

	if len(outputs) > 0 {
		t.Error("logger detach error")
	}
}

func TestLogger_LoggerLevel(t *testing.T) {

	logger := NewLogger()

	level := logger.LoggerLevel("emerGENCY")
	if level != LOGGER_LEVEL_EMERGENCY {
		t.Error("logger loggerLevel error")
	}
	level = logger.LoggerLevel("ALERT")
	if level != LOGGER_LEVEL_ALERT {
		t.Error("logger loggerLevel error")
	}
	level = logger.LoggerLevel("cRITICAL")
	if level != LOGGER_LEVEL_CRITICAL {
		t.Error("logger loggerLevel error")
	}
	level = logger.LoggerLevel("DEBUG")
	if level != LOGGER_LEVEL_DEBUG {
		t.Error("logger loggerLevel error")
	}
	level = logger.LoggerLevel("ooox")
	if level != LOGGER_LEVEL_DEBUG {
		t.Error("logger loggerLevel error")
	}
}