package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func Example() {
	env, argv, _, err := shelltoken.Parse("PATH=/bin ls -l", false)
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{"PATH=/bin"}
	// argv: []string{"ls", "-l"}
}
