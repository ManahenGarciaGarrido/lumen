package evaluator

import (
	"fmt"

	"github.com/ManahenGarciaGarrido/lumen/ast"
	"github.com/ManahenGarciaGarrido/lumen/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval recursively evaluates an AST node within the given environment.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.AssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if !env.Update(node.Name.Value, val) {
			return newError("RuntimeError — undefined variable %q", node.Name.Value)
		}
		return val

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBool(node.Value)

	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elements, env)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		idx := Eval(node.Index, env)
		if isError(idx) {
			return idx
		}
		return evalIndexExpression(left, idx)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	}

	return newError("unknown node type: %T", node)
}

func evalProgram(prog *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range prog.Statements {
		result = Eval(stmt, env)
		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("RuntimeError — undefined variable %q", node.Value)
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalNegateExpression(right)
	}
	return newError("RuntimeError — unknown prefix operator: %s%s", op, right.Type())
}

func evalBangExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegateExpression(right object.Object) object.Object {
	i, ok := right.(*object.Integer)
	if !ok {
		return newError("RuntimeError — operator - not supported for %s", right.Type())
	}
	return &object.Integer{Value: -i.Value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfix(op, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfix(op, left, right)
	case left.Type() == object.STRING_OBJ && op == "+":
		return &object.String{Value: left.(*object.String).Value + right.Inspect()}
	case right.Type() == object.STRING_OBJ && op == "+":
		return &object.String{Value: left.Inspect() + right.(*object.String).Value}
	case op == "==":
		return nativeBool(left == right)
	case op == "!=":
		return nativeBool(left != right)
	case op == "&&":
		return nativeBool(isTruthy(left) && isTruthy(right))
	case op == "||":
		return nativeBool(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return newError("RuntimeError — type mismatch: %s %s %s", left.Type(), op, right.Type())
	}
	return newError("RuntimeError — unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalIntegerInfix(op string, left, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		if r == 0 {
			return newError("RuntimeError — division by zero")
		}
		return &object.Integer{Value: l / r}
	case "%":
		if r == 0 {
			return newError("RuntimeError — modulo by zero")
		}
		return &object.Integer{Value: l % r}
	case "<":
		return nativeBool(l < r)
	case ">":
		return nativeBool(l > r)
	case "<=":
		return nativeBool(l <= r)
	case ">=":
		return nativeBool(l >= r)
	case "==":
		return nativeBool(l == r)
	case "!=":
		return nativeBool(l != r)
	}
	return newError("RuntimeError — unknown integer operator: %s", op)
}

func evalStringInfix(op string, left, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: l + r}
	case "==":
		return nativeBool(l == r)
	case "!=":
		return nativeBool(l != r)
	case "<":
		return nativeBool(l < r)
	case ">":
		return nativeBool(l > r)
	}
	return newError("RuntimeError — operator %s not supported for strings", op)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	}
	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exprs {
		val := Eval(e, env)
		if isError(val) {
			return []object.Object{val}
		}
		result = append(result, val)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError("RuntimeError — wrong number of arguments: want %d, got %d",
				len(fn.Parameters), len(args))
		}
		extended := object.NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			extended.Set(param.Value, args[i])
		}
		result := Eval(fn.Body, extended)
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value
		}
		return result
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return newError("RuntimeError — not a function: %s", fn.Type())
}

func evalIndexExpression(left, idx object.Object) object.Object {
	if left.Type() == object.ARRAY_OBJ && idx.Type() == object.INTEGER_OBJ {
		arr := left.(*object.Array).Elements
		i := idx.(*object.Integer).Value
		if i < 0 || i >= int64(len(arr)) {
			return NULL
		}
		return arr[i]
	}
	return newError("RuntimeError — index operator not supported for %s", left.Type())
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBool(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
