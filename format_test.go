package iecbytes

import (
	"math"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		bytes uint64
		want  string
	}{
		{name: "zero", bytes: 0, want: "0 B"},
		{name: "bytes", bytes: 500, want: "500 B"},
		{name: "maximum bytes unit", bytes: 1023, want: "1023 B"},
		{name: "one kibibyte", bytes: 1 << 10, want: "1.00 KiB"},
		{name: "one mebibyte", bytes: 1 << 20, want: "1.00 MiB"},
		{name: "one gibibyte", bytes: 1 << 30, want: "1.00 GiB"},
		{name: "one tebibyte", bytes: 1 << 40, want: "1.00 TiB"},
		{name: "one pebibyte", bytes: 1 << 50, want: "1.00 PiB"},
		{name: "one exbibyte", bytes: 1 << 60, want: "1.00 EiB"},
		{name: "common storage value", bytes: 5_000_000, want: "4.77 MiB"},
		{name: "100 gibibytes", bytes: 100 << 30, want: "100.00 GiB"},
		{name: "decimal trillion bytes", bytes: 1_234_567_890_123, want: "1.12 TiB"},
		{name: "maximum uint64", bytes: math.MaxUint64, want: "16.00 EiB"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Format(tt.bytes); got != tt.want {
				t.Fatalf("Format(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		bytes uint64
		want  string
	}{
		{1024 - 1, "1023 B"},
		{1024, "1.00 KiB"},
		{1024 + 1, "1.00 KiB"},
		{(1 << 20) - 1, "1.00 MiB"},
		{1 << 20, "1.00 MiB"},
		{(1 << 20) + 1, "1.00 MiB"},
		{(1 << 30) - 1, "1.00 GiB"},
		{1 << 30, "1.00 GiB"},
		{(1 << 40) - 1, "1.00 TiB"},
		{1 << 40, "1.00 TiB"},
		{(1 << 50) - 1, "1.00 PiB"},
		{1 << 50, "1.00 PiB"},
		{(1 << 60) - 1, "1.00 EiB"},
		{1 << 60, "1.00 EiB"},
	}

	for _, tt := range tests {
		if got := Format(tt.bytes); got != tt.want {
			t.Errorf("Format(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestFormatWithOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		bytes uint64
		opts  Options
		want  string
	}{
		{
			name:  "zero decimal places rounds half up",
			bytes: 1536,
			opts:  Options{Precision: 0},
			want:  "2 KiB",
		},
		{
			name:  "one decimal place",
			bytes: 1536,
			opts:  Options{Precision: 1},
			want:  "1.5 KiB",
		},
		{
			name:  "four decimal places",
			bytes: 1536,
			opts:  Options{Precision: 4},
			want:  "1.5000 KiB",
		},
		{
			name:  "trim trailing zeros",
			bytes: 1536,
			opts: Options{
				Precision:         4,
				TrimTrailingZeros: true,
			},
			want: "1.5 KiB",
		},
		{
			name:  "trim all fractional zeros",
			bytes: 2048,
			opts: Options{
				Precision:         4,
				TrimTrailingZeros: true,
			},
			want: "2 KiB",
		},
		{
			name:  "preserve trailing zeros",
			bytes: 1536,
			opts:  Options{Precision: 2},
			want:  "1.50 KiB",
		},
		{
			name:  "bytes ignore fractional precision",
			bytes: 999,
			opts: Options{
				Precision:         10,
				TrimTrailingZeros: false,
			},
			want: "999 B",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FormatWithOptions(tt.bytes, tt.opts); got != tt.want {
				t.Fatalf("FormatWithOptions(%d, %+v) = %q, want %q", tt.bytes, tt.opts, got, tt.want)
			}
		})
	}
}

func TestDecimalRoundHalfUp(t *testing.T) {
	t.Parallel()

	// 1280 bytes = exactly 1.25 KiB. At one decimal place the documented
	// round-half-up policy must produce 1.3, not banker's rounding to 1.2.
	got := FormatWithOptions(1280, Options{Precision: 1})
	if want := "1.3 KiB"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRoundingPromotion(t *testing.T) {
	t.Parallel()

	got := FormatWithOptions(
		(1<<20)-1,
		Options{Precision: 2, PromoteOnRounding: true},
	)
	if want := "1.00 MiB"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNoRoundingPromotion(t *testing.T) {
	t.Parallel()

	got := FormatWithOptions(
		(1<<20)-1,
		Options{Precision: 2, PromoteOnRounding: false},
	)
	if want := "1024.00 KiB"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNearBoundaryDoesNotPromoteTooEarly(t *testing.T) {
	t.Parallel()

	got := FormatWithOptions(
		1_048_000,
		Options{Precision: 2, PromoteOnRounding: true},
	)
	if want := "1023.44 KiB"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPrecisionNormalization(t *testing.T) {
	t.Parallel()

	if got, want := FormatWithOptions(1536, Options{Precision: -1}), "2 KiB"; got != want {
		t.Fatalf("negative precision: got %q, want %q", got, want)
	}

	if got, want := FormatWithOptions(1536, Options{Precision: 999}), "1.5000000000 KiB"; got != want {
		t.Fatalf("high precision: got %q, want %q", got, want)
	}
}

func TestOptionsValidate(t *testing.T) {
	t.Parallel()

	for _, precision := range []int{0, 1, MaxPrecision} {
		if err := (Options{Precision: precision}).Validate(); err != nil {
			t.Fatalf("precision %d unexpectedly invalid: %v", precision, err)
		}
	}

	for _, precision := range []int{-1, MaxPrecision + 1, 999} {
		err := (Options{Precision: precision}).Validate()
		if err == nil {
			t.Fatalf("precision %d unexpectedly valid", precision)
		}
		if !strings.Contains(err.Error(), "precision") {
			t.Fatalf("error %q does not describe precision", err)
		}
	}
}

func TestAppendPreservesPrefix(t *testing.T) {
	t.Parallel()

	got := string(Append([]byte("used="), 5_000_000))
	if want := "used=4.77 MiB"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestUnitSelectionAtMaxUint64(t *testing.T) {
	t.Parallel()

	if got := FormatWithOptions(math.MaxUint64, Options{Precision: MaxPrecision}); !strings.HasSuffix(got, " EiB") {
		t.Fatalf("expected EiB suffix, got %q", got)
	}
}
