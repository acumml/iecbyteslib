package iecbytes

import "strconv"

const base uint64 = 1024

var (
	unitNames = [...]string{
		"B",
		"KiB",
		"MiB",
		"GiB",
		"TiB",
		"PiB",
		"EiB",
	}

	unitDivisors = [...]uint64{
		1,
		1 << 10,
		1 << 20,
		1 << 30,
		1 << 40,
		1 << 50,
		1 << 60,
	}
)

// Format converts bytes to a human-readable IEC string using DefaultOptions.
//
// Examples:
//
//	Format(0)          // "0 B"
//	Format(1024)       // "1.00 KiB"
//	Format(5_000_000)  // "4.77 MiB"
func Format(bytes uint64) string {
	return string(Append(nil, bytes))
}

// FormatWithOptions converts bytes to a human-readable IEC string using opts.
//
// Precision outside the supported range is clamped defensively. Call
// Options.Validate during configuration loading if an invalid precision should
// instead fail fast.
func FormatWithOptions(bytes uint64, opts Options) string {
	return string(AppendWithOptions(nil, bytes, opts))
}

// Append appends the DefaultOptions representation of bytes to dst and returns
// the extended buffer. It is useful in allocation-sensitive paths.
func Append(dst []byte, bytes uint64) []byte {
	return AppendWithOptions(dst, bytes, DefaultOptions())
}

// AppendWithOptions appends the IEC representation of bytes to dst and returns
// the extended buffer.
func AppendWithOptions(dst []byte, bytes uint64, opts Options) []byte {
	opts = opts.normalized()

	unitIndex := selectUnit(bytes)

	// Bytes are always integral. Decimal precision is only meaningful once a
	// binary prefix is present.
	if unitIndex == 0 {
		dst = strconv.AppendUint(dst, bytes, 10)
		dst = append(dst, ' ', 'B')
		return dst
	}

	whole, fraction := roundParts(bytes, unitDivisors[unitIndex], opts.Precision)

	if opts.PromoteOnRounding && unitIndex < len(unitNames)-1 && whole >= base {
		unitIndex++
		whole, fraction = roundParts(bytes, unitDivisors[unitIndex], opts.Precision)
	}

	dst = strconv.AppendUint(dst, whole, 10)

	if opts.Precision > 0 {
		fractionEnd := opts.Precision
		if opts.TrimTrailingZeros {
			for fractionEnd > 0 && fraction[fractionEnd-1] == '0' {
				fractionEnd--
			}
		}

		if fractionEnd > 0 {
			dst = append(dst, '.')
			dst = append(dst, fraction[:fractionEnd]...)
		}
	}

	dst = append(dst, ' ')
	dst = append(dst, unitNames[unitIndex]...)
	return dst
}

func selectUnit(bytes uint64) int {
	// Iterating from the largest unit makes selection independent of
	// multiplication and therefore immune to overflow near math.MaxUint64.
	for i := len(unitDivisors) - 1; i > 0; i-- {
		if bytes >= unitDivisors[i] {
			return i
		}
	}
	return 0
}

// roundParts returns the rounded whole and fractional decimal digits for
// bytes/divisor. Rounding is decimal round-half-up.
//
// No floating-point arithmetic is used. The remainder is always smaller than
// 2^60 (the largest divisor), so remainder*10 cannot overflow uint64:
// (2^60-1)*10 < math.MaxUint64.
func roundParts(bytes, divisor uint64, precision int) (uint64, [MaxPrecision]byte) {
	whole := bytes / divisor
	remainder := bytes % divisor

	var fraction [MaxPrecision]byte
	for i := 0; i < precision; i++ {
		remainder *= 10
		digit := remainder / divisor
		remainder %= divisor
		fraction[i] = byte('0' + digit)
	}

	// Generate one guard digit. For non-negative values, guard >= 5 is exactly
	// decimal round-half-up; any later digits cannot change that decision.
	remainder *= 10
	guard := remainder / divisor
	if guard < 5 {
		return whole, fraction
	}

	if precision == 0 {
		return whole + 1, fraction
	}

	for i := precision - 1; i >= 0; i-- {
		if fraction[i] < '9' {
			fraction[i]++
			return whole, fraction
		}
		fraction[i] = '0'
	}

	return whole + 1, fraction
}
