package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/parser"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

var libsEnv map[string]*object.Environment

func newError(errorMsg string) *object.Error {
	return &object.Error{Message: errorMsg}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.UseStatement:
		return evalUseStatement(node.Filename)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Array:
		arr := &object.Array{}

		arr.Elements = evalExpressions(node.Elements, env)

		return arr
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		switch node.Operator {
		case "!":
			if right == TRUE {
				return FALSE
			} else if right == FALSE {
				return TRUE
			} else {
				return FALSE
			}
		case "-":
			if right.Type() != object.INTEGER_OBJ {
				return newError("expected the right member to be an integer, got a " + right.Type() + " instead")
			}

			originalValue := right.(*object.Integer).Value
			return &object.Integer{Value: -originalValue}
		default:
			return newError("unknown operator: " + node.Operator)
		}
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
			leftInt := left.(*object.Integer).Value
			rightInt := right.(*object.Integer).Value

			switch node.Operator {
			case "+":
				return &object.Integer{Value: leftInt + rightInt}
			case "-":
				return &object.Integer{Value: leftInt - rightInt}
			case "*":
				return &object.Integer{Value: leftInt * rightInt}
			case "/":
				return &object.Integer{Value: leftInt / rightInt}
			case "==":
				if leftInt == rightInt {
					return TRUE
				} else {
					return FALSE
				}
			case "!=":
				if leftInt == rightInt {
					return FALSE
				} else {
					return TRUE
				}
			case ">":
				if leftInt > rightInt {
					return TRUE
				} else {
					return FALSE
				}
			case "<":
				if leftInt < rightInt {
					return TRUE
				} else {
					return FALSE
				}
			}
		} else if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
			switch node.Operator {
			case "==":
				if left == right {
					return TRUE
				} else {
					return FALSE
				}
			case "!=":
				if left != right {
					return TRUE
				} else {
					return FALSE
				}
			}
		} else if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
			leftString := left.(*object.String).Value
			rightString := right.(*object.String).Value

			switch node.Operator {
			case "+":
				return &object.String{Value: leftString + rightString}
			case "==":
				if leftString == rightString {
					return TRUE
				} else {
					return FALSE
				}
			case "!=":
				if leftString != rightString {
					return TRUE
				} else {
					return FALSE
				}
			}
		} else {
			return newError("left and right values have different types")
		}
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if condition != NULL && condition != FALSE {
			return Eval(&node.TrueBlock, env)
		} else if len(node.FalseBlock.Statements) > 0 {
			return Eval(&node.FalseBlock, env)
		} else {
			return NULL
		}
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:

		if varExists := env.Get(node.Name.Value); varExists != nil {
			return newError("variable already declared")
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return env.Set(node.Name.Value, val)
	case *ast.Identifier:
		val := env.Get(node.Value)
		if val != nil {
			return val
		}

		if builtin, ok := builtins[node.Value]; ok {
			return builtin
		}

		return newError("invalid identifier: " + node.Value)
	case *ast.ReassignmentStatement:

		if varExists := env.Get(node.Variable.Value); varExists == nil {
			return newError("invalid identifier: " + node.Variable.Value)
		}

		val := Eval(node.NewValue, env)
		if isError(val) {
			return val
		}

		return env.Set(node.Variable.Value, val)

	case *ast.FunctionLiteral:
		funcLiteral := &object.Function{Parameters: node.Parameters, Body: node.Block, Env: *env}

		if node.Name.Value != "" {
			env.Set(node.Name.Value, funcLiteral)
		}

		return funcLiteral
	case *ast.CallExpression:

		var function object.Object

		if node.Lib != "" {
			function = Eval(node.Function, libsEnv[node.Lib])
		} else {
			function = Eval(node.Function, env)
		}

		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return callFunction(function, args)
	case *ast.ArrayAccessExpression:
		array := Eval(node.Array, env)
		if isError(array) {
			return array
		}

		if array.Type() != object.ARRAY_OBJ {
			return newError("expected left member to be an array, got " + array.Type() + " instead")
		}
		arrayElements := array.(*object.Array).Elements

		position := Eval(node.Position, env)
		if isError(position) {
			return position
		}

		if position.Type() != object.INTEGER_OBJ {
			return newError("expected position to be an integer, got " + position.Type() + " instead")
		}

		pos := position.(*object.Integer).Value

		if pos >= int64(len(arrayElements)) {
			return newError(fmt.Sprintf("index out of range, array's length is %d", len(arrayElements)))
		}

		value := arrayElements[pos]
		return value
	case *ast.WhileStatement:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		for condition != NULL && condition != FALSE {

			result := Eval(&node.Block, env)
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}

			condition = Eval(node.Condition, env)
			if isError(condition) {
				return condition
			}
		}

		return NULL
	case *ast.ExternalReferenceExpression:

		switch ref := node.Referece.(type) {
		case *ast.CallExpression:
			ref.Lib = node.Module

			result := Eval(ref, env)

			return result
		case *ast.Identifier:
			return Eval(ref, libsEnv[node.Module])
		}

	}
	return nil
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	libsEnv = make(map[string]*object.Environment)

	var result object.Object

	for _, s := range statements {
		result = Eval(s, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalUseStatement(filename string) object.Object {
	fileAst := parser.ParseFile(filename)

	fileEnv := object.NewEnvironment()

	result := Eval(fileAst, fileEnv)

	libsEnv[filename] = fileEnv

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range block.Statements {
		result = Eval(s, env)

		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func evalExpressions(nodes []ast.Expression, env *object.Environment) []object.Object {

	var arguments []object.Object

	for _, a := range nodes {
		result := Eval(a, env)

		if isError(result) {
			return []object.Object{result}
		}

		arguments = append(arguments, result)
	}

	return arguments
}

func callFunction(fn object.Object, args []object.Object) object.Object {

	if builtinFn, ok := fn.(*object.Builtin); ok {
		return builtinFn.Fn(args...)
	}

	function, ok := fn.(*object.Function)
	if !ok {
		return newError("expected a function, got " + function.Type() + " instead")
	}

	extendedEnv := object.NewExtendedEnvironment(&function.Env)

	for i, arg := range args {
		extendedEnv.Set(function.Parameters[i].Value, arg)
	}

	evaluated := Eval(&function.Body, extendedEnv)

	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluated
}
