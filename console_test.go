package go_logger

import (
	"testing"
	"time"
)

func TestNewAdapterConsole(t *testing.T) {
	NewAdapterConsole()
}

func TestAdapterConsole_Init(t *testing.T) {

	consoleAdapter := NewAdapterConsole()

	consoleConfig := &ConsoleConfig{}
	consoleAdapter.Init(NewConfigConsole(consoleConfig))
}

func TestAdapterConsole_Name(t *testing.T) {

	consoleAdapter := NewAdapterConsole()

	if consoleAdapter.Name() != CONSOLE_ADAPTER_NAME {
		t.Error("consoleAdapter name error")
	}
}

func TestAdapterConsole_WriteColor(t *testing.T) {

	consoleAdapter := NewAdapterConsole()

	consoleConfig := &ConsoleConfig{
		Color:true,
	}
	consoleAdapter.Init(NewConfigConsole(consoleConfig))

	loggerMsg := &loggerMessage {
		Timestamp : time.Now().Unix(),
		TimestampFormat : time.Now().Format("2006-01-02 15:04:05"),
		Millisecond : time.Now().UnixNano()/1e6,
		MillisecondFormat : time.Now().Format("2006-01-02 15:04:05.999"),
		Level : LOGGER_LEVEL_DEBUG,
		LevelString: "debug",
		Body: "logger console adapter test color",
		File : "console_test.go",
		Line : 50,
		Function: "TestAdapterConsole_WriteIsColor",
	}
	err := consoleAdapter.Write(loggerMsg)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAdapterConsole_WriteJsonFormat(t *testing.T) {

	consoleAdapter := NewAdapterConsole()

	consoleConfig := &ConsoleConfig{
		JsonFormat:true,
	}
	consoleAdapter.Init(NewConfigConsole(consoleConfig))

	loggerMsg := &loggerMessage {
		Timestamp : time.Now().Unix(),
		TimestampFormat : time.Now().Format("2006-01-02 15:04:05"),
		Millisecond : time.Now().UnixNano()/1e6,
		MillisecondFormat : time.Now().Format("2006-01-02 15:04:05.999"),
		Level : LOGGER_LEVEL_DEBUG,
		LevelString: "debug",
		Body: "logger console adapter test jsonFormat",
		File : "console_test.go",
		Line : 77,
		Function: "TestAdapterConsole_WriteJsonFormat",
	}
	err := consoleAdapter.Write(loggerMsg)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAdapterConsole_WriteShowFileLine(t *testing.T) {

	consoleAdapter := NewAdapterConsole()

	consoleConfig := &ConsoleConfig{
		ShowFileLine:true,
	}
	consoleAdapter.Init(NewConfigConsole(consoleConfig))

	loggerMsg := &loggerMessage {
		Timestamp : time.Now().Unix(),
		TimestampFormat : time.Now().Format("2006-01-02 15:04:05"),
		Millisecond : time.Now().UnixNano()/1e6,
		MillisecondFormat : time.Now().Format("2006-01-02 15:04:05.999"),
		Level : LOGGER_LEVEL_DEBUG,
		LevelString: "info",
		Body: "logger console adapter test show file line",
		File : "console_test.go",
		Line : 104,
		Function: "TestAdapterConsole_WriteJsonShowFileLine",
	}
	err := consoleAdapter.Write(loggerMsg)
	if err != nil {
		t.Error(err.Error())
	}
}