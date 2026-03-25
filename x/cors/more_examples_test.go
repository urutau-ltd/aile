package cors_test

import (
	"fmt"

	"codeberg.org/urutau-ltd/aile/v2/x/cors"
)

func ExampleConfig() {
	cfg := cors.Config{
		AllowOrigins: []string{"https://app.example.com"},
		MaxAge:       600,
	}

	fmt.Println(cfg.AllowOrigins[0])
	fmt.Println(cfg.MaxAge)
	// Output:
	// https://app.example.com
	// 600
}
