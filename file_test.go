package go_logger

import (
	"testing"
	"time"
)

func TestNewAdapterFile(t *testing.T) {
	NewAdapterFile()
}

func TestAdapterFile_Name(t *testing.T) {
	fileAdapter := NewAdapterFile()

	if fileAdapter.Name() != FILE_ADAPTER_NAME {
		t.Error("file adapter name error")
	}
}

func TestAdapterFile_Write(t *testing.T) {

	fileAdapter := NewAdapterFile()

	fileConfig := &FileConfig{
		Filename:      "./test.log",
		LevelFileName: map[int]string{},
		MaxLine:       2000,
		MaxSize:       10000 * 4,
		JsonFormat:    true,
		DateSlice:     "d",
	}
	err := fileAdapter.Init(fileConfig)
	if err != nil {
		t.Fatal(err.Error())
	}

	loggerMsg := &loggerMessage{
		Timestamp:         time.Now().Unix(),
		TimestampFormat:   time.Now().Format("2006-01-02 15:04:05"),
		Millisecond:       time.Now().UnixNano() / 1e6,
		MillisecondFormat: time.Now().Format("2006-01-02 15:04:05.999"),
		Level:             LOGGER_LEVEL_DEBUG,
		LevelString:       "debug",
		Body:              "logger test file adapter write",
		File:              "file_test.go",
		Line:              50,
		Function:          "TestAdapterFile_Write",
	}
	err = fileAdapter.Write(loggerMsg)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAdapterFile_WriteLevelFile(t *testing.T) {

	fileAdapter := NewAdapterFile()

	fileConfig := &FileConfig{
		Filename: "./test.log",
		LevelFileName: map[int]string{
			LOGGER_LEVEL_DEBUG: "./debug.log",
			LOGGER_LEVEL_INFO:  "./info.log",
			LOGGER_LEVEL_ERROR: "./error.log",
		},
		MaxLine:    2000,
		MaxSize:    10000 * 4,
		JsonFormat: true,
		DateSlice:  "d",
	}
	err := fileAdapter.Init(fileConfig)
	if err != nil {
		t.Fatal(err.Error())
	}

	loggerMsg := &loggerMessage{
		Timestamp:         time.Now().Unix(),
		TimestampFormat:   time.Now().Format("2006-01-02 15:04:05"),
		Millisecond:       time.Now().UnixNano() / 1e6,
		MillisecondFormat: time.Now().Format("2006-01-02 15:04:05.999"),
		Level:             LOGGER_LEVEL_DEBUG,
		LevelString:       "debug",
		Body:              "logger test file adapter write",
		File:              "file_test.go",
		Line:              50,
		Function:          "TestAdapterFile_Write",
	}
	fileAdapter.Write(loggerMsg)
	loggerMsg.Level = LOGGER_LEVEL_INFO
	fileAdapter.Write(loggerMsg)
	loggerMsg.Level = LOGGER_LEVEL_ERROR
	fileAdapter.Write(loggerMsg)
}
