args
====

command line argument parser

Given a string (the "command line") it splits into words separated by white spaces, respecting quotes (single and double) and the escape character (backslash)

## Installation

    $ go get github.com/gobs/args

## Documentation
http://godoc.org/github.com/gobs/args

## Example

    import "github.com/gobs/args"
    
    s := `one two three "double quotes" 'single quotes' arg\ with\ spaces "\"quotes\" in 'quotes'" '"quotes" in \'quotes'"`

	for i, arg := range GetArgs(s) {
		fmt.Println(i, arg)
	}
