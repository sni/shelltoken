# shelltoken

[![Go Reference](https://pkg.go.dev/badge/github.com/sni/shelltoken.svg)](https://pkg.go.dev/github.com/sni/shelltoken)
[![License](https://img.shields.io/github/license/sni/shelltoken)](https://github.com/sni/shelltoken/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/sni/shelltoken)](https://goreportcard.com/report/github.com/sni/shelltoken)
[![CICD Pipeline](https://github.com/sni/shelltoken/actions/workflows/citest.yml/badge.svg)](https://github.com/sni/shelltoken/actions/workflows/citest.yml)

Go library to split a command line into env, command and arguments.

## Installation

    %> go get github.com/sni/shelltoken

## Example

    package main

    import (
        "github.com/sni/shelltoken"
    )

    env, argv, err := shelltoken.Parse("./command line with 'spaces' and back\slashes")
