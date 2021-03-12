package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type Object interface {
	Type() string
	Inspect() string
}

type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
)

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() string {
	return INTEGER_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Type() string {
	return BOOLEAN_OBJ
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}
func (s *String) Type() string {
	return STRING_OBJ
}

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	var result bytes.Buffer

	var elements []string
	for _, a := range a.Elements {
		elements = append(elements, a.Inspect())
	}

	result.WriteString("[")
	result.WriteString(strings.Join(elements, ", "))
	result.WriteString("]")

	return result.String()
}

func (a *Array) Type() string {
	return ARRAY_OBJ
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() string {
	return NULL_OBJ
}

type ReturnValue struct {
	Value Object
}

func (ro *ReturnValue) Type() string {
	return RETURN_VALUE_OBJ
}
func (ro *ReturnValue) Inspect() string {
	return ro.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() string {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

type Function struct {
	Parameters []ast.Identifier
	Body       ast.BlockStatement
	Env        Environment
}

func (f *Function) Type() string {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var result bytes.Buffer

	parameters := []string{}
	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}

	result.WriteString("fn(")
	result.WriteString(strings.Join(parameters, ", "))
	result.WriteString(") {\n")
	result.WriteString(f.Body.String())
	result.WriteString("\n}")

	return result.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() string {
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string {
	return "builtin function"
}
