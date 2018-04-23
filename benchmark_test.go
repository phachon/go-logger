package go_logger

import (
	"testing"
)

// go test -run=benchmark -bench=. -benchtime="3s"

func BenchmarkLoggerConsole(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark logger message")
		}
	})
}
