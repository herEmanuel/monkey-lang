package parser

import (
	"bufio"
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"os"
)

func ParseFile(filename string) *ast.Program {
	file, err := os.Open(fmt.Sprintf("./%s.mk", filename))
	if err != nil {
		fmt.Print("Error: ", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	l := lexer.New(scanner.Text(), scanner)
	p := New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println("Parsing error: " + err)
		}
	}

	return program
}
