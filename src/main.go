package main

import "fmt"

func main() {
	code := ReadCode()

	compiler := CreateCompiler(code)
	project := compiler.Compile()
	compiler.Messanger.Report()

	innerCompiler := CreateInnerCompiler(*project)
	innerCompiler.InnerParse()
	if compiler.IsParsedSuccessfully && innerCompiler.IsParsedSuccessfully {
		fmt.Println("parse successful!")
	}
	innerCompiler.Messenger.Report()
}
