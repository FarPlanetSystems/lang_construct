package main

import "fmt"

func main() {
	code := ReadCode()

	compiler := CreateCompiler(code)
	project := compiler.Compile()

	innerCompiler := CreateInnerCompiler(*project)
	//fmt.Println("external compilation complete")
	if compiler.IsParsedSuccessfully {
		innerCompiler.InnerParse()
	}
	//fmt.Println("internal compilation complete")
	if compiler.IsParsedSuccessfully && innerCompiler.IsParsedSuccessfully {
		fmt.Println("parse successful! Verifying project...")
		isVerified := project.Verify(*innerCompiler)
		if isVerified {
			fmt.Println("coherence verified!")
		} else {
			fmt.Println("coherence not verified :(")
		}

	}
	compiler.Messanger.Report()
	innerCompiler.Messenger.Report()
}
