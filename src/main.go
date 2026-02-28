package main

import "fmt"

func main() {
	code := ReadCode()

	compiler := CreateCompiler(code)
	project := compiler.Compile()
	innerCompiler := CreateInnerCompiler(*project)

	if compiler.IsParsedSuccessfully {
		innerCompiler.InnerParse()
	}
	if compiler.IsParsedSuccessfully && innerCompiler.IsParsedSuccessfully {
		fmt.Println("parse successful!")
	}
	compiler.Messanger.Report()
	innerCompiler.Messenger.Report()
}

// first set returns an empty array for add_num option but intended to return ["", "+"]. Figure out why
