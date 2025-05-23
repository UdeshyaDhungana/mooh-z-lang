package parser

import (
	"testing"

	"github.com/udeshyadhungana/interprerer/app/ast"
	"github.com/udeshyadhungana/interprerer/app/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("encountered %d errors while parsing\n", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q\n", msg)
	}
	t.FailNow()
}

/* THOOS MUJI */

func TestThoosMujiStatements(t *testing.T) {
	input := `
	thoos_muji x = 34;
	thoos_muji y = x;
	thoos_muji foobar = 234543;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testThoosMujiStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testThoosMujiStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "thoos_muji" {
		t.Errorf("s.TokenLiteral not 'thoos_muji'. got=%q", s.TokenLiteral())
		return false
	}

	thoosStmt, ok := s.(*ast.ThoosMujiStatement)
	if !ok {
		t.Errorf("s not *ast.ThoosMujiStatement. got=%T", s)
		return false
	}

	if thoosStmt.Name.Value != name {
		t.Errorf("thoosStmt.Name.Value not '%s'. got=%s", name, thoosStmt.Name.Value)
		return false
	}

	if thoosStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, thoosStmt.Name)
		return false
	}

	return true
}

/* PATHA MUJI */
func TestPathaMujiStatements(t *testing.T) {
	input := `
	patha_muji 2;
	patha_muji a;
	patha_muji x + y;`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		ptmjStmt, ok := stmt.(*ast.PathaMujiStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnState")
			continue
		}
		if ptmjStmt.TokenLiteral() != "patha_muji" {
			t.Errorf("ptmjStmt.TokenLiteral not 'patha_muji', got %q", program.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "555;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 555 {
		t.Errorf("literal.Value is not %d. got=%d", 555, literal.Value)
	}
	if literal.TokenLiteral() != "555" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "555", literal.TokenLiteral())
	}
}
