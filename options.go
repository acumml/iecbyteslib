package iecbytes

import "fmt"

const (
	// MaxPrecision is the largest supported number of fractional decimal digits.
	// The limit keeps formatting bounded when Options values originate from
	// configuration or other untrusted input.
	MaxPrecision = 10
)

// Options controls IEC byte formatting.
//
// The zero value is valid and formats non-byte units with zero fractional
// digits, without trimming or rounding-based unit promotion. Use
// DefaultOptions for the recommended human-readable logging defaults.
type Options struct {
	// Precision is the number of digits after the decimal point for KiB and
	// larger units. Values passed to FormatWithOptions outside [0,
	// MaxPrecision] are clamped defensively. Validate can be used when invalid
	// configuration should be rejected instead.
	Precision int

	// TrimTrailingZeros removes insignificant zeroes from the fractional part.
	// For example, 1.50 KiB becomes 1.5 KiB.
	TrimTrailingZeros bool

	// PromoteOnRounding promotes to the next IEC unit when rounding would make
	// the displayed value 1024 in the current unit. For example, 1,048,575
	// bytes formats as 1.00 MiB at precision 2 instead of 1024.00 KiB.
	PromoteOnRounding bool
}

// DefaultOptions returns the recommended defaults for logs and user-facing
// operational output.
func DefaultOptions() Options {
	return Options{
		Precision:         2,
		TrimTrailingZeros: false,
		PromoteOnRounding: true,
	}
}

// Validate reports whether o can be used without normalization.
//
// FormatWithOptions deliberately clamps out-of-range precision so formatting
// remains safe in logging and error paths. Validate is intended for startup
// configuration checks where fail-fast behaviour is preferable.
func (o Options) Validate() error {
	if o.Precision < 0 || o.Precision > MaxPrecision {
		return fmt.Errorf(
			"iecbytes: precision %d outside supported range 0..%d",
			o.Precision,
			MaxPrecision,
		)
	}
	return nil
}

func (o Options) normalized() Options {
	if o.Precision < 0 {
		o.Precision = 0
	} else if o.Precision > MaxPrecision {
		o.Precision = MaxPrecision
	}
	return o
}
