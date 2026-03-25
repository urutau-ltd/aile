package bearerauth_test

import (
	"fmt"

	bearerauth "codeberg.org/urutau-ltd/aile/v2/x/bearer_auth"
)

func ExampleValidator() {
	var validate bearerauth.Validator = func(token string) bool {
		return token == "good-token"
	}

	fmt.Println(validate("good-token"))
	// Output: true
}
