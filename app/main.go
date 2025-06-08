package main

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/udeshyadhungana/interprerer/app/eval"
	"github.com/udeshyadhungana/interprerer/app/lexer"
	"github.com/udeshyadhungana/interprerer/app/object"
	"github.com/udeshyadhungana/interprerer/app/parser"
	"github.com/udeshyadhungana/interprerer/app/repl"
)

func main() {
	if len(os.Args) == 1 {
		startRepl()
	} else {
		filePath := os.Args[1]
		interpret(filePath)
	}
}

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("नमस्कार %s मुजी!\n", user.Username)
	fmt.Println("यो \"मुजी\" भाषा हो। तल लेख् मुजी 👇")
	repl.Start(os.Stdin, os.Stdout)
}

func interpret(filepath string) {
	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	l := lexer.NewLexer(string(fileContents))
	p := parser.NewParser(l)
	program := p.ParseProgram()
	hasErrs := p.CheckAndReportErrors()
	if hasErrs {
		return
	}

	env := object.NewEnvironment()
	evaluated := eval.Eval(program, env)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}
