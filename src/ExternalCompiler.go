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
	allAxioms            []compiler_objects.Axiom
	allRules             []compiler_objects.Rule
	allSubstitutions     []compiler_objects.Substitution
	allSubstitutes       []compiler_objects.Substitute
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
		compiler.Messanger.InsertMessage("Compilation failed: no finished statements were found", 0)
		compiler.IsParsedSuccessfully = false
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
	project := CreateProject(compiler.allSyntaxRules, compiler.allPropositions, compiler.allAxioms, compiler.allRules, compiler.allSubstitutions)
	return project
}

func (compiler *Compiler) ScanCode() *compiler_objects.StatementQueue {
	fmt.Println("scanning")
	res := compiler_objects.CreateStatementQueue()
	lexer := CreateLexer(compiler.code)
	token := lexer.getNextToken(compiler.Messanger)
	for token.TokenType != compiler_objects.REPORT_SECTION && token.TokenType != compiler_objects.EOF {
		if token.TokenType == compiler_objects.UNEXPECTED_SYMBOL {
			compiler.Messanger.InsertMessage("Unexpected symbol: "+token.Value, lexer.current_line)
			compiler.IsParsedSuccessfully = false
			return res
		}
		//fmt.Println("token: " + token.token_type)
		statement := compiler.scanStatement(lexer, token)

		if statement.IsEmpty() {
			fmt.Println("error 2")
			compiler.IsParsedSuccessfully = false
			return res
		}
		fmt.Println("enqueuing statement:")
		fmt.Println(statement)
		res.Enqueue(statement)
		token = lexer.getNextToken(compiler.Messanger)
	}
	return res
}

func (compiler *Compiler) scanStatement(lexer *Lexer, token compiler_objects.Token) compiler_objects.Formula {
	statement := compiler_objects.CreateStatement(lexer.current_line)

	for token.TokenType != compiler_objects.SEMI_COLON {
		fmt.Println("scanning token " + token.TokenType)
		fmt.Println("token2: " + token.TokenType)
		if token.TokenType == compiler_objects.UNEXPECTED_SYMBOL {
			compiler.Messanger.InsertMessage("Unexpected symbol: "+token.Value, lexer.current_line)
			compiler.IsParsedSuccessfully = false
			return compiler_objects.CreateStatement(lexer.current_line)
		}
		if token.TokenType == compiler_objects.EOF {
			compiler.Messanger.InsertMessage("Semi colon in the end of the code missing", lexer.current_line)
			compiler.IsParsedSuccessfully = false
			return compiler_objects.CreateStatement(lexer.current_line)
		}
		//fmt.Println(token)
		for token.TokenType == compiler_objects.COMMENT || token.TokenType == compiler_objects.NEW_LINE {
			token = lexer.getNextToken(compiler.Messanger)
		}
		if token.TokenType == compiler_objects.SEMI_COLON {
			token = lexer.getNextToken(compiler.Messanger)
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

func (compiler *Compiler) ParseStatement(statement compiler_objects.Formula) {
	fmt.Println("parsing...")
	stmt := &statement
	switch stmt.Face().TokenType {
	case compiler_objects.ID:
		grammar := compiler.Parser.parseGrammar(stmt)
		if !compiler.Parser.isParsedSuccessfully {
			compiler.IsParsedSuccessfully = false
			return
		}
		compiler.allSyntaxRules = append(compiler.allSyntaxRules, grammar)
		//fmt.Println(grammar)
	case compiler_objects.HAVE:
		fmt.Println("parsing have...")
		proposition := compiler.Parser.have(stmt)
		if !compiler.Parser.isParsedSuccessfully {
			compiler.IsParsedSuccessfully = false
			return
		}
		compiler.allPropositions = append(compiler.allPropositions, proposition)
	case compiler_objects.AXIOM:
		fmt.Println("parsing axiom...")
		axiom := compiler.Parser.axiom(stmt)
		if !compiler.Parser.isParsedSuccessfully {
			compiler.IsParsedSuccessfully = false
			return
		}
		compiler.allAxioms = append(compiler.allAxioms, axiom)
	case compiler_objects.RULE:
		fmt.Println("parsing rule...")
		rule := compiler.Parser.rule(stmt)
		if !compiler.Parser.isParsedSuccessfully {
			compiler.IsParsedSuccessfully = false
			return
		}
		compiler.allRules = append(compiler.allRules, rule)
	case compiler_objects.SUBSTITUTION:
		fmt.Println("parsing substitution...")
		substitution := compiler.Parser.substitution(stmt)
		if !compiler.Parser.isParsedSuccessfully {
			compiler.IsParsedSuccessfully = false
			return
		}
		compiler.allSubstitutions = append(compiler.allSubstitutions, substitution)
	default:
		compiler.Messanger.InsertMessage("Error: ID was expected.", stmt.Line)
		compiler.IsParsedSuccessfully = false
	}

}
