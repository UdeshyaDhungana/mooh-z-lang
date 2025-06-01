package eval

import (
	"fmt"

	"github.com/udeshyadhungana/interprerer/app/ast"
	"github.com/udeshyadhungana/interprerer/app/object"
	"github.com/udeshyadhungana/interprerer/app/utils"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.YediMujiExpression:
		return evalYediMujiStatement(node.Condition, node.Consequent, node.Alternative, env)
	case *ast.PathaMujiStatement:
		return evalPathaMujiStatement(node.Value, env)
	case *ast.ThoosMujiStatement:
		return evalThoosMujiStatement(node.Name, node.Value, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringExpression:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return object.TRUE
		}
		return object.FALSE
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.KaamGarMujiExpression:
		return evalKaamGarMujiExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(left, node.Operator, right)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statment := range program.Statements {
		result = Eval(statment, env)
		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.PATHA_MUJI_OBJ || rt == object.GALAT_MUJI_OBJ {
				return result
			}
		}
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
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer)
	return &object.Integer{Value: -value.Value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	if left.Type() == object.INTEGER_OBJ {
		return evalForInteger(left, operator, right)
	}
	if left.Type() == object.BOOLEAN_OBJ {
		return evalForBoolean(left, operator, right)
	}
	if left.Type() == object.STRING {
		return evalForString(left, operator, right)
	}
	return newError("type unsupported: %s %s %s", left.Type(), operator, right.Type())
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
		return utils.GetBoolRef(l.Value > r.Value)
	case "<":
		return utils.GetBoolRef(l.Value < r.Value)
	case "==":
		return utils.GetBoolRef(l.Value == r.Value)
	case "!=":
		return utils.GetBoolRef(l.Value != r.Value)
	}
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalForBoolean(left object.Object, operator string, right object.Object) object.Object {
	l := left.(*object.Boolean)
	r := right.(*object.Boolean)
	switch operator {
	case "==":
		return utils.GetBoolRef(l.Value == r.Value)
	case "!=":
		return utils.GetBoolRef(l.Value != r.Value)
	}
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalForString(left object.Object, operator string, right object.Object) object.Object {
	l := left.(*object.String)
	r := right.(*object.String)
	switch operator {
	case "+":
		return &object.String{Value: l.Value + r.Value}
	}
	return newError("unkown operator: %s %s %s", l.Type(), operator, r.Type())
}

func evalYediMujiStatement(condition ast.Node, consequent ast.Node, alternative ast.Node, env *object.Environment) object.Object {
	cc := Eval(condition, env)
	if cc == nil {
		return nil
	}
	if isError(cc) {
		return cc
	}

	if utils.IsTruthy(cc) {
		b, ok := consequent.(*ast.BlockStatement)
		if !ok {
			// maybe we will support non block statements in the future
			return nil
		}
		return evalStatements(b.Statements, env)
	} else {
		b, ok := alternative.(*ast.BlockStatement)
		if !ok {
			return nil
		}
		if b == nil {
			return object.NULL
		}
		return evalStatements(b.Statements, env)
	}
}

func evalPathaMujiStatement(value ast.Node, env *object.Environment) object.Object {
	val := Eval(value, env)
	if isError(val) {
		return val
	}
	return &object.Return{
		Value: Eval(value, env),
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalThoosMujiStatement(name *ast.Identifier, value ast.Expression, env *object.Environment) object.Object {
	if name != nil {
		v := Eval(value, env)
		if v.Type() == object.GALAT_MUJI_OBJ {
			return v
		}
		return env.Set(name.Value, v)
	}
	return newError("identifier is nil")
}

func evalKaamGarMujiExpression(node *ast.KaamGarMujiExpression, env *object.Environment) object.Object {
	var result object.KaamGar
	result.Body = node.Body
	result.Env = env
	result.Parameters = node.Arguments
	return &result
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	var fn object.Object
	fn, ok := env.Get(node.Function.TokenLiteral())
	if !ok {
		kgr, ok := node.Function.(*ast.KaamGarMujiExpression)
		if !ok {
			return newError("cannot apply expression: %s", node.Function.String())
		}
		fn = evalKaamGarMujiExpression(kgr, env)
	} else {
		if fn.Type() != object.KAAM_GAR_MUJI_OBJ {
			return newError("identifier not found: %s", node.Function.TokenLiteral())
		}
	}
	f, ok := fn.(*object.KaamGar)
	if !ok {
		return newError("cannot apply %s; not a function", node.Function.TokenLiteral())
	}

	if len(f.Parameters) != len(node.Arguments) {
		return newError("arguments length mismatch for %s", node.TokenLiteral())
	}
	f.Env = object.NewEnclosedEnvironment(f.Env)
	for i, v := range f.Parameters {
		evaluatedArg := Eval(node.Arguments[i], env)
		if isError(evaluatedArg) {
			return evaluatedArg
		}
		f.Env.Set(v.Value, evaluatedArg)
	}
	result := Apply(f, node)
	f.Env = f.Env.PopStack()
	return result
}

func Apply(f *object.KaamGar, callExp *ast.CallExpression) object.Object {
	res := Eval(f.Body, f.Env)

	if returnValue, ok := res.(*object.Return); ok {
		return returnValue.Value
	}
	return res
}

/* Utils */
func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.GALAT_MUJI_OBJ
	}
	return false
}
