package main

import (
	"bufio"
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Starting...")
		repl.Start(os.Stdin, os.Stdout)
	} else {

		file, err := os.Open(fmt.Sprintf("./%s", os.Args[1]))
		if err != nil {
			fmt.Print("Error: ", err.Error())
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		l := lexer.New(scanner.Text(), scanner)
		p := parser.New(l)

		program := p.ParseProgram()
		// fmt.Println(program.String())
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Println("Parsing error: " + err)
			}
		}

		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
		}

	}
}
