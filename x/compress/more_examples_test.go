package compress_test

import (
	"fmt"

	"codeberg.org/urutau-ltd/aile/v2/x/compress"
)

func ExampleConfig() {
	cfg := compress.Config{
		MinSize: 256,
	}

	fmt.Println(cfg.MinSize)
	// Output: 256
}
