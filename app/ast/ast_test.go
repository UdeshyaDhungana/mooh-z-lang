package ast

import (
	"testing"

	"github.com/udeshyadhungana/interprerer/app/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ThoosMujiStatement{
				Token: token.Token{Type: token.THOOS_MUJI, Literal: "thoos_muji"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "thoos_muji myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
