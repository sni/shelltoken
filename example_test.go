package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func ExampleParseLinux() {
	env, argv, hasShell, err := shelltoken.ParseLinux("PATH=/bin ls -l")
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:   %#v\nargv:  %#v\nshell: %#v\n", env, argv, hasShell)
	// Output:
	// env:   []string{"PATH=/bin"}
	// argv:  []string{"ls", "-l"}
	// shell: false
}

func ExampleParseWindows() {
	env, argv, hasShell, err := shelltoken.ParseWindows(`'C:\Program Files\Vim\vim90\vim.exe' -n test.txt`)
	if err != nil {
		panic("parse error: " + err.Error())
	}

	fmt.Printf("env:   %#v\nargv:  %#v\nshell: %#v\n", env, argv, hasShell)
	// Output:
	// env:   []string{}
	// argv:  []string{"C:\\Program Files\\Vim\\vim90\\vim.exe", "-n", "test.txt"}
	// shell: false
}
