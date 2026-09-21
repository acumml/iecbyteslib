package iecbytes

import "log/slog"

// Attr returns a string-valued slog attribute containing the DefaultOptions
// IEC representation of bytes.
//
// For observability data where the exact machine-readable byte count should be
// retained alongside the display value, prefer GroupAttr.
func Attr(key string, bytes uint64) slog.Attr {
	return slog.String(key, Format(bytes))
}

// AttrWithOptions is Attr with explicit formatting options.
func AttrWithOptions(key string, bytes uint64, opts Options) slog.Attr {
	return slog.String(key, FormatWithOptions(bytes, opts))
}

// GroupAttr returns a slog group containing both the exact raw byte count and
// its DefaultOptions IEC representation.
//
// With slog.JSONHandler, for example:
//
//	"memory":{"bytes":1234567890,"iec":"1.15 GiB"}
func GroupAttr(key string, bytes uint64) slog.Attr {
	return GroupAttrWithOptions(key, bytes, DefaultOptions())
}

// GroupAttrWithOptions is GroupAttr with explicit formatting options.
func GroupAttrWithOptions(key string, bytes uint64, opts Options) slog.Attr {
	return slog.Group(
		key,
		slog.Uint64("bytes", bytes),
		slog.String("iec", FormatWithOptions(bytes, opts)),
	)
}
