// Package shelltoken implements a command line tokenizer.
//
// The shelltoken package splits a command line into token by whitespace
// characters while honoring single and double quotes.
// Backslashes and escaped quotes are supported as well.
package shelltoken

import (
	"errors"
	"fmt"
	"strings"
)

type ShellCharactersFoundError struct {
	pos int
}

func (e *ShellCharactersFoundError) Error() string {
	return fmt.Sprintf("shell character at position %d", e.pos)
}

type UnbalancedQuotesError struct{}

func (e *UnbalancedQuotesError) Error() string {
	return "unbalanced quotes"
}

var ErrUnbalancedQuotes = errors.New("unbalanced quotes")

const WHITESPACE = " \t\n\r"

// SplitLinux will tokenize a string the way the linux /bin/sh would do.
// A successful parse will return the env list with
// parsed environment variable definitions along with
// the argv list. The argv list will always contain at
// least one element (which can be empty).
// The argv[0] contains the command and all following elements
// are the arguments.
// It uses
// - separator: " \t\n\r".
// - keep backslashes: false.
// - keep quotes: false.
// - keep separator: false.
// returns error if shell characters were found.
func SplitLinux(str string) (env, argv []string, err error) {
	argv, err = SplitQuotes(strings.TrimSpace(str), WHITESPACE, false, false, false, false)
	if err != nil {
		return nil, nil, err
	}

	if len(argv) == 0 {
		argv = append(argv, "")
	}

	env, argv = ExtractEnvFromArgv(argv)

	return env, argv, nil
}

// SplitWindows will tokenize a string the way windows would do.
// A successful parse will return the env list with
// parsed environment variable definitions along with
// the argv list. The argv list will always contain at
// least one element (which can be empty).
// The argv[0] contains the command and all following elements
// are the arguments.
// It uses
// - separator: " \t\n\r".
// - keep backslashes: true.
// - keep quotes: false.
// - keep separator: false.
// returns error if shell characters were found.
func SplitWindows(str string) (env, argv []string, err error) {
	argv, err = SplitQuotes(strings.TrimSpace(str), WHITESPACE, true, false, false, false)
	if err != nil {
		return nil, nil, err
	}

	if len(argv) == 0 {
		argv = append(argv, "")
	}

	env, argv = ExtractEnvFromArgv(argv)

	return env, argv, nil
}

// ExtractEnvFromArgv splits list of arguments into env and args.
func ExtractEnvFromArgv(argv []string) (envs, args []string) {
	for i := range argv {
		if !strings.Contains(argv[i], "=") {
			return argv[0:i], argv[i:]
		}
	}

	return
}

// SplitQuotes will tokenize text into chunks honoring quotes.
// keepBackslash controls wether backslashes are kept or removed.
// - true: keep them (ex. useful for windows commands)
// - false: (default) parse backslashes like the sh/bash shell
// keepSep controls wether separators are kept or removed.
// keepQuote controls wether quotes are kept or removed
// ignoreShellChars controls wether shell characters lead to a ShellCharactersFoundError
// An unsuccessful parse will return an error. The error will be either
// ErrUnbalancedQuotes or ShellCharactersFoundError.
func SplitQuotes(str, sep string, keepBackSlash, keepSep, keepQuote, ignoreShellChars bool) (argv []string, err error) {
	argv = []string{}
	state := &parseState{
		hasToken:       false,
		escaped:        false,
		hasShellCode:   false,
		inSingleQuotes: false,
		inDoubleQuotes: false,
		token:          strings.Builder{},
	}

	for pos, char := range str {
		if state.hasShellCode && !ignoreShellChars {
			return nil, &ShellCharactersFoundError{pos: pos}
		}

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

	// append last token
	if state.hasToken {
		argv = append(argv, state.token.String())
	}

	switch {
	case state.inSingleQuotes, state.inDoubleQuotes:
		return nil, ErrUnbalancedQuotes
	default:
		return argv, nil
	}
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
		case '$', '`', '!', '&', '*', '(', ')', '~', '[', ']', '|', '{', '}', ';', '<', '>', '?':
			p.hasShellCode = true
		}
	}

	p.hasToken = true

	p.token.WriteRune(char)
}
