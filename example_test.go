package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func Example() {
	fmt.Println(shelltoken.Parse("PATH=/bin ls -l"))
	// Output:
	// [PATH=/bin] [ls -l] <nil>
}
