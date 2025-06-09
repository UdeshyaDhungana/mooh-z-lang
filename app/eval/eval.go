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
		return evalYediMujiStatement(node, env)
	case *ast.PathaMujiStatement:
		return evalPathaMujiStatement(node.Value, env)
	case *ast.ThoosMujiStatement:
		return evalThoosMujiStatement(node.Name, node.Value, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringExpression:
		return &object.String{Value: node.Value}
	case *ast.ArrayExpression:
		return evalArrayExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.JabasammaMujiExpression:
		return evalJabasammaMujiExpression(node.Condition, node.Consequent, env)
	case *ast.GhumaMujiExpression:
		return evalGhumaMujiExpression(node.Initialization, node.Condition, node.Update, node.Body, env)
	case *ast.HashExpression:
		return evalHashExpression(node, env)
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
		/* If assignment operator, evaluate the right and assign to right */
		/* can be either x = 34, or x[24] = 53 */
		if node.Operator == "=" {
			return evalAssignment(node.Left, node.Right, env)
		}
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(left, node.Operator, right)
	}
	fmt.Printf("FATAL: Eval() does not implement %s node\n", node.String())
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
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

/* Infix begin */
func evalArithmetic(left object.Object, operator string, right object.Object) object.Object {
	if left.Type() == right.Type() && operator == "+" {
		switch l := left.(type) {
		case *object.String:
			r := right.(*object.String)
			return &object.String{Value: l.Value + r.Value}
		case *object.Array:
			r := right.(*object.Array)
			return &object.Array{Arr: append(l.Arr, r.Arr...)}
		}
	}
	if !areBothNumbers(left, right) {
		return newError("unsupported operation %s %s %s", left.Type(), operator, right.Type())
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		l := left.(*object.Integer).Value
		r := right.(*object.Integer).Value
		var res object.Integer
		switch operator {
		case "+":
			res.Value = l + r
		case "-":
			res.Value = l - r
		case "*":
			res.Value = l * r
		case "/":
			res.Value = l / r
		case "%":
			res.Value = l % r
		default:
			return newError("unsupported operation %s %s %s", left.Type(), operator, right.Type())
		}
		return &res
	}
	var leftVal, rightVal float64
	var result object.Float

	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	switch operator {
	case "+":
		result.Value = leftVal + rightVal
	case "-":
		result.Value = leftVal - rightVal
	case "*":
		result.Value = leftVal * rightVal
	case "/":
		result.Value = leftVal / rightVal
	default:
		return newError("unsupported operation %s %s %s", left.Type(), operator, right.Type())
	}
	return &result
}

func evalEQ(left object.Object, right object.Object) *object.Boolean {
	if left.Type() != right.Type() {
		return object.FALSE
	}
	switch l := left.(type) {
	case *object.Integer:
		r := right.(*object.Integer)
		if l.Value == r.Value {
			return object.TRUE
		}
		return object.FALSE
	case *object.Float:
		r := right.(*object.Float)
		if l.Value == r.Value {
			return object.TRUE
		}
		return object.FALSE
	case *object.Boolean:
		r := right.(*object.Boolean)
		if l.Value == r.Value {
			return object.TRUE
		}
		return object.FALSE
	default:
		// we can either go with checking if they are same objects
		// or take the python's approach of checking each element
		// we do neither ¯\_(ツ)_/¯
		return object.FALSE
	}
}

func evalLT(left object.Object, right object.Object) object.Object {
	// ensure both of them are either float or int
	if !areBothNumbers(left, right) {
		return newError("cannot use '<' operator for %s", left.Type())
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		l := left.(*object.Integer).Value
		r := right.(*object.Integer).Value
		return utils.GetBoolRef(l < r)
	}
	var leftVal, rightVal float64
	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = (left.(*object.Float).Value)
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	return utils.GetBoolRef(leftVal < rightVal)
}

