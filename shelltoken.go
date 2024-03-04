// Package shelltoken implements a command line parser.
//
// The shelltoken package splits a command line into token by whitespace
// characters while honoring single and double quotes.
// Backslashes and escaped quotes are supported as well.
package shelltoken

import (
	"errors"
	"strings"
)

var ErrUnbalancedQuotes = errors.New("unbalanced quotes")

// ParseLinux splits a string the way the linux /bin/sh would do.
// It uses
// - separator: " \t\n\r".
// - keep backslashes: false.
// - keep separator: false.
func ParseLinux(str string) (env, argv []string, hasShellCode bool, err error) {
	separator := " \t\n\r"

	return Parse(str, separator, false, false)
}

// ParseWindows splits a string the way windows would do
// It uses
// - separator: " \t\n\r".
// - keep backslashes: true.
// - keep separator: false.
func ParseWindows(str string) (env, argv []string, hasShellCode bool, err error) {
	separator := " \t\n\r"

	return Parse(str, separator, true, false)
}

// Parse parses command into list of envs and argv.
// A successful parse will return the env list with
// parsed environment variable definitions along with
// the argv list. The argv list will always contain at
// least one element (which can be empty).
// The argv[0] contains the command and all following elements
// are the arguments.
// keepBackslash controls wether backslashes are kept or parsed.
// - true: keep them (ex. useful for windows commands)
// - false: (default) parse backslashes like the sh/bash shell
// keepSep controls wether separators are kept or removed.
// hasShellCode is set to true if any shell special characters are found, ex.: sub shells like $(cmd)
// An unsuccessful parse will return an error.
func Parse(str, sep string, keepBackSlash, keepSep bool) (env, argv []string, hasShellCode bool, err error) {
	var token []rune

	inSingleQuotes := false
	inDoubleQuotes := false
	escaped := false
	str = strings.TrimSpace(str)

	addToken := func(char rune) {
		// reset escaped flag
		escaped = false

		if token == nil {
			token = make([]rune, 0)
		}

		switch {
		case inSingleQuotes:
		case inDoubleQuotes:
			switch char {
			case '$', '`':
				hasShellCode = true
			}
		default:
			switch char {
			case '$', '`', '!', '&', '*', '(', ')', '~', '[', ']', '\\', '|', '{', '}', ';', '<', '>', '?':
				hasShellCode = true
			}
		}

		token = append(token, char)
	}

	for pos, char := range str {
		switch {
		case !escaped && char == '\\':
			escaped = true

			switch {
			case keepBackSlash, inSingleQuotes:
				// backslashes are kept in single quotes
				addToken(char)
			case inDoubleQuotes:
				// or in double quotes except...
				if len(str) > pos {
					switch str[pos+1] {
					// next character is a double quote again
					case '"':
					// or a backslash
					case '\\':
					default:
						addToken(char)
					}
				}
			}

		case !escaped && char == '"':
			if token == nil {
				token = make([]rune, 0)
			}

			if !inSingleQuotes {
				inDoubleQuotes = !inDoubleQuotes
			} else {
				addToken(char)
			}
		case !escaped && char == '\'':
			if token == nil {
				token = make([]rune, 0)
			}

			if !inDoubleQuotes {
				inSingleQuotes = !inSingleQuotes
			} else {
				addToken(char)
			}
		case !escaped && strings.ContainsRune(sep, char):
			switch {
			case inSingleQuotes, inDoubleQuotes:
				addToken(char)
			case keepSep:
				addToken(char)
			case token != nil:
				argv = append(argv, string(token))
				token = nil
			}
		default:
			addToken(char)
		}
	}

	if token == nil {
		// append empty token if no token found so far
		argv = append(argv, "")
	} else {
		// append last token
		argv = append(argv, string(token))
	}

	switch {
	case inSingleQuotes:
		return nil, nil, false, ErrUnbalancedQuotes
	case inDoubleQuotes:
		return nil, nil, false, ErrUnbalancedQuotes
	}

	env, argv = extractEnvFromArgv(argv)

	return env, argv, hasShellCode, nil
}

func extractEnvFromArgv(argv []string) (envs, args []string) {
	for i := range argv {
		if !strings.Contains(argv[i], "=") {
			return argv[0:i], argv[i:]
		}
	}

	return
}
