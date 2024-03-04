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

const WHITESPACE = " \t\n\r"

// ParseLinux splits a string the way the linux /bin/sh would do.
// It uses
// - separator: " \t\n\r".
// - keep backslashes: false.
// - keep quotes: false.
// - keep separator: false.
func ParseLinux(str string) (env, argv []string, hasShellCode bool, err error) {
	return Parse(str, WHITESPACE, false, false, false)
}

// ParseWindows splits a string the way windows would do.
// It uses
// - separator: " \t\n\r".
// - keep backslashes: true.
// - keep quotes: false.
// - keep separator: false.
func ParseWindows(str string) (env, argv []string, hasShellCode bool, err error) {
	return Parse(str, WHITESPACE, true, false, false)
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
// keepQuote controls wether quotes are kept or removed.
// hasShellCode is set to true if any shell special characters are found, ex.: sub shells like $(cmd)
// An unsuccessful parse will return an error.
func Parse(str, sep string, keepBackSlash, keepSep, keepQuote bool) (env, argv []string, hasShellCode bool, err error) {
	state := &parseState{
		hasToken:       false,
		escaped:        false,
		hasShellCode:   false,
		inSingleQuotes: false,
		inDoubleQuotes: false,
		token:          strings.Builder{},
	}

	str = strings.TrimSpace(str)

	for pos, char := range str {
		switch {
		case !state.escaped && char == '\\':
			state.escaped = true

			switch {
			case keepBackSlash, state.inSingleQuotes:
				// backslashes are kept in single quotes
				state.addToken(char)
			case state.inDoubleQuotes:
				// or in double quotes except...
				if len(str) > pos {
					switch str[pos+1] {
					// next character is a double quote again
					case '"':
					// or a backslash
					case '\\':
					default:
						state.addToken(char)
					}
				}
			}

		case !state.escaped && char == '"':
			state.hasToken = true

			if !state.inSingleQuotes {
				state.inDoubleQuotes = !state.inDoubleQuotes
				if keepQuote {
					state.addToken(char)
				}
			} else {
				state.addToken(char)
			}
		case !state.escaped && char == '\'':
			state.hasToken = true

			if !state.inDoubleQuotes {
				state.inSingleQuotes = !state.inSingleQuotes
				if keepQuote {
					state.addToken(char)
				}
			} else {
				state.addToken(char)
			}
		case !state.escaped && strings.ContainsRune(sep, char):
			switch {
			case state.inSingleQuotes, state.inDoubleQuotes:
				state.addToken(char)
			case keepSep:
				if state.hasToken {
					argv = append(argv, state.token.String())
					state.token.Reset()

					state.hasToken = false
				}

				argv = append(argv, string(char))
			case state.hasToken:
				argv = append(argv, state.token.String())
				state.token.Reset()

				state.hasToken = false
			}
		default:
			state.addToken(char)
		}
	}

	if !state.hasToken {
		// append empty token if no token found so far
		argv = append(argv, "")
	} else {
		// append last token
		argv = append(argv, state.token.String())
	}

	switch {
	case state.inSingleQuotes:
		return nil, nil, false, ErrUnbalancedQuotes
	case state.inDoubleQuotes:
		return nil, nil, false, ErrUnbalancedQuotes
	}

	env, argv = extractEnvFromArgv(argv)

	return env, argv, state.hasShellCode, nil
}

type parseState struct {
	hasToken       bool
	escaped        bool
	hasShellCode   bool
	inSingleQuotes bool
	inDoubleQuotes bool
	token          strings.Builder
}

func (p *parseState) addToken(char rune) {
	// reset escaped flag
	p.escaped = false

	switch {
	case p.inSingleQuotes:
	case p.inDoubleQuotes:
		switch char {
		case '$', '`':
			p.hasShellCode = true
		}
	default:
		switch char {
		case '$', '`', '!', '&', '*', '(', ')', '~', '[', ']', '\\', '|', '{', '}', ';', '<', '>', '?':
			p.hasShellCode = true
		}
	}

	p.hasToken = true

	p.token.WriteRune(char)
}

func extractEnvFromArgv(argv []string) (envs, args []string) {
	for i := range argv {
		if !strings.Contains(argv[i], "=") {
			return argv[0:i], argv[i:]
		}
	}

	return
}