func evalGT(left object.Object, right object.Object) object.Object {
	if !areBothNumbers(left, right) {
		return newError("cannot use '>' operator for %s", left.Type())
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		l := left.(*object.Integer).Value
		r := right.(*object.Integer).Value
		return utils.GetBoolRef(l > r)
	}
	var leftVal, rightVal float64
	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = (left.(*object.Float).Value)
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	return utils.GetBoolRef(leftVal > rightVal)
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch operator {
	case "+", "-", "*", "/", "%":
		return evalArithmetic(left, operator, right)
	case "==":
		return evalEQ(left, right)
	case "!=":
		eq := evalEQ(left, right)
		if eq == object.TRUE {
			return object.FALSE
		}
		return object.TRUE
	case ">":
		return evalGT(left, right)
	case ">=":
		isGT := evalGT(left, right)
		if isGT.Type() == object.GALAT_MUJI_OBJ {
			return isGT
		}
		if isGT == object.TRUE {
			return isGT
		}
		isEQ := evalEQ(left, right)
		if isEQ.Type() == object.GALAT_MUJI_OBJ {
			return isEQ
		}
		if isEQ == object.TRUE {
			return isEQ
		}
		return object.FALSE
	case "<":
		return evalLT(left, right)
	case "<=":
		isLT := evalLT(left, right)
		if isLT.Type() == object.GALAT_MUJI_OBJ {
			return isLT
		}
		if isLT == object.TRUE {
			return isLT
		}
		isEQ := evalEQ(left, right)
		if isEQ.Type() == object.GALAT_MUJI_OBJ {
			return isEQ
		}
		if isEQ == object.TRUE {
			return object.TRUE
		}
		return object.FALSE
	default:
		return newError("unsupported operator %s", operator)
	}
}

func areBothNumbers(left object.Object, right object.Object) bool {
	return (left.Type() == object.INTEGER_OBJ || left.Type() == object.FLOAT_OBJ) &&
		(right.Type() == object.INTEGER_OBJ || right.Type() == object.FLOAT_OBJ)
}

/* End Infix */

func evalYediMujiStatement(yediMujiExpr *ast.YediMujiExpression, env *object.Environment) object.Object {
	if isConditionTrue(yediMujiExpr.Condition, env) {
		return evalStatements(yediMujiExpr.Consequent.Statements, env)
	}

	// Evaluate alternatives now
	if yediMujiExpr.Alternatives != nil {
		for _, alt := range yediMujiExpr.Alternatives {
			if isConditionTrue(alt.Condition, env) {
				return evalStatements(alt.Consequent.Statements, env)
			}
		}
	}

	if yediMujiExpr.Fallback != nil {
		blockEnv := object.NewEnclosedEnvironment(env)
		return evalStatements(yediMujiExpr.Fallback.Statements, blockEnv)
	}
	return object.NULL
}

func evalJabasammaMujiExpression(condition ast.Node, consequent ast.Node, env *object.Environment) object.Object {
	var result object.Object
	newEnv := object.NewEnclosedEnvironment(env)
	for isConditionTrue(condition, newEnv) {
		body, ok := consequent.(*ast.BlockStatement)
		if !ok {
			return newError("body of jaba samma muji must be a block statement")
		}
		newEnv := object.NewEnclosedEnvironment(newEnv)
		result = Eval(body, newEnv)
	}
	return result
}

func evalGhumaMujiExpression(initialization ast.Node, condition ast.Node, update ast.Node, body ast.Node, env *object.Environment) object.Object {
	newEnv := object.NewEnclosedEnvironment(env)
	init := Eval(initialization, newEnv)
	if init.Type() == object.GALAT_MUJI_OBJ {
		return init
	}
	for isConditionTrue(condition, newEnv) {
		// check condition
		result := Eval(body, newEnv)
		if result.Type() == object.GALAT_MUJI_OBJ || result.Type() == object.PATHA_MUJI_OBJ {
			return result
		}
		up := Eval(update, newEnv)
		if up.Type() == object.GALAT_MUJI_OBJ {
			return result
		}
	}
	return object.NULL
}

