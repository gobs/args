package args

import (
	"fmt"
	"testing"
)

const (
	TEST_STRING  = `the   quick 	  "brown  'fox'"  jumps 'o v e r' \"the\"\ lazy dog`
	PARSE_STRING = "-l --number=42 -where=here one two three"
)

func TestScanner(test *testing.T) {
	scanner := NewScannerString(TEST_STRING)

	for {
		token, delim, err := scanner.NextToken()
		if err != nil {
			test.Log(err)
			break
		}

		test.Log(delim, token)
	}
}

func TestGetArgs(test *testing.T) {

	test.Log(GetArgs(TEST_STRING))
}

func TestParseArgs(test *testing.T) {

	test.Log(ParseArgs(PARSE_STRING))
}

func ExampleGetArgs() {
	s := `one two three "double quotes" 'single quotes' arg\ with\ spaces "\"quotes\" in 'quotes'" '"quotes" in \'quotes'"`

	for i, arg := range GetArgs(s) {
		fmt.Println(i, arg)
	}
	// Output:
	// 0 one
	// 1 two
	// 2 three
	// 3 double quotes
	// 4 single quotes
	// 5 arg with spaces
	// 6 "quotes" in 'quotes'
	// 7 "quotes" in 'quotes
}

func ExampleParseArgs() {
	arguments := "-l --number=42 -where=here one two three"

	parsed := ParseArgs(arguments)

	fmt.Println("options:", parsed.Options)
	fmt.Println("arguments:", parsed.Arguments)
	// Output:
	// options: map[l: number:42 where:here]
	// arguments: [one two three]
}
