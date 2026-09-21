// Package iecbytes formats byte counts using IEC binary units.
//
// It supports B, KiB, MiB, GiB, TiB, PiB, and EiB, uses integer-only
// arithmetic for deterministic rounding across the full uint64 range, and
// includes helpers for structured log/slog output.
package iecbytes
