package evaluator

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

// Eval evaluates the given node and traverses recursive over its children
func Eval(node ast.Node, env object.Environment) object.RubyObject {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions

	// Literals
	case (*ast.IntegerLiteral):
		return object.NewInteger(node.Value)
	case (*ast.Boolean):
		return nativeBoolToBooleanObject(node.Value)
	case (*ast.Nil):
		return object.NIL
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.SymbolLiteral:
		return &object.Symbol{Value: node.Value}
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		function := &object.Function{Parameters: params, Env: env, Body: body}
		env.Set(node.Name.Value, function)
		return function
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.VariableAssignment:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val
	case *ast.ContextCallExpression:
		context := Eval(node.Context, env)
		if IsError(context) {
			return context
		}
		if context == nil {
			context, _ = env.Get("self")
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		if node.Context == nil {
			function := Eval(node.Function, env)
			if IsError(function) {
				return function
			}
			return applyFunction(function, args)
		}
		return object.Send(context, node.Function.Value, args...)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if IsError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.RequireExpression:
		return evalRequireExpression(node, env)
	case nil:
		return nil
	default:
		return object.NewException("Unknown AST: %T", node)
	}

}

func evalProgram(stmts []ast.Statement, env object.Environment) object.RubyObject {
	var result object.RubyObject
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Builtin:
			return result.Fn()
		}

		if IsError(result) {
			return result
		}
	}
	return result
}

func evalExpressions(exps []ast.Expression, env object.Environment) []object.RubyObject {
	var result []object.RubyObject

	for _, e := range exps {
		evaluated := Eval(e, env)
		if IsError(evaluated) {
			return []object.RubyObject{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalRequireExpression(expr *ast.RequireExpression, env object.Environment) object.RubyObject {
	filename := expr.Name.Value
	if !strings.HasSuffix(filename, "rb") {
		filename += ".rb"
	}
	loadedFeatures, ok := env.Get("$LOADED_FEATURES")
	if !ok {
		loadedFeatures = object.NewArray()
		env.SetGlobal("$LOADED_FEATURES", loadedFeatures)
	}
	arr, ok := loadedFeatures.(*object.Array)
	if !ok {
		arr = object.NewArray()
	}
	loaded := false
	for _, feat := range arr.Elements {
		if feat.Inspect() == filename {
			loaded = true
			break
		}
	}
	if loaded {
		return object.FALSE
	}

	arr.Elements = append(arr.Elements, &object.String{Value: filename})
	file, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return object.NewLoadError(expr.Name.Value)
	}
	l := lexer.New(string(file))
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return object.NewSyntaxError(err.Error())
	}
	Eval(prog, env)
	return object.TRUE
}

func evalPrefixExpression(operator string, right object.RubyObject) object.RubyObject {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.NewException("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.RubyObject) object.RubyObject {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NIL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.RubyObject) object.RubyObject {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	default:
		return object.NewException("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.RubyObject) object.RubyObject {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return object.NewException("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.RubyObject) object.RubyObject {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.RubyObject) object.RubyObject {
	if operator != "+" {
		return object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIfExpression(ie *ast.IfExpression, env object.Environment) object.RubyObject {
	condition := Eval(ie.Condition, env)
	if IsError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.NIL
	}
}

func evalIndexExpression(left, index object.RubyObject) object.RubyObject {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return object.NewException("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.RubyObject) object.RubyObject {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return object.NIL
	}
	return arrayObject.Elements[idx]
}

func evalBlockStatement(block *ast.BlockStatement, env object.Environment) object.RubyObject {
	var result object.RubyObject
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || IsError(result) {
				return result
			}

		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env object.Environment) object.RubyObject {
	val, ok := env.Get(node.Value)
	if ok {
		if fn, ok := val.(*object.Function); ok {
			if len(fn.Parameters) != 0 {
				return val
			}
			return applyFunction(fn, []object.RubyObject{})
		}
		return val
	}
	self, _ := env.Get("self")
	val = object.Send(self, node.Value)
	if IsError(val) {
		return object.NewNameError(self, node.Value)
	}
	return val
}

func applyFunction(fn object.RubyObject, args []object.RubyObject) object.RubyObject {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return object.NewWrongNumberOfArgumentsError(len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return object.NewSyntaxError(fmt.Sprintf("not a function: %s", fn.Type()))
	}
}

func extendFunctionEnv(fn *object.Function, args []object.RubyObject) object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.RubyObject) object.RubyObject {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.RubyObject) bool {
	switch obj {
	case object.NIL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

// IsError returns true if the given RubyObject is an object.Error or an
// object.Exception (or any subclass of object.Exception)
func IsError(obj object.RubyObject) bool {
	if obj != nil {
		return obj.Type() == object.EXCEPTION_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) object.RubyObject {
	if input {
		return object.TRUE
	}
	return object.FALSE
}
