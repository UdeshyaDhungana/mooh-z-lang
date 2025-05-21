package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/udeshyadhungana/interprerer/app/lexer"
	"github.com/udeshyadhungana/interprerer/app/token"
)

const PROMPT = "(lekh_muji) >> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.NewLexer(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
