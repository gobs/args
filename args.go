/*
 This package provides methods to parse a shell-like command line string into a list of arguments.

 Words are split on white spaces, respecting quotes (single and double) and the escape character (backslash)
*/
package args

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	ESCAPE_CHAR  = '\\'
	QUOTE_CHARS  = "`'\""
	SYMBOL_CHARS = `|><#{([`
	NO_QUOTE     = unicode.ReplacementChar
)

var (
	BRACKETS = map[rune]rune{
		'{': '}',
		'[': ']',
		'(': ')',
	}
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
	quote := NO_QUOTE    // invalid character - not a quote
	brackets := []rune{} // stack of open brackets

	for {
		if c, _, e := this.in.ReadRune(); e == nil {
			//
			// check escape character
			//
			if c == ESCAPE_CHAR && !escape {
				escape = true
				first = false
				continue
			}

			//
			// if escaping, just add the character
			//
			if escape {
				escape = false
				buf.WriteString(string(c))
				continue
			}

			//
			// checks for beginning of token
			//
			if first {
				if unicode.IsSpace(c) {
					//
					// skip leading spaces
					//
					continue
				}

				first = false

				if strings.ContainsRune(QUOTE_CHARS, c) {
					//
					// start quoted token
					//
					quote = c
					continue
				}

				if b, ok := BRACKETS[c]; ok {
					//
					// start a bracketed session
					//
					delim = int(c)
					brackets = append(brackets, b)
					buf.WriteString(string(c))
					continue
				}

				if strings.ContainsRune(SYMBOL_CHARS, c) {
					//
					// if it's a symbol, return  all the remaining characters
					//
					buf.WriteString(string(c))
					_, err = io.Copy(buf, this.in)
					s = buf.String()
					return // (token, delim, err)
				}
			}

			if len(brackets) == 0 {
				//
				// terminate on spaces
				//
				if unicode.IsSpace(c) && quote == NO_QUOTE {
					s = buf.String()
					delim = int(c)
					return // (token, delim, nil)
				}

				//
				// close quote and terminate
				//
				if c == quote {
					quote = NO_QUOTE
					s = buf.String()
					delim = int(c)
					return // (token, delim, nil)
				}

				//
				// append to buffer
				//
				buf.WriteString(string(c))
			} else {
				//
				// append to buffer
				//
				buf.WriteString(string(c))

				last := len(brackets) - 1

				if quote == NO_QUOTE {
					if c == brackets[last] {
						brackets = brackets[:last] // pop

						if len(brackets) == 0 {
							s = buf.String()
							return // (token, delim, nil)
						}
					} else if strings.ContainsRune(QUOTE_CHARS, c) {
						//
						// start quoted token
						//
						quote = c
					} else if b, ok := BRACKETS[c]; ok {
						brackets = append(brackets, b)
					}
				} else if c == quote {
					quote = NO_QUOTE
				}
			}
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

// Create a new FlagSet to be used with ParseFlags
func NewFlags(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)

	flags.Usage = func() {
		fmt.Printf("Usage of %s:\n", name)
		flags.PrintDefaults()
	}

	return flags
}

// Parse the input line through the (initialized) FlagSet
func ParseFlags(flags *flag.FlagSet, line string) error {
	return flags.Parse(GetArgs(line))
}
