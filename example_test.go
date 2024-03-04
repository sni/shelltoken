package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func ExampleSplitLinux() {
	env, argv, err := shelltoken.SplitLinux("PATH=/bin ls -l")
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{"PATH=/bin"}
	// argv: []string{"ls", "-l"}
}

func ExampleSplitWindows() {
	env, argv, err := shelltoken.SplitWindows(`'C:\Program Files\Vim\vim90\vim.exe' -n test.txt`)
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{}
	// argv: []string{"C:\\Program Files\\Vim\\vim90\\vim.exe", "-n", "test.txt"}
}
