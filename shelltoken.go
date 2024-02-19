package main

import (
	"errors"
	"strings"
)

var ErrUnbalancedQuotes = errors.New("unbalanced quotes")

// Parse parses command into list of envs and argv.
func Parse(str string) (env, argv []string, err error) {
	var token []rune

	separator := " \t\n\r"
	inQuotes := false
	inDbl := false
	escaped := false
	str = strings.TrimSpace(str)

	addToken := func(char rune) {
		escaped = false

		if token == nil {
			token = make([]rune, 0)
		}

		token = append(token, char)
	}

	for pos, char := range str {
		switch {
		case !escaped && char == '\\':
			escaped = true

			switch {
			case inQuotes:
				// backslashes are kept in single quotes
				addToken(char)
			case inDbl:
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

			if !inQuotes {
				inDbl = !inDbl
			} else {
				addToken(char)
			}
		case !escaped && char == '\'':
			if token == nil {
				token = make([]rune, 0)
			}

			if !inDbl {
				inQuotes = !inQuotes
			} else {
				addToken(char)
			}
		case !escaped && strings.ContainsRune(separator, char):
			switch {
			case inQuotes, inDbl:
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
	case inQuotes:
		return nil, nil, ErrUnbalancedQuotes
	case inDbl:
		return nil, nil, ErrUnbalancedQuotes
	}

	env, argv = extractEnvFromArgv(argv)

	return env, argv, nil
}

func extractEnvFromArgv(argv []string) (envs, args []string) {
	for i, s := range argv {
		if !strings.Contains(s, "=") {
			return argv[0:i], argv[i:]
		}
	}

	return
}
