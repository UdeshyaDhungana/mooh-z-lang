package utils

import (
	"io"
	"unicode"

	"github.com/udeshyadhungana/interprerer/app/object"
)

func IsLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func IsDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func IsTruthy(o object.Object) bool {
	switch o {
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	case object.NULL:
		return false
	default:
		return true
	}
}

func GetBoolRef(x bool) *object.Boolean {
	if x {
		return object.TRUE
	}
	return object.FALSE
}

/* Interpreter */
func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
