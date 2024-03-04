package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func ExampleParseLinux() {
	env, argv, err := shelltoken.ParseLinux("PATH=/bin ls -l")
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{"PATH=/bin"}
	// argv: []string{"ls", "-l"}
}

func ExampleParseWindows() {
	env, argv, err := shelltoken.ParseWindows(`'C:\Program Files\Vim\vim90\vim.exe' -n test.txt`)
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{}
	// argv: []string{"C:\\Program Files\\Vim\\vim90\\vim.exe", "-n", "test.txt"}
}
