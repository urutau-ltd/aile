package secureheaders_test

import (
	"fmt"

	secureheaders "codeberg.org/urutau-ltd/aile/v2/x/secure_headers"
)

func ExampleConfig() {
	cfg := secureheaders.Config{
		ContentTypeNosniff: true,
		FrameDeny:          true,
	}

	fmt.Println(cfg.ContentTypeNosniff)
	fmt.Println(cfg.FrameDeny)
	// Output:
	// true
	// true
}
