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
func Eval(node ast.Node, env object.Environment) (object.RubyObject, error) {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val, err := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val, err
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions

	// Literals
	case (*ast.IntegerLiteral):
		return object.NewInteger(node.Value), nil
	case (*ast.Boolean):
		return nativeBoolToBooleanObject(node.Value), nil
	case (*ast.Nil):
		return object.NIL, nil
	case (*ast.Self):
		self, _ := env.Get("self")
		return self, nil
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}, nil
	case *ast.SymbolLiteral:
		return &object.Symbol{Value: node.Value}, nil
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		context, _ := env.Get("self")
		function := &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
			CallFn:     applyFunction,
		}
		object.AddMethod(context, node.Name.Value, function)
		return function, nil
	case *ast.ArrayLiteral:
		elements, err := evalExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0], err
		}
		return &object.Array{Elements: elements}, nil
	case *ast.VariableAssignment:
		val, err := Eval(node.Value, env)
		if IsError(val) {
			return val, err
		}
		env.Set(node.Name.Value, val)
		return val, nil
	case *ast.ContextCallExpression:
		context, err := Eval(node.Context, env)
		if IsError(context) {
			return context, err
		}
		if context == nil {
			context, _ = env.Get("self")
		}
		args, err := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0], err
		}
		if function, ok := env.Get(node.Function.Value); ok {
			return applyFunction(function, args), nil
		}
		return object.Send(context, node.Function.Value, args...), nil
	case *ast.IndexExpression:
		left, err := Eval(node.Left, env)
		if IsError(left) {
			return left, err
		}
		index, err := Eval(node.Index, env)
		if IsError(index) {
			return index, err
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		right, err := Eval(node.Right, env)
		if IsError(right) {
			return right, err
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, err := Eval(node.Left, env)
		if IsError(left) {
			return left, err
		}

		right, err := Eval(node.Right, env)
		if IsError(right) {
			return right, err
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.RequireExpression:
		return evalRequireExpression(node, env)
	case nil:
		return nil, nil
	default:
		err := object.NewException("Unknown AST: %T", node)
		return err, err
	}

}

func evalProgram(stmts []ast.Statement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range stmts {
		result, err = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		case *object.Builtin:
			return result.Fn(), nil
		}

		if IsError(result) {
			return result, err
		}
	}
	return result, nil
}

func evalExpressions(exps []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	for _, e := range exps {
		evaluated, err := Eval(e, env)
		if IsError(evaluated) {
			return []object.RubyObject{evaluated}, err
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalRequireExpression(expr *ast.RequireExpression, env object.Environment) (object.RubyObject, error) {
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
		return object.FALSE, nil
	}

	arr.Elements = append(arr.Elements, &object.String{Value: filename})
	file, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		errObj := object.NewLoadError(expr.Name.Value)
		return errObj, errObj
	}
	l := lexer.New(string(file))
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		errObj := object.NewSyntaxError(err.Error())
		return errObj, errObj
	}
	Eval(prog, env)
	return object.TRUE, nil
}

func evalPrefixExpression(operator string, right object.RubyObject) (object.RubyObject, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		err := object.NewException("unknown operator: %s%s", operator, right.Type())
		return err, err
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

func evalMinusPrefixOperatorExpression(right object.RubyObject) (object.RubyObject, error) {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}, nil
	default:
		err := object.NewException("unknown operator: -%s", right.Type())
		return err, err
	}
}

func evalInfixExpression(operator string, left, right object.RubyObject) (object.RubyObject, error) {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right), nil
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right), nil
	case left.Type() != right.Type():
		err := object.NewException("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		return err, err
	default:
		err := object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		return err, err
	}
}

func evalIntegerInfixExpression(operator string, left, right object.RubyObject) (object.RubyObject, error) {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}, nil
	case "-":
		return &object.Integer{Value: leftVal - rightVal}, nil
	case "*":
		return &object.Integer{Value: leftVal * rightVal}, nil
	case "/":
		return &object.Integer{Value: leftVal / rightVal}, nil
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal), nil
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal), nil
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal), nil
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal), nil
	default:
		err := object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		return err, err
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.RubyObject,
) (object.RubyObject, error) {
	if operator != "+" {
		err := object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		return err, err
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}, nil
}

func evalIfExpression(ie *ast.IfExpression, env object.Environment) (object.RubyObject, error) {
	condition, err := Eval(ie.Condition, env)
	if IsError(condition) {
		return condition, err
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.NIL, nil
	}
}

func evalIndexExpression(left, index object.RubyObject) (object.RubyObject, error) {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index), nil
	default:
		err := object.NewException("index operator not supported: %s", left.Type())
		return err, err
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

func evalBlockStatement(block *ast.BlockStatement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range block.Statements {
		result, err = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || IsError(result) {
				return result, err
			}

		}
	}
	return result, nil
}

func evalIdentifier(node *ast.Identifier, env object.Environment) (object.RubyObject, error) {
	val, ok := env.Get(node.Value)
	if ok {
		if fn, ok := val.(*object.Function); ok {
			if len(fn.Parameters) != 0 {
				return val, nil
			}
			return applyFunction(fn, []object.RubyObject{}), nil
		}
		return val, nil
	}
	self, _ := env.Get("self")
	val = object.Send(self, node.Value)
	if IsError(val) {
		err := object.NewNameError(self, node.Value)
		return err, err
	}
	return val, nil
}

func applyFunction(fn object.RubyObject, args []object.RubyObject) object.RubyObject {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return object.NewWrongNumberOfArgumentsError(len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated, _ := Eval(fn.Body, extendedEnv)
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
