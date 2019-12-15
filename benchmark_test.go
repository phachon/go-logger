package go_logger

import (
	"testing"
)

// go test -run=benchmark -cpu=1,2,4 -benchmem -benchtime=3s -bench="ConsoleText"
func BenchmarkLoggerConsoleText(b *testing.B) {
	logger := NewLogger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
}

// go test -run=benchmark -cpu=1,2,4 -benchmem -benchtime=3s -bench="ConsoleAsyncText"
func BenchmarkLoggerConsoleAsyncText(b *testing.B) {
	logger := NewLogger()
	logger.SetAsync()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
	logger.Flush()
}

// go test -run=benchmark -cpu=1,2,4 -benchmem -benchtime=3s -bench="ConsoleJson"
func BenchmarkLoggerConsoleJson(b *testing.B) {
	logger := NewLogger()
	logger.Detach("console")
	logger.Attach("console", LOGGER_LEVEL_DEBUG, &ConsoleConfig{
		JsonFormat: true,
	})
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
}

// go test -run=benchmark -cpu=1,2,4 -benchmem -benchtime=3s -bench="FileText"
func BenchmarkLoggerFileText(b *testing.B) {
	logger := NewLogger()
	logger.Detach("console")
	logger.Attach("file", LOGGER_LEVEL_DEBUG, &FileConfig{
		Filename:  "./test.log",
		DateSlice: "d",
	})
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
}

// go test -run=benchmark -cpu=1,2,4 -benchmem -benchtime=3s -bench="AsyncText"
func BenchmarkLoggerFileAsyncText(b *testing.B) {
	logger := NewLogger()
	logger.Detach("console")
	logger.Attach("file", LOGGER_LEVEL_DEBUG, &FileConfig{
		Filename:  "./test.log",
		DateSlice: "d",
	})
	logger.SetAsync()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
	logger.Flush()
}

// go test -run=benchmark -cpu=1,2,4 -benchmem -benchtime=3s -bench="FileJson"
func BenchmarkLoggerFileJson(b *testing.B) {
	logger := NewLogger()
	logger.Detach("console")
	logger.Attach("file", LOGGER_LEVEL_DEBUG, &FileConfig{
		Filename:   "./test.log",
		DateSlice:  "d",
		JsonFormat: true,
	})
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
}
