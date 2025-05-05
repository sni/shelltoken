// Package shelltoken implements a command line tokenizer.
//
// The shelltoken package splits a command line into token by whitespace
// characters while honoring single and double quotes.
// Backslashes and escaped quotes are supported as well.
// Whitespace is defined as " \t\n\r"
// Shell-characters within single quotes are "$`"
// Shell-characters within double quotes are "$`!&*()~[]|\{};<>?"
package shelltoken

import (
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

const (
	Whitespace                  = " \t\n\r"
	DoubleQuoteShellCharacters  = "$`"
	OutsideQuoteShellCharacters = "$`!&*()~[]|{};<>?"
)

// SplitOption sets available parse options.
type SplitOption uint8

const (
	// SplitNoOptions is the zero value for options.
	SplitNoOptions SplitOption = 0

	// SplitKeepBackslashes: Do not remove backslashes.
	SplitKeepBackslashes SplitOption = 1 << iota

	// SplitIgnoreBackslashes does not escape characters by backslash.
	SplitIgnoreBackslashes

	// SplitKeepQuotes: Keep quotes in the final argv list.
	SplitKeepQuotes

	// SplitKeepSeparator: do not remove the split characters. They end up as a separate element in the argv list.
	SplitKeepSeparator

	// SplitStopOnShellCharacters cancels parsing upon the first shell character and returns ShellCharactersFoundError.
	SplitStopOnShellCharacters

	// SplitContinueOnShellCharacters returns ShellCharactersFoundError on shell characters but will parse to the end.
	SplitContinueOnShellCharacters

	// SplitIgnoreShellCharacters will ignore shell characters.
	SplitIgnoreShellCharacters

	// SplitKeepAndIgnoreAll just splits but keeps all characters.
	SplitKeepAll = SplitKeepQuotes | SplitKeepBackslashes | SplitKeepSeparator | SplitIgnoreShellCharacters
)

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
	argv, err = SplitQuotes(strings.TrimSpace(str), Whitespace, SplitStopOnShellCharacters)
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
	windowsOptions := SplitKeepBackslashes | SplitIgnoreBackslashes | SplitStopOnShellCharacters

	argv, err = SplitQuotes(strings.TrimSpace(str), Whitespace, windowsOptions)
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
// Options are a list of SplitOption(s) or a bitmask of SplitOption(s)
// An unsuccessful parse will return an error. The error will be either
// UnbalancedQuotesError or ShellCharactersFoundError.
func SplitQuotes(str, sep string, options ...SplitOption) (argv []string, err error) {
	pst := newParseState(options)
	argv = []string{}

	for pos, char := range str {
		if pst.stopShell && pst.firstShellPos != -1 {
			return nil, &ShellCharactersFoundError{pos: pst.firstShellPos}
		}

		switch {
		case pst.escaped:
			// reset escaped flag
			pst.escaped = false
			pst.addToken(char, pos)
		case char == '\\':
			if !pst.ignBackslashes {
				pst.escaped = true
			}

			switch {
			case pst.keepBackSlash, pst.inSingleQuotes:
				// backslashes are kept in single quotes
				pst.addToken(char, pos)
			case pst.inDoubleQuotes:
				// or in double quotes except...
				if len(str) > pos {
					switch str[pos+1] {
					// next character is a double quote again
					case '"':
					// or a backslash
					case '\\':
					default:
						pst.addToken(char, pos)
					}
				}
			}

		case char == '"':
			pst.hasToken = true

			if !pst.inSingleQuotes {
				pst.inDoubleQuotes = !pst.inDoubleQuotes
				if pst.keepQuote {
					pst.addToken(char, pos)
				}
			} else {
				pst.addToken(char, pos)
			}
		case char == '\'':
			pst.hasToken = true

			if !pst.inDoubleQuotes {
				pst.inSingleQuotes = !pst.inSingleQuotes
				if pst.keepQuote {
					pst.addToken(char, pos)
				}
			} else {
				pst.addToken(char, pos)
			}
		case strings.ContainsRune(sep, char):
			switch {
			case pst.inSingleQuotes, pst.inDoubleQuotes:
				pst.addToken(char, pos)
			case pst.keepSep:
				if pst.hasToken {
					argv = append(argv, pst.token.String())
					pst.token.Reset()

					pst.hasToken = false
				}

				argv = append(argv, string(char))
			case pst.hasToken:
				argv = append(argv, pst.token.String())
				pst.token.Reset()

				pst.hasToken = false
			}
		default:
			pst.addToken(char, pos)
		}
	}

	// in case the last character was a shell char
	if pst.stopShell && pst.firstShellPos != -1 {
		return nil, &ShellCharactersFoundError{pos: pst.firstShellPos}
	}

	// append last token
	if pst.hasToken {
		argv = append(argv, pst.token.String())
	}

	switch {
	case pst.inSingleQuotes, pst.inDoubleQuotes:
		return nil, &UnbalancedQuotesError{}
	case pst.contShell && pst.firstShellPos != -1:
		return argv, &ShellCharactersFoundError{pos: pst.firstShellPos}
	default:
		return argv, nil
	}
}

type parseState struct {
	token strings.Builder

	// current state flags
	hasToken       bool
	escaped        bool
	inSingleQuotes bool
	inDoubleQuotes bool
	firstShellPos  int // position of first shell character found
	// parse flags
	keepBackSlash  bool
	keepQuote      bool
	keepSep        bool
	stopShell      bool
	contShell      bool
	ignShell       bool
	ignBackslashes bool
}

func newParseState(options []SplitOption) *parseState {
	pst := &parseState{
		hasToken:       false,
		escaped:        false,
		inSingleQuotes: false,
		inDoubleQuotes: false,
		token:          strings.Builder{},
		firstShellPos:  -1,
		keepBackSlash:  false,
		keepQuote:      false,
		keepSep:        false,
		stopShell:      false,
		contShell:      false,
		ignShell:       false,
		ignBackslashes: false,
	}

	option := SplitNoOptions
	for _, o := range options {
		option |= o
		if o == SplitNoOptions {
			option = SplitNoOptions
		}
	}

	pst.keepBackSlash = option&SplitKeepBackslashes > 0
	pst.keepQuote = option&SplitKeepQuotes > 0
	pst.keepSep = option&SplitKeepSeparator > 0
	pst.stopShell = option&SplitStopOnShellCharacters > 0
	pst.contShell = option&SplitContinueOnShellCharacters > 0
	pst.ignBackslashes = option&SplitIgnoreBackslashes > 0
	pst.ignShell = (!pst.stopShell && !pst.contShell) || option&SplitIgnoreShellCharacters > 0

	return pst
}

func (p *parseState) addToken(char rune, pos int) {
	p.hasToken = true

	// exit early if we do not search for shell characters (anymore)
	switch {
	case p.ignShell, p.inSingleQuotes, p.firstShellPos != -1:
		p.token.WriteRune(char)

		return
	}

	switch {
	case p.inDoubleQuotes:
		if strings.ContainsRune(DoubleQuoteShellCharacters, char) {
			p.firstShellPos = pos
		}
	case strings.ContainsRune(OutsideQuoteShellCharacters, char):
		p.firstShellPos = pos
	case char == '\\':
		if !p.keepBackSlash && !p.ignBackslashes {
			p.firstShellPos = pos
		}
	}

	p.token.WriteRune(char)
}
