package iprestriction_test

import (
	"fmt"
	"net"

	iprestriction "codeberg.org/urutau-ltd/aile/v2/x/ip_restriction"
)

func ExampleConfig() {
	_, allow, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		panic(err)
	}

	cfg := iprestriction.Config{
		Allow:      []*net.IPNet{allow},
		TrustProxy: true,
	}

	fmt.Println(len(cfg.Allow))
	fmt.Println(cfg.TrustProxy)
	// Output:
	// 1
	// true
}
