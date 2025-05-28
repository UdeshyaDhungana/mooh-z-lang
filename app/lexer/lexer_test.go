package lexer

import (
	"testing"

	"github.com/udeshyadhungana/interprerer/app/token"
)

func TestNextToken(t *testing.T) {
	input := `thoos_muji 界 = 5;
	thoos_muji ten = 10;
	
	thoos_muji add = kaam_gar_muji(x, y) {
		x + y;
	};
	
	thoos_muji result = add(five, ten);
	!-/*5;
	5 < 10 > 66;
	
	yedi_muji (5 < 10) {
		patha_muji sacho_muji;
	} nabhae_chikne {
	 	patha_muji jhut_muji;
	}
	
	10 == 10;
	10 != 9;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.THOOS_MUJI, "thoos_muji"},
		{token.IDFIER, "界"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.THOOS_MUJI, "thoos_muji"},
		{token.IDFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.THOOS_MUJI, "thoos_muji"},
		{token.IDFIER, "add"},
		{token.ASSIGN, "="},
		{token.KAAM_GAR_MUJI, "kaam_gar_muji"},
		{token.LPAREN, "("},
		{token.IDFIER, "x"},
		{token.COMMA, ","},
		{token.IDFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDFIER, "x"},
		{token.PLUS, "+"},
		{token.IDFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.THOOS_MUJI, "thoos_muji"},
		{token.IDFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDFIER, "add"},
		{token.LPAREN, "("},
		{token.IDFIER, "five"},
		{token.COMMA, ","},
		{token.IDFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "66"},
		{token.SEMICOLON, ";"},
		{token.YEDI_MUJI, "yedi_muji"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.PATHA_MUJI, "patha_muji"},
		{token.SACHO_MUJI, "sacho_muji"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.NABHAE_CHIKNE, "nabhae_chikne"},
		{token.LBRACE, "{"},
		{token.PATHA_MUJI, "patha_muji"},
		{token.JHUT_MUJI, "jhut_muji"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] failed - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] failed - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
