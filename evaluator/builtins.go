package evaluator

import (
	"bufio"
	"fmt"
	"monkey/object"
	"os"
	"strconv"
)

func b_puts(args ...object.Object) object.Object {

	if len(args) == 0 {
		return newError(fmt.Sprintf("Invalid number of arguments, want at least 1, got %d", len(args)))
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			fmt.Println(arg.Value)
		case *object.Integer:
			fmt.Println(arg.Value)
		case *object.Boolean:
			fmt.Println(arg.Value)
		}
	}

	return NULL
}

func b_read(args ...object.Object) object.Object {

	if len(args) != 1 {
		return newError(fmt.Sprintf("Invalid number of arguments, want 1, got %d", len(args)))
	}

	b_puts(args...)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		return newError("could not read the user input, " + err.Error())
	}

	return &object.String{Value: scanner.Text()}
}

func b_len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(fmt.Sprintf("Invalid number of arguments, want 1, got %d", len(args)))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("len only supports string and array arguments")
	}
}

func b_int(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(fmt.Sprintf("Invalid number of arguments, want 1, got %d", len(args)))
	}

	switch arg := args[0].(type) {
	case *object.String:
		value, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			return newError("could not convert the string to an integer, " + err.Error())
		}
		return &object.Integer{Value: value}
	case *object.Integer:
		return arg
	default:
		return newError("Invalid argument, want a string, got " + arg.Type())
	}
}

func b_str(args ...object.Object) object.Object {

	if len(args) != 1 {
		return newError(fmt.Sprintf("Invalid number of arguments, want 1, got %d", len(args)))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		value := strconv.FormatInt(arg.Value, 10)
		return &object.String{Value: value}
	default:
		return newError("argument type not supported")
	}

}

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: b_len,
	},
	"puts": &object.Builtin{
		Fn: b_puts,
	},
	"read": &object.Builtin{
		Fn: b_read,
	},
	"int": &object.Builtin{
		Fn: b_int,
	},
	"str": &object.Builtin{
		Fn: b_str,
	},
}
