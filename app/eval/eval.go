package eval

import (
	"github.com/udeshyadhungana/interprerer/app/ast"
	"github.com/udeshyadhungana/interprerer/app/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return object.TRUE
		}
		return object.FALSE
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(left, node.Operator, right)
	case *ast.YediMujiStatement:
		return evalYediMujiStatement(node.Condition, node.Consequent, node.Alternative)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NULL
	}
	value := right.(*object.Integer)
	return &object.Integer{Value: -value.Value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return object.NULL
	}
	if left.Type() == object.INTEGER_OBJ {
		return evalForInteger(left, operator, right)
	}
	if left.Type() == object.BOOLEAN_OBJ {
		return evalForBoolean(left, operator, right)
	}
	return object.NULL
}

func evalForInteger(left object.Object, operator string, right object.Object) object.Object {
	l := left.(*object.Integer)
	r := right.(*object.Integer)
	switch operator {
	case "+":
		return &object.Integer{Value: l.Value + r.Value}
	case "-":
		return &object.Integer{Value: l.Value - r.Value}
	case "*":
		return &object.Integer{Value: l.Value * r.Value}
	case "/":
		return &object.Integer{Value: l.Value / r.Value}
	case ">":
		return getBoolRef(l.Value > r.Value)
	case "<":
		return getBoolRef(l.Value < r.Value)
	case "==":
		return getBoolRef(l.Value == r.Value)
	case "!=":
		return getBoolRef(l.Value != r.Value)
	}
	return object.NULL
}

func evalForBoolean(left object.Object, operator string, right object.Object) object.Object {
	l := left.(*object.Boolean)
	r := right.(*object.Boolean)
	switch operator {
	case "==":
		return getBoolRef(l.Value == r.Value)
	case "!=":
		return getBoolRef(l.Value != r.Value)
	}
	return object.NULL
}

func evalYediMujiStatement(condition ast.Node, consequent ast.Node, alternative ast.Node) object.Object {
	cc := Eval(condition)
	if cc == nil {
		return nil
	}

	if isTruthy(cc) {
		b, ok := consequent.(*ast.BlockStatement)
		if !ok {
			// maybe we will support non block statements in the future
			return nil
		}
		return evalStatements(b.Statements)
	} else {
		b, ok := alternative.(*ast.BlockStatement)
		if !ok {
			return nil
		}
		if b == nil {
			return object.NULL
		}
		return evalStatements(b.Statements)
	}
}

func isTruthy(o object.Object) bool {
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

func getBoolRef(x bool) *object.Boolean {
	if x {
		return object.TRUE
	}
	return object.FALSE
}
