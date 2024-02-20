package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func Example() {
	env, argv, err := shelltoken.Parse("PATH=/bin ls -l")
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Println(env, argv)
	// Output:
	// [PATH=/bin] [ls -l]
}
