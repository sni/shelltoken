package shelltoken_test

import (
	"fmt"

	"github.com/sni/shelltoken"
)

func ExampleSplitLinux() {
	env, argv, err := shelltoken.SplitLinux("PATH=/bin ls -l")
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{"PATH=/bin"}
	// argv: []string{"ls", "-l"}
}

func ExampleSplitWindows() {
	env, argv, err := shelltoken.SplitWindows(`'C:\Program Files\Vim\vim90\vim.exe' -n test.txt`)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("env:  %#v\nargv: %#v\n", env, argv)
	// Output:
	// env:  []string{}
	// argv: []string{"C:\\Program Files\\Vim\\vim90\\vim.exe", "-n", "test.txt"}
}

func ExampleSplitQuotes() {
	token, err := shelltoken.SplitQuotes(`ls -la | grep xyz;echo ok`, `|;`, shelltoken.SplitIgnoreShellCharacters|shelltoken.SplitKeepSeparator)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("token: %#v\n", token)
	// Output:
	// token: []string{"ls -la ", "|", " grep xyz", ";", "echo ok"}
}
