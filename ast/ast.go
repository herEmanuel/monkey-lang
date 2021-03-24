package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var result bytes.Buffer

	for _, s := range p.Statements {
		result.WriteString(s.String())
	}

	return result.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	var result bytes.Buffer

	result.WriteString(ls.TokenLiteral() + " ")
	result.WriteString(ls.Name.String())
	result.WriteString(" = ")
	result.WriteString(ls.Value.String())
	result.WriteString(";")

	return result.String()

}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	var result bytes.Buffer

	result.WriteString(rs.TokenLiteral() + " ")
	result.WriteString(rs.Value.String())
	result.WriteString(";")

	return result.String()
}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Block     BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}
func (ws *WhileStatement) String() string {
	var result bytes.Buffer

	result.WriteString("while (")
	result.WriteString(ws.Condition.String())
	result.WriteString(") {")
	result.WriteString(ws.Block.String())
	result.WriteString("}")

	return result.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

type Array struct {
	Token    token.Token
	Elements []Expression
}

func (a *Array) expressionNode() {}
func (a *Array) TokenLiteral() string {
	return a.Token.Literal
}
func (a *Array) String() string {
	var result bytes.Buffer

	var elements []string
	for _, a := range a.Elements {
		elements = append(elements, a.String())
	}

	result.WriteString("[")
	result.WriteString(strings.Join(elements, ", "))
	result.WriteString("]")

	return result.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var result bytes.Buffer

	result.WriteString(pe.Operator + " ")
	result.WriteString(pe.Right.String())

	return result.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var result bytes.Buffer

	result.WriteString("( ")
	result.WriteString(ie.Left.String())
	result.WriteString(" " + ie.Operator + " ")
	result.WriteString(ie.Right.String())
	result.WriteString(" )")

	return result.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var result bytes.Buffer

	for _, s := range bs.Statements {
		result.WriteString(s.String() + " ")
	}

	return result.String()
}

type IfExpression struct {
	Token      token.Token
	Condition  Expression
	TrueBlock  BlockStatement
	FalseBlock BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var result bytes.Buffer
	result.WriteString("if")
	result.WriteString(ie.Condition.String())
	result.WriteString(" ")
	result.WriteString(ie.TrueBlock.String())

	if ie.FalseBlock.Statements != nil {
		result.WriteString("else ")
		result.WriteString(ie.FalseBlock.String())
	}

	return result.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Name       Identifier
	Parameters []Identifier
	Block      BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var result bytes.Buffer

	result.WriteString(fl.TokenLiteral())
	if fl.Name.Value != "" {
		result.WriteString(fl.Name.Value)
	}
	result.WriteString(" (")
	for _, s := range fl.Parameters {
		result.WriteString(s.String() + ", ")
	}
	result.WriteString(") ")

	result.WriteString(fl.Block.String())

	return result.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
	Lib       string
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var result bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	result.WriteString(ce.Function.String())
	result.WriteString("(")
	result.WriteString(strings.Join(args, ", "))
	result.WriteString(")")

	return result.String()
}

type ArrayAccessExpression struct {
	Token    token.Token
	Array    Expression
	Position Expression
}

func (aae *ArrayAccessExpression) expressionNode() {}
func (aae *ArrayAccessExpression) TokenLiteral() string {
	return aae.Token.Literal
}
func (aae *ArrayAccessExpression) String() string {
	var result bytes.Buffer

	result.WriteString(aae.Array.String())
	result.WriteString("[")
	result.WriteString(aae.Position.String())
	result.WriteString("]")

	return result.String()
}

type ReassignmentStatement struct {
	Token    token.Token
	Variable Identifier
	NewValue Expression
}

func (rs *ReassignmentStatement) statementNode() {}
func (rs *ReassignmentStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReassignmentStatement) String() string {
	var result bytes.Buffer

	result.WriteString(rs.Variable.String())
	result.WriteString(" = ")
	result.WriteString(rs.NewValue.String())

	return result.String()
}

type UseStatement struct {
	Token    token.Token
	Filename string
}

func (us *UseStatement) statementNode() {}
func (us *UseStatement) TokenLiteral() string {
	return us.Token.Literal
}
func (us *UseStatement) String() string {
	return "use " + us.Filename
}

type ExternalReferenceExpression struct {
	Token    token.Token
	Module   string
	Referece Expression
}

func (ere *ExternalReferenceExpression) expressionNode() {}
func (ere *ExternalReferenceExpression) TokenLiteral() string {
	return ere.Token.Literal
}
func (ere *ExternalReferenceExpression) String() string {
	var result bytes.Buffer

	result.WriteString(ere.Module)
	result.WriteString(".")
	result.WriteString(ere.Referece.String())

	return result.String()
}
