package iecbytes

import (
	"log/slog"
	"testing"
)

func TestAttr(t *testing.T) {
	t.Parallel()

	attr := Attr("used_memory", 5_000_000)
	if attr.Key != "used_memory" {
		t.Fatalf("key = %q, want %q", attr.Key, "used_memory")
	}
	if attr.Value.Kind() != slog.KindString {
		t.Fatalf("kind = %v, want %v", attr.Value.Kind(), slog.KindString)
	}
	if got, want := attr.Value.String(), "4.77 MiB"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestGroupAttr(t *testing.T) {
	t.Parallel()

	const bytes = uint64(1_234_567_890)
	attr := GroupAttr("memory", bytes)
	if attr.Key != "memory" {
		t.Fatalf("key = %q, want %q", attr.Key, "memory")
	}
	if attr.Value.Kind() != slog.KindGroup {
		t.Fatalf("kind = %v, want %v", attr.Value.Kind(), slog.KindGroup)
	}

	group := attr.Value.Group()
	if len(group) != 2 {
		t.Fatalf("group len = %d, want 2", len(group))
	}
	if group[0].Key != "bytes" || group[0].Value.Uint64() != bytes {
		t.Fatalf("raw bytes attr = %#v", group[0])
	}
	if group[1].Key != "iec" || group[1].Value.String() != "1.15 GiB" {
		t.Fatalf("IEC attr = %#v", group[1])
	}
}
