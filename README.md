# shelltoken

[![Go Reference](https://pkg.go.dev/badge/github.com/sni/shelltoken.svg)](https://pkg.go.dev/github.com/sni/shelltoken)
[![License](https://img.shields.io/github/license/sni/shelltoken)](https://github.com/sni/shelltoken/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/sni/shelltoken)](https://goreportcard.com/report/github.com/sni/shelltoken)
[![CICD Pipeline](https://github.com/sni/shelltoken/actions/workflows/citest.yml/badge.svg)](https://github.com/sni/shelltoken/actions/workflows/citest.yml)

Go library to split a command line into env, command and arguments.

## Installation

    %> go get github.com/sni/shelltoken

## Documentation

The documenation can be found on [pkg.go.dev](https://pkg.go.dev/github.com/sni/shelltoken).

## Example

```golang
package main

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
```