// true should be checked against err == nil && bool == true
// the only case where bool = false and err == nil is when Eval() does not implement case for the node
func isConditionTrue(condition ast.Node, env *object.Environment) bool {
	cc := Eval(condition, env)
	if cc == nil {
		return false
	}
	if isError(cc) {
		return false
	}
	return utils.IsTruthy(cc)
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
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
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

func evalAssignment(name ast.Expression, value ast.Expression, env *object.Environment) object.Object {
	if name == nil {
		return newError("left operand of assignment operator is nil")
	}
	v := Eval(value, env)
	if v.Type() == object.GALAT_MUJI_OBJ {
		return v
	}
	switch n := name.(type) {
	case *ast.Identifier:
		_, ok := env.Get(n.Value)
		if !ok {
			return newError("reassignment to an undefined variable %s", n.Value)
		}
		/* Get the environment in which the varible lies */
		e := env.GetEnv(n.Value)
		return e.Set(n.Value, v)
	case *ast.IndexExpression:
		return evalAssignmentForIndexExpression(n, v, env)
	default:
		return newError("left operand of assignment operator is neither identifier nor indexexpression. got=%T", name)
	}
}

func evalAssignmentForIndexExpression(ie *ast.IndexExpression, value object.Object, env *object.Environment) object.Object {
	operand := Eval(ie.Operand, env)
	index := Eval(ie.Index, env)
	switch operand.Type() {
	case object.ARRAY_OBJECT:
		if index.Type() != object.INTEGER_OBJ {
			return newError("cannot index an array using %s", index.Type())
		}
		a := operand.(*object.Array)
		i := index.(*object.Integer)
		if i.Value >= int64(len(a.Arr)) {
			return newError("array index out of bounds")
		}
		a.Arr[i.Value] = value
		return value
	case object.HASHMAP_OBJECT:
		if index.Type() != object.STRING {
			return newError("cannot index a hashmap using %s", index.Type())
		}
		h := operand.(*object.HashMap)
		i := index.(*object.String)
		h.Pairs[i.Value] = value
		return value
	default:
		return newError("only arrays and hashmaps can be indexed")
	}
}

func evalKaamGarMujiExpression(node *ast.KaamGarMujiExpression, env *object.Environment) object.Object {
	var result object.KaamGar
	result.Body = node.Body
	result.Env = env
	result.Parameters = node.Arguments
	return &result
}

func evalCallExpression(name *ast.CallExpression, env *object.Environment) object.Object {
	// evaluate arguments
	var evaluatedArgs []*object.Object
	for _, v := range name.Arguments {
		e := Eval(v, env)
		if isError(e) {
			return e
		}
		evaluatedArgs = append(evaluatedArgs, &e)
	}

	fn := evalIdentifier(&ast.Identifier{Value: name.Function.TokenLiteral()}, env)
	switch fn.Type() {
	case object.GALAT_MUJI_OBJ:
		// might be a function expression itself
		if kgr, ok := name.Function.(*ast.KaamGarMujiExpression); ok {
			fn = evalKaamGarMujiExpression(kgr, env)
		} else {
			return newError("Cannot apply: %s", name.Function.String())
		}
	case object.BUILTIN_OBJECT:
		f := fn.(*object.Builtin)
		return evalBuiltin(f, evaluatedArgs)
	}

	if f, ok := fn.(*object.KaamGar); ok {
		return evalUserDefinedCall(f, evaluatedArgs)
	}
	return newError("cannot apply %s; not a function or a builtin", name.Function.TokenLiteral())
}

func evalBuiltin(b *object.Builtin, args []*object.Object) object.Object {
	converted := make([]object.Object, len(args))
	for i, v := range args {
		converted[i] = *v
	}
	return b.Fn(converted...)
}

func evalUserDefinedCall(f *object.KaamGar, args []*object.Object) object.Object {
	// check parameters
	if len(f.Parameters) != len(args) {
		return newError("arguments length mismatch")
	}

	// set it for function's environment
	f.Env = object.NewEnclosedEnvironment(f.Env)
	for i, v := range f.Parameters {
		f.Env.Set(v.Value, *args[i])
	}
	result := Apply(f)
	f.Env = f.Env.PopStack()
	return result
}

func evalArrayExpression(a *ast.ArrayExpression, env *object.Environment) object.Object {
	var result object.Array
	for i := range a.Elements {
		evaluated := Eval(a.Elements[i], env)
		if evaluated.Type() == object.GALAT_MUJI_OBJ {
			return evaluated
		}
		result.Arr = append(result.Arr, evaluated)
	}
	return &result
}

func evalIndexExpression(a *ast.IndexExpression, env *object.Environment) object.Object {
	idxEvaluated := Eval(a.Index, env)
	if idxEvaluated.Type() == object.GALAT_MUJI_OBJ {
		return idxEvaluated
	}
	operand := Eval(a.Operand, env)
	switch operand.Type() {
	case object.ARRAY_OBJECT:
		if idxEvaluated.Type() != object.INTEGER_OBJ {
			return newError("array index must be an integer, got %s", idxEvaluated.Type())
		}
		arr := operand.(*object.Array)
		idx := idxEvaluated.(*object.Integer)
		return arr.Arr[idx.Value]
	case object.HASHMAP_OBJECT:
		if idxEvaluated.Type() != object.STRING {
			return newError("hashmap index must be a string, got %s", idxEvaluated.Type())
		}
		hmap := operand.(*object.HashMap)
		idx := idxEvaluated.(*object.String)
		return hmap.Pairs[idx.Value]
	case object.GALAT_MUJI_OBJ:
		return operand
	default:
		return newError("cannot index %s", operand.Type())
	}
}

func evalHashExpression(node *ast.HashExpression, env *object.Environment) object.Object {
	pairs := node.Pairs
	result := object.HashMap{Pairs: make(map[string]object.Object)}
	for k, v := range pairs {
		key := Eval(k, env)
		if key.Type() != object.STRING {
			return newError("key must be a string")
		}
		keyStr := key.(*object.String)
		val := Eval(v, env)
		result.Pairs[keyStr.Value] = val
	}
	return &result
}

func Apply(f *object.KaamGar) object.Object {
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
