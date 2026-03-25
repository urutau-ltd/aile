package basicauth_test

import (
	"fmt"

	basicauth "codeberg.org/urutau-ltd/aile/v2/x/basic_auth"
)

func ExampleValidator() {
	var validate basicauth.Validator = func(user, pass string) bool {
		return user == "admin" && pass == "secret"
	}

	fmt.Println(validate("admin", "secret"))
	// Output: true
}
