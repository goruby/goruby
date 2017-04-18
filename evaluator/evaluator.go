package evaluator

import (
	"fmt"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/object"
)

type callContext struct {
	object.CallContext
}

func (c *callContext) Eval(node ast.Node, env object.Environment) (object.RubyObject, error) {
	return Eval(node, env)
}

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
		if err != nil {
			return nil, err
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

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
		}
		object.AddMethod(context, node.Name.Value, function)
		return &object.Symbol{node.Name.Value}, nil
	case *ast.ArrayLiteral:
		elements, err := evalExpressions(node.Elements, env)
		if err != nil {
			return nil, err
		}
		return &object.Array{Elements: elements}, nil

	// Expressions
	case *ast.VariableAssignment:
		val, err := Eval(node.Value, env)
		if err != nil {
			return nil, err
		}
		env.Set(node.Name.Value, val)
		return val, nil
	case *ast.ModuleExpression:
		module := object.NewModule(node.Name.Value, nil)
		env.Set("self", &object.Self{module, node.Name.Value})
		defer env.Unset("self")
		bodyReturn, err := Eval(node.Body, env)
		if err != nil {
			return nil, err
		}
		selfObject, _ := env.Get("self")
		self := selfObject.(*object.Self)
		env.SetGlobal(node.Name.Value, self.RubyObject)
		return bodyReturn, nil
	case *ast.ContextCallExpression:
		context, err := Eval(node.Context, env)
		if err != nil {
			return nil, err
		}
		if context == nil {
			context, _ = env.Get("self")
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return nil, err
		}
		callContext := &callContext{object.NewCallContext(env, context)}
		return object.Send(callContext, node.Function.Value, args...)
	case *ast.IndexExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, err
		}
		index, err := Eval(node.Index, env)
		if err != nil {
			return nil, err
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, err
		}

		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case nil:
		return nil, nil
	default:
		err := object.NewException("Unknown AST: %T", node)
		return nil, err
	}

}

func evalProgram(stmts []ast.Statement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range stmts {
		result, err = Eval(statement, env)

		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		}

	}
	return result, nil
}

func evalExpressions(exps []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	for _, e := range exps {
		evaluated, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalPrefixExpression(operator string, right object.RubyObject) (object.RubyObject, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil, object.NewException("unknown operator: %s%s", operator, right.Type())
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
		return nil, object.NewException("unknown operator: -%s", right.Type())
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
		return nil, object.NewException("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return nil, object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return nil, object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.RubyObject,
) (object.RubyObject, error) {
	if operator != "+" {
		return nil, object.NewException("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}, nil
}

func evalIfExpression(ie *ast.IfExpression, env object.Environment) (object.RubyObject, error) {
	condition, err := Eval(ie.Condition, env)
	if err != nil {
		return nil, err
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
		return nil, object.NewException("index operator not supported: %s", left.Type())
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
		if err != nil {
			return nil, err
		}
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ {
				return result, nil
			}

		}
	}
	if result == nil {
		return object.NIL, nil
	}
	return result, nil
}

func evalIdentifier(node *ast.Identifier, env object.Environment) (object.RubyObject, error) {
	val, ok := env.Get(node.Value)
	if ok {
		return val, nil
	}

	if node.IsConstant() {
		return nil, object.NewUninitializedConstantNameError(node.Value)
	}

	self, _ := env.Get("self")
	context := &callContext{object.NewCallContext(env, self)}
	val, err := object.Send(context, node.Value)
	if err != nil {
		return nil, object.NewUndefinedLocalVariableOrMethodNameError(self, node.Value)
	}
	return val, nil
}

func applyFunction(fn object.CallContext, args []object.RubyObject) (object.RubyObject, error) {
	receiver := fn.Receiver()
	switch fn := receiver.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return nil, object.NewWrongNumberOfArgumentsError(len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated, err := Eval(fn.Body, extendedEnv)
		if err != nil {
			return nil, err
		}
		return unwrapReturnValue(evaluated), nil
	default:
		return nil, object.NewSyntaxError(fmt.Errorf("not a function: %s", fn.Type()))
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
