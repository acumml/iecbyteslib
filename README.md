# iecbyteslib

`iecbyteslib` is a small, dependency-free Go library for formatting byte counts
with IEC binary units (`B`, `KiB`, `MiB`, `GiB`, `TiB`, `PiB`, `EiB`).

It is designed for production logging, monitoring, diagnostics, CLI output, and
other places where exact byte counts need a concise human-readable companion.

## Why this library

- **Exact across the full `uint64` input range**: formatting and rounding use
  integer arithmetic, not `float64`.
- **IEC-correct units**: powers of 1024 with the IEC names `KiB`, `MiB`, etc.
- **Defined rounding**: decimal **round-half-up**.
- **Boundary-aware promotion**: values that round to `1024.00 KiB` can become
  `1.00 MiB` instead.
- **Bounded precision**: 0-10 decimal places, preventing accidentally unbounded
  output from configuration.
- **Allocation-aware API**: `Append` and `AppendWithOptions` can reuse caller
  buffers.
- **`log/slog` helpers**: emit either a human-readable string or a structured
  group preserving both raw bytes and the IEC value.
- **No third-party dependencies**.
- **Concurrency-safe**: all formatting operations are stateless.

## Compatibility

The module targets Go 1.23 and later. It uses only the Go standard library.

## Install

```bash
go get github.com/acumml/iecbyteslib
```

## Basic formatting

```go
package main

import (
    "fmt"

    iecbytes "github.com/acumml/iecbyteslib"
)

func main() {
    fmt.Println(iecbytes.Format(500))
    // 500 B

    fmt.Println(iecbytes.Format(1024))
    // 1.00 KiB

    fmt.Println(iecbytes.Format(1_234_567_890))
    // 1.15 GiB
}
```

The default options are:

```go
iecbytes.Options{
    Precision:          2,
    TrimTrailingZeros:  false,
    PromoteOnRounding:  true,
}
```

## Structured `slog` usage

For operational telemetry, retaining the exact byte count as well as the
human-readable value is usually preferable:

```go
slog.Info(
    "redis memory sampled",
    iecbytes.GroupAttr("memory", usedMemory),
)
```

With `slog.JSONHandler`, the relevant part is shaped like:

```json
{
  "memory": {
    "bytes": 1234567890,
    "iec": "1.15 GiB"
  }
}
```

If a string-only field is desired:

```go
slog.Info(
    "redis memory sampled",
    iecbytes.Attr("used_memory", usedMemory),
)
```

## Custom formatting

```go
value := iecbytes.FormatWithOptions(
    1536,
    iecbytes.Options{
        Precision:          4,
        TrimTrailingZeros:  true,
        PromoteOnRounding:  true,
    },
)

// value == "1.5 KiB"
```

### Precision validation

`FormatWithOptions` clamps precision into the supported 0-10 range. This is
intentional so diagnostic/logging code remains safe even when passed a bad
runtime value.

For configuration loaded at application startup, validate first and fail fast:

```go
opts := iecbytes.Options{Precision: configuredPrecision}
if err := opts.Validate(); err != nil {
    return err
}
```

## Rounding and unit promotion

Rounding uses decimal **round-half-up**. For example, 1280 bytes is exactly
`1.25 KiB`, which formats as `1.3 KiB` at one decimal place.

With `PromoteOnRounding: true`, a value immediately below the next binary unit
can be promoted after rounding:

```text
1,048,575 bytes -> 1.00 MiB
```

With promotion disabled, the same value at precision 2 is:

```text
1,048,575 bytes -> 1024.00 KiB
```

Promotion affects presentation only. The raw byte count is never modified.

## Why integer-only rounding matters

Converting a large `uint64` to `float64` can lose integer precision because a
`float64` has 53 bits of integer precision. This library avoids that issue.
It performs long-division-style decimal digit generation from the integer
remainder, then rounds using a guard digit. The result is deterministic for the
entire `uint64` input range.

## Allocation-aware formatting

Callers in hot paths can reuse their own buffer:

```go
buf := make([]byte, 0, 32)
buf = iecbytes.Append(buf[:0], byteCount)
```

## API summary

```go
func Format(bytes uint64) string
func FormatWithOptions(bytes uint64, opts Options) string
func Append(dst []byte, bytes uint64) []byte
func AppendWithOptions(dst []byte, bytes uint64, opts Options) []byte

func Attr(key string, bytes uint64) slog.Attr
func AttrWithOptions(key string, bytes uint64, opts Options) slog.Attr
func GroupAttr(key string, bytes uint64) slog.Attr
func GroupAttrWithOptions(key string, bytes uint64, opts Options) slog.Attr

func DefaultOptions() Options
func (Options) Validate() error
```

## Development and verification

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
go test -run=^$ -bench=. -benchmem ./...
```

Fuzzing:

```bash
go test -fuzz='^FuzzFormat$' -fuzztime=30s
go test -fuzz='^FuzzFormatWithOptions$' -fuzztime=30s
```

## Design boundaries

This package deliberately does **not**:

- format SI decimal units such as `kB`, `MB`, or `GB`;
- parse human-readable sizes back into bytes;
- localise decimal separators;
- insert digit-grouping separators;
- silently accept negative byte counts (the primary API uses `uint64`).

Those features are separate concerns, and keeping them out preserves a small,
predictable API for byte-count observability.

## License

This project is licensed under the MIT License. See [`LICENSE`](LICENSE) for details.