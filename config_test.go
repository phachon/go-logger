package go_logger

import "testing"

func TestNewConfigConsole(t *testing.T) {
	consoleConfig := &ConsoleConfig{}
	config := NewConfigConsole(consoleConfig)
	if config.Console != consoleConfig {
		t.Error("new config console error")
	}
}

func TestNewConfigFile(t *testing.T) {
	fileConfig := &FileConfig{}
	config := NewConfigFile(fileConfig)
	if config.File != fileConfig {
		t.Error("new config file error")
	}
}

func TestNewConfigApi(t *testing.T) {
	apiConfig := &ApiConfig{}
	config := NewConfigApi(apiConfig)
	if config.Api != apiConfig {
		t.Error("new config api error")
	}
}