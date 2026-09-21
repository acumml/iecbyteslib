package iecbytes

import (
	"strings"
	"testing"
)

func FuzzFormat(f *testing.F) {
	for _, seed := range []uint64{
		0,
		1,
		1023,
		1024,
		(1 << 20) - 1,
		1 << 60,
		^uint64(0),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, bytes uint64) {
		got := Format(bytes)
		if got == "" {
			t.Fatal("expected non-empty output")
		}
		if !hasKnownUnit(got) {
			t.Fatalf("unknown unit in %q", got)
		}
	})
}

func FuzzFormatWithOptions(f *testing.F) {
	f.Add(uint64(1536), int(2), true, true)
	f.Add(^uint64(0), int(10), false, true)
	f.Add(uint64(0), int(-999), true, false)
	f.Add(uint64(1024), int(999), false, false)

	f.Fuzz(func(t *testing.T, bytes uint64, precision int, trim, promote bool) {
		got := FormatWithOptions(bytes, Options{
			Precision:         precision,
			TrimTrailingZeros: trim,
			PromoteOnRounding: promote,
		})
		if got == "" {
			t.Fatal("expected non-empty output")
		}
		if !hasKnownUnit(got) {
			t.Fatalf("unknown unit in %q", got)
		}
	})
}

func hasKnownUnit(value string) bool {
	for _, unit := range unitNames {
		if strings.HasSuffix(value, " "+unit) {
			return true
		}
	}
	return false
}
