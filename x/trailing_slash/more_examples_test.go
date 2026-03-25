package trailingslash_test

import (
	"fmt"

	trailingslash "codeberg.org/urutau-ltd/aile/v2/x/trailing_slash"
)

func ExampleMode() {
	var mode trailingslash.Mode = trailingslash.RedirectTrim
	fmt.Println(mode == trailingslash.RedirectTrim)
	// Output: true
}

func ExampleRedirectTrim() {
	fmt.Println(trailingslash.RedirectTrim == trailingslash.Mode(1))
	// Output: true
}

func ExampleRedirectAppend() {
	fmt.Println(trailingslash.RedirectAppend == trailingslash.Mode(2))
	// Output: true
}
