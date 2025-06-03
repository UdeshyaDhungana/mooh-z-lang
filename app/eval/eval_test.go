package eval

import (
	"fmt"
	"testing"

	"github.com/udeshyadhungana/interprerer/app/lexer"
	"github.com/udeshyadhungana/interprerer/app/object"
	"github.com/udeshyadhungana/interprerer/app/parser"
)

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	if len(p.Errors()) != 0 {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return nil
	}
	return Eval(program, env)
}

func TestEvalIntegerStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
		{"-5;", -5},
		{"-10;", -10},
		{"5 + 5 + 5 + 5 - 10;", 10},
		{"2 * 2 * 2 * 2 * 2;", 32},
		{"-50 + 100 + -50;", 0},
		{"5 * 2 + 10;", 20},
		{"5 + 2 * 10;", 25},
		{"20 + 2 * -10;", 0},
		{"50 / 2 * 2 + 10;", 60},
		{"2 * (5 + 10);", 30},
		{"3 * 3 * 3 + 10;", 37},
		{"3 * (3 * 3) + 10;", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBoolStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"sacho_muji;", true},
		{"jhut_muji;", false},
		{"1 < 2;", true},
		{"1 > 2;", false},
		{"1 < 1;", false},
		{"1 > 1;", false},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"1 == 2;", false},
		{"1 != 2;", true},
		{"sacho_muji == sacho_muji;", true},
		{"jhut_muji == jhut_muji;", true},
		{"sacho_muji == jhut_muji;", false},
		{"sacho_muji != jhut_muji;", true},
		{"jhut_muji != sacho_muji;", true},
		{"(1 < 2) == sacho_muji;", true},
		{"(1 < 2) == jhut_muji;", false},
		{"(1 > 2) == sacho_muji;", false},
		{"(1 > 2) == jhut_muji;", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!sacho_muji;", false},
		{"!jhut_muji;", true},
		{"!5;", false},
		{"!!sacho_muji;", true},
		{"!!jhut_muji;", false},
		{"!!5;", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func TestYediMujiExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"yedi_muji(sacho_muji) { 10; }", 10},
		{"yedi_muji (jhut_muji) { 10; }", nil},
		{"yedi_muji (1) { 10; }", 10},
		{"yedi_muji (1 < 2) { 10; }", 10},
		{"yedi_muji (1 > 2) { 10; }", nil},
		{"yedi_muji (1 > 2) { 10; } nabhae_chikne { 20; }", 20},
		{"yedi_muji (1 < 2) { 10; } nabhae_chikne { 20; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return false
}

func TestPathaMujiStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"patha_muji 10;", 10},
		{"patha_muji 233;", 233},
		{"yedi_muji(10 > 1) { patha_muji 1; } nabhae_chikne { patha_muji 0; }", 1},
		{"yedi_muji(10 < 1) { patha_muji 1; } nabhae_chikne { patha_muji 0; }", 0},
		{`
		yedi_muji (10 > 1) {
			yedi_muji (10 > 1) {
				patha_muji 1;
			}
			patha_muji 2;
		}
		`, 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + sacho_muji;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + jhut_muji; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-sacho_muji;",
			"unknown operator: -BOOLEAN",
		},
		{
			"sacho_muji + jhut_muji;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; sacho_muji + jhut_muji; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"yedi_muji (10 > 1) { sacho_muji + jhut_muji; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`yedi_muji (10 > 1) {
				yedi_muji (20 > 1) {
					patha_muji sacho_muji + jhut_muji;
				}
				patha_muji 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}

	}
}

/* TestThoosMuji */
func TestThoosMujiStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"thoos_muji a = 5; a;", 5},
		{"thoos_muji a = 5 * 5; a;", 25},
		{"thoos_muji a = 5; thoos_muji b = a; b;", 5},
		{"thoos_muji a = 5; thoos_muji b = a; thoos_muji c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"thoos_muji identity = kaam_gar_muji(x) { x; }; identity(5);", 5},
		{"thoos_muji identity = kaam_gar_muji(x) { patha_muji x; }; identity(5);", 5},
		{"thoos_muji double = kaam_gar_muji(x) { x * 2; }; double(5);", 10},
		{"thoos_muji add = kaam_gar_muji(x, y) { x + y; }; add(5, 5);", 10},
		{"thoos_muji add = kaam_gar_muji(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"kaam_gar_muji(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestRecursion(t *testing.T) {
	tests := []struct {
		program  string
		expected int64
	}{
		{`
		thoos_muji recursion = kaam_gar_muji(x) {
			yedi_muji (x == 0) {
				patha_muji 1;
			} nabhae_chikne {
			 	patha_muji x * recursion(x - 1);
			}
		};
		recursion(4);
		`,
			24},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.program)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestArrayIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`thoos_muji x = [1,2,3];
			x[0]`,
			1,
		},
		{
			`thoos_muji y = [sacho_muji, jhut_muji, jhut_muji];
			y[2]`,
			false,
		},
		{
			`thoos_muji y = ["Udeshya", "Dhungana"];
			y[0]`,
			"Udeshya",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBoolObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestJabasammaMujiExpressionEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`
			thoos_muji i = 0;
			thoos_muji sum = 0;
			jaba_samma_muji(i < 10) {
				thoos_muji sum = sum + i + 1;
			}
			`,
			55,
		},
		// {
		// 	``,
		// 	100,
		// },
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			t.Fatalf("failed to test JabasammaMujiExpression")
		}
	}
}

// re-assignment

func TestCustom(t *testing.T) {
	program := `
	thoos_muji makeGreeter = kaam_gar_muji(greeting) {
		patha_muji kaam_gar_muji(name) {
			patha_muji greeting + " " + name + "!";
		};
	};
	thoos_muji hello = makeGreeter("Hello");
	hello("Udeshya");
	`

	evaluated := testEval(program)
	if evaluated.Inspect() != "Hello Udeshya!" {
		t.Fatalf("test failed")
	}
}
