package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/udeshyadhungana/interprerer/app/token"
	"github.com/udeshyadhungana/interprerer/app/utils"
)

type Lexer struct {
	input        string
	position     int  // current position
	readPosition int  // current position + 1
	ch           rune // character under examination
	inComment    bool
	errors       []string
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readRune()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
	// symbols
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readRune()
			tok = token.NewTokenFromStr(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = token.NewToken(token.ASSIGN, l.ch)
		}
	// arithmetic
	case '+':
		tok = token.NewToken(token.PLUS, l.ch)
	case '-':
		tok = token.NewToken(token.MINUS, l.ch)
	case '*':
		tok = token.NewToken(token.ASTERISK, l.ch)
	case '/':
		tok = token.NewToken(token.SLASH, l.ch)
	case '%':
		tok = token.NewToken(token.MOD, l.ch)
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
	//array
	case '[':
		tok = token.NewToken(token.LBRACKET, l.ch)
	case ']':
		tok = token.NewToken(token.RBRACKET, l.ch)
	// logical
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readRune()
			tok = token.NewTokenFromStr(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = token.NewToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			lteq := fmt.Sprintf("%c%c", l.ch, l.peekChar())
			l.readRune()
			tok = token.NewTokenFromStr(token.LT_EQ, lteq)
		} else {
			tok = token.NewToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			gteq := fmt.Sprintf("%c%c", l.ch, l.peekChar())
			l.readRune()
			tok = token.NewTokenFromStr(token.GT_EQ, gteq)
		} else {
			tok = token.NewToken(token.GT, l.ch)
		}
	case ':':
		tok = token.NewToken(token.COLON, l.ch)
	// delimiters
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.ch)
	case ',':
		tok = token.NewToken(token.COMMA, l.ch)
	case '$':
		if l.inComment {
			l.inComment = false
			l.readRune()
			return l.NextToken()
		} else {
			l.inComment = true
			l.readRune()
			for l.ch != '$' && l.ch != 0 {
				l.readRune()
			}
			if l.ch == 0 {
				l.errors = append(l.errors, "unterminated comment")
				return token.Token{Type: token.EOF, Literal: ""}
			} else {
				return l.NextToken()
			}
		}
	case '(':
		tok = token.NewToken(token.LPAREN, l.ch)
	case ')':
		tok = token.NewToken(token.RPAREN, l.ch)
	case '{':
		tok = token.NewToken(token.LBRACE, l.ch)
	case '}':
		tok = token.NewToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if utils.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if utils.IsDigit(l.ch) {
			tok = l.readNumber()
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, l.ch)
		}
	}
	l.readRune()
	return tok
}

func (l *Lexer) Errors() []string {
	return l.errors
}

func (l *Lexer) ReportErrors() {
	for _, s := range l.errors {
		fmt.Println(s)
	}
}

func (l *Lexer) readRune() {
	var r rune
	var size int
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.ch = r
	l.position = l.readPosition
	l.readPosition += size
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for utils.IsLetter(l.ch) {
		l.readRune()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() token.Token {
	position := l.position
	for utils.IsDigit(l.ch) {
		l.readRune()
	}
	if l.ch == '.' {
		l.readRune()
		for utils.IsDigit(l.ch) {
			l.readRune()
		}
		return token.NewTokenFromStr(token.FLOAT, l.input[position:l.position])
	} else {
		return token.NewTokenFromStr(token.INT, l.input[position:l.position])
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readRune()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readRune()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// maybe we will need peekRune?
