/*
 This package provides methods to parse a shell-like command line string into a list of arguments.

 Words are split on white spaces, respecting quotes (single and double) and the escape character (backslash)
*/
package args

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

const (
	NO_QUOTE = unicode.ReplacementChar
)

type Scanner struct {
	in *bufio.Reader
}

// Creates a new Scanner with io.Reader as input source
func NewScanner(r io.Reader) *Scanner {
	sc := Scanner{in: bufio.NewReader(r)}
	return &sc
}

// Creates a new Scanner with a string as input source
func NewScannerString(s string) *Scanner {
	sc := Scanner{in: bufio.NewReader(strings.NewReader(s))}
	return &sc
}

// Get the next token from the Scanner, return io.EOF when done
func (this *Scanner) NextToken() (s string, delim int, err error) {
	buf := bytes.NewBufferString("")
	first := true
	escape := false
	quote := NO_QUOTE // invalid character - not a quote

	for {
		if c, _, e := this.in.ReadRune(); e == nil {
			if unicode.IsSpace(c) && !escape {
				if first { // skip leading spaces
					continue
				}

				if quote == NO_QUOTE { // not in quotes
					s = buf.String()
					delim = int(c)
					return // (token, delim, nil)
				}

				// otherwise we treat it as a regular character
			}

			if first {
				first = false

				if c == '"' || c == '\'' {
					quote = c // we are quoting
					first = false
					continue
				}
			}

			if !escape {
				if c == quote { // close quote
					s = buf.String()
					delim = int(c)
					return // (token, delim, nil)
				}

				if /* quote != NO_QUOTE && */ c == '\\' { // escape next
					escape = true
					continue
				}
			} else {
				escape = false
			}

			buf.WriteString(string(c))
		} else {
			if e == io.EOF {
				if buf.Len() > 0 {
					s = buf.String()
					return // (token, 0, nil)
				}
			}
			err = e
			return // ("", 0, io.EOF)
		}
	}

	return
}

// Return all tokens as an array of strings
func (this *Scanner) GetTokens() (tokens []string, err error) {

	tokens = make([]string, 0, 10) // an arbitrary initial capacity

	for {
		tok, _, err := this.NextToken()
		if err != nil {
			break
		}

		tokens = append(tokens, tok)
	}

	return
}

// Parse the input line into an array of arguments
func GetArgs(line string) (args []string) {
	scanner := NewScannerString(line)
	args, _ = scanner.GetTokens()
	return
}

type Args struct {
	Options   map[string]string
	Arguments []string
}

func ParseArgs(line string) (parsed Args) {
	parsed = Args{Options: map[string]string{}, Arguments: []string{}}
	args := GetArgs(line)
	if len(args) == 0 {
		return
	}

	for len(args) > 0 {
		arg := args[0]

		if !strings.HasPrefix(arg, "-") {
			break
		}

		args = args[1:]
		if arg == "--" { // stop parsing options
			break
		}

		arg = strings.TrimLeft(arg, "-")
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key := parts[0]
			value := parts[1]

			parsed.Options[key] = value
		} else {
			parsed.Options[arg] = ""
		}
	}

	parsed.Arguments = args
	return
}
