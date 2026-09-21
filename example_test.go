package iecbytes_test

import (
	"fmt"

	iecbytes "github.com/acumml/iecbyteslib"
)

func ExampleFormat() {
	fmt.Println(iecbytes.Format(500))
	fmt.Println(iecbytes.Format(1024))
	fmt.Println(iecbytes.Format(1_234_567_890))
	// Output:
	// 500 B
	// 1.00 KiB
	// 1.15 GiB
}

func ExampleFormatWithOptions() {
	fmt.Println(iecbytes.FormatWithOptions(
		1536,
		iecbytes.Options{
			Precision:         4,
			TrimTrailingZeros: true,
			PromoteOnRounding: true,
		},
	))
	// Output:
	// 1.5 KiB
}
