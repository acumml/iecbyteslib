package iecbytes

import "testing"

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Format(1_234_567_890_123)
	}
}

func BenchmarkAppend(b *testing.B) {
	buf := make([]byte, 0, 32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = Append(buf, 1_234_567_890_123)
	}
}
