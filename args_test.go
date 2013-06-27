package args

import (
    "testing"
)

const (
    TEST_STRING = `the   quick 	  "brown  'fox'"  jumps 'o v e r' \"the\"\ lazy dog`
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
