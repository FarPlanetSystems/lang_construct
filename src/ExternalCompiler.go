package main

import (
	"fmt"

	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

type Compiler struct {
	code                 string
	Statements           compiler_objects.StatementQueue
	Parser               Parser
	Messanger            *Messenger
	allPropositions      []compiler_objects.Proposition
	allSyntaxRules       []compiler_objects.SyntaxRule
	IsParsedSuccessfully bool
}

func CreateCompiler(code string) *Compiler {
	messager := Messenger{}
	res := Compiler{
		code:                 code,
		Messanger:            &messager,
		Parser:               *createParser(&messager),
		IsParsedSuccessfully: true,
	}
	return &res
}

func (compiler *Compiler) Compile() *Project {
	statements := compiler.ScanCode()
	//fmt.Println("statements")
	//fmt.Println(statements)
	if (*statements).IsEmpty {
		compiler.Messanger.InsertMessage("Compilation failed", 0)
		compiler.Messanger.Report()
		return nil
	}

	//interprete the syntax rule. Convert them to a grammar tree. There must be the "statement" and "term" grammar rules
	for !statements.IsEmpty {
		statement := statements.Face()
		if !statement.IsEmpty() {
			compiler.ParseStatement(*statement)
			compiler.Parser.Clear()
		}
		statements.Dequeue()
		//fmt.Println(statement)
	}
	if compiler.IsParsedSuccessfully {
		fmt.Println("parsed successfully")
	}
	project := createProject(compiler.allSyntaxRules, compiler.allPropositions)
	return project
}

func (compiler *Compiler) ScanCode() *compiler_objects.StatementQueue {
	fmt.Println("scanning")
	res := compiler_objects.CreateStatementQueue()
	lexer := CreateLexer(compiler.code)
	token := lexer.getNextToken(compiler.Messanger)
	for token.TokenType != compiler_objects.REPORT_SECTION && token.TokenType != compiler_objects.EOF {
		if token.TokenType == compiler_objects.UNEXPECTED_SYMBOL {
			fmt.Println("error 1")
			compiler.Messanger.InsertMessage("Unexpected symbol: "+token.Value, lexer.current_line)
			return res
		}
		//fmt.Println("token: " + token.token_type)
		statement := compiler.scanStatement(lexer, token)

		if statement.IsEmpty() {
			fmt.Println("error 2")
			return res
		}
		fmt.Println("enqueuing statement:")
		fmt.Println(statement)
		res.Enqueue(statement)
		token = lexer.getNextToken(compiler.Messanger)
	}
	return res
}

func (compiler *Compiler) scanStatement(lexer *Lexer, token compiler_objects.Token) compiler_objects.Statement {
	statement := compiler_objects.CreateStatement()

	for token.TokenType != compiler_objects.SEMI_COLON {
		//fmt.Println("scanning token " + token.token_type)
		//fmt.Println("token2: " + token.token_type)
		if token.TokenType == compiler_objects.UNEXPECTED_SYMBOL {
			compiler.Messanger.InsertMessage("Unexpected symbol: "+token.Value, lexer.current_line)

			return compiler_objects.CreateStatement()
		}
		if token.TokenType == compiler_objects.EOF {
			compiler.Messanger.InsertMessage("Semi colon in the end of the code missing", lexer.current_line)
			return compiler_objects.CreateStatement()
		}
		//fmt.Println(token)
		for token.TokenType == compiler_objects.COMMENT || token.TokenType == compiler_objects.NEW_LINE {
			token = lexer.getNextToken(compiler.Messanger)
		}
		if token.TokenType == compiler_objects.SEMI_COLON {
			break
		}
		statement.Enqueue(token)
		token = lexer.getNextToken(compiler.Messanger)
	}
	statement.Line = lexer.current_line
	//fmt.Println(statement.headToken)
	//fmt.Println(statement.isEmpty)
	return statement
}

func (compiler *Compiler) ParseStatement(statement compiler_objects.Statement) {
	fmt.Println("parsing...")
	stmt := &statement
	switch stmt.Face().TokenType {
	case compiler_objects.ID:
		grammar := compiler.Parser.parseGrammar(stmt)
		compiler.allSyntaxRules = append(compiler.allSyntaxRules, grammar)
		fmt.Println(grammar)
	case compiler_objects.HAVE:
		fmt.Println("parsing have...")
		proposition := compiler.Parser.have(stmt)
		compiler.allPropositions = append(compiler.allPropositions, proposition)
	default:
		compiler.Messanger.InsertMessage("Error: ID was expected.", stmt.Line)
	}

}
