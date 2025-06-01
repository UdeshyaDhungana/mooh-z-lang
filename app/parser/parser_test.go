package parser

import (
	"fmt"
	"testing"

	"github.com/udeshyadhungana/interprerer/app/ast"
	"github.com/udeshyadhungana/interprerer/app/lexer"
	"github.com/udeshyadhungana/interprerer/app/token"
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
	thoos_muji x = 4;
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

func TestBooleanLiteralExpression(t *testing.T) {
	input := "jhut_muji;"

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

	literal, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != false {
		t.Errorf("literal.Value is not %t. got=%t", false, literal.Value)
	}
	if literal.TokenLiteral() != "jhut_muji" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "555", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", int64(5)},
		{"-15;", "-", int64(15)},
		{"!jhut_muji;", "!", false},
		{"!sacho_muji;", "!", true},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. Got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteral(t, exp.Right, tt.value) {
			return
		}
	}
}

func testLiteral(t *testing.T, il ast.Expression, value any) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if ok {
		v, ok := value.(int64)
		if !ok {
			t.Errorf("value is not int64. got=%T", v)
			return false
		}
		if integ.Value != v {
			t.Errorf("integ.Value not (%d). got=(%d)", v, integ.Value)
			return false
		}

		if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
			t.Errorf("integ.TokenLiteral not %s. got=%s", value, integ.TokenLiteral())
			return false
		}
		return true
	}
	b, ok := il.(*ast.Boolean)
	if ok {
		v, ok := value.(bool)
		if !ok {
			t.Errorf("value is not bool. got=%T", b)
			return false
		}
		if b.Value != v {
			t.Errorf("boolean.Value not %t, got=%t", v, b.Value)
			return false
		}

		if b.Token.Type != token.JHUT_MUJI && b.Token.Type != token.SACHO_MUJI {
			t.Errorf("boolean.TokenLiteral not %s or %s. got=%s", token.SACHO_MUJI, token.JHUT_MUJI, b.TokenLiteral())
			return false
		}
		return true
	}

	id, ok := il.(*ast.Identifier)
	if ok {
		v, ok := value.(string)
		if !ok {
			t.Errorf("value is not string. got=%T", b)
			return false
		}
		if id.Value != v {
			t.Errorf("id.Value not %s, got=%s", v, id.Value)
		}

		if id.Token.Type != token.IDFIER {
			t.Errorf("identifier.TokenLiteral not %s, got=%s\n", token.IDFIER, id.Token.Type)
			return false
		}
		return true
	}

	t.Errorf("il is not *ast.IntegerLiteral or *ast.Boolean. got=%T", il)
	return false
}

func TestParsingInixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d elements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + (2 + 3) + 4;", "((1 + (2 + 3)) + 4);"},
		{"-a * b;",
			"((-a) * b);"},
		{"!-a;",
			"(!(-a));"},
		{
			"a + b + c;",
			"((a + b) + c);",
		},
		{
			"a + b - c;",
			"((a + b) - c);",
		},
		{
			"a * b * c;",
			"((a * b) * c);",
		},
		{
			"a * b / c;",
			"((a * b) / c);",
		},
		{
			"a + b / c;",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4; -5 * 5;",
			"(3 + 4);((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"sacho_muji;",
			"sacho_muji;",
		},
		{
			"jhut_muji;",
			"jhut_muji;",
		},
		{
			"3 < 5 == jhut_muji;",
			"((3 < 5) == jhut_muji);",
		},
		{
			"3 > 5 == sacho_muji;",
			"((3 > 5) == sacho_muji);",
		},
		{
			"a + add(b * c) + d;",
			"((a + add((b * c))) + d);",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		},
		{
			"add(a + b + c * d / f + g);",
			"add((((a + b) + ((c * d) / f)) + g));",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=[%q], got=[%q]", tt.expected, actual)
		}
	}
}

func TestYediMujiStatementParsing(t *testing.T) {
	tests := []struct {
		program  string
		expected string
	}{
		{
			`yedi_muji (2 + 3 == 5) {
		1;
		};`,
			"",
		},
		{
			`yedi_muji (sacho_muji) {
			23;
			} nabhae_chikne {
			43;
			};`,
			"",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.program)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		_, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
	}
}

func TestKaamGarMujiStatement(t *testing.T) {
	tests := []struct {
		literal        string
		expectedParams []string
	}{
		{
			"kaam_gar_muji(x, y) { patha_muji x + y; };",
			[]string{"x", "y"},
		},
		{
			"kaam_gar_muji() { patha_muji 3; };",
			[]string{},
		},
		{
			"kaam_gar_muji(a, b, c) { patha_muji a + b + c; };",
			[]string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.literal)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok || stmt == nil {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement")
		}
		f, ok := stmt.Expression.(*ast.KaamGarMujiExpression)
		if !ok || f == nil {
			t.Fatalf("stmt.Expression is not *ast.KaamGarMujiExpression")
		}

		if len(f.Arguments) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(f.Arguments))
		}

		for i, ident := range tt.expectedParams {
			testLiteral(t, f.Arguments[i], ident)
		}
	}
}

func TestCustom(t *testing.T) {
	tests := []struct {
		statement string
		expected  string
	}{
		{
			"thoos_muji x = kaam_gar_muji(x) { patha_muji x; };",
			"thoos_muji x = kaam_gar_muji(x) { patha_muji x; };",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.statement)
		p := NewParser(l)

		program := p.ParseProgram()
		fmt.Println(program.String())
	}
}
