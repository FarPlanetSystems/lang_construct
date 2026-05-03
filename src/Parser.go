package main

import (
	"fmt"
	"strconv"

	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

/*
 */

type Parser struct {
	isThereReportSection bool
	isParsedSuccessfully bool
	messager             *Messenger
	currentStatement     *compiler_objects.Formula
}

func createParser(messager *Messenger) *Parser {

	res := Parser{
		messager:             messager,
		isThereReportSection: false,
		isParsedSuccessfully: true,
	}
	return &res
}

func (parser *Parser) Clear() {
	parser.isThereReportSection = false
	parser.currentStatement = nil
}

func (parser *Parser) currentToken() compiler_objects.Token {
	return parser.currentStatement.Face()
}

func (parser *Parser) parseGrammar(statement *compiler_objects.Formula) compiler_objects.SyntaxRule {
	fmt.Println("parsing grammar...")
	parser.currentStatement = statement
	grammarRuleName := parser.currentStatement.Face().Value
	syntaxRule := compiler_objects.CreateSyntaxRule(grammarRuleName, []compiler_objects.SyntaxOption{})
	parser.eatToken(compiler_objects.ID)
	parser.eatToken(compiler_objects.COLON_EQUAL)
	syntaxRule.Options = parser.parseGrammarAddRule()
	return syntaxRule
}

func (parser *Parser) parseGrammarRule() compiler_objects.SyntaxOption {
	fmt.Println("parsing grammar rule...")
	option := compiler_objects.CreateSyntaxOption()
	for parser.currentToken().TokenType != compiler_objects.COMMA && !parser.currentStatement.IsEmpty() {
		fmt.Println(parser.currentToken())
		(&option).Enqueue(parser.parseGrammarRuleWord())
		if !parser.isParsedSuccessfully {
			return compiler_objects.CreateSyntaxOption()
		}
	}
	return option
}

// reserved symobls like := must be checked
func (parser *Parser) parseGrammarRuleWord() compiler_objects.GrammarWord {
	//fmt.Println(token.token_type)
	switch parser.currentToken().TokenType {
	case compiler_objects.STRING:
		word := parser.currentToken()
		parser.eatToken(compiler_objects.STRING)
		return compiler_objects.GrammarWord{Content: compiler_objects.Token{Value: word.Value, TokenType: word.TokenType}}
	case compiler_objects.ID:
		word := parser.currentToken()
		parser.eatToken(compiler_objects.ID)
		return compiler_objects.GrammarWord{Content: compiler_objects.Token{Value: word.Value, TokenType: word.TokenType}}

	case compiler_objects.INNER_EOF:
		parser.eatToken(compiler_objects.INNER_EOF)
		return compiler_objects.GrammarWord{Content: compiler_objects.Token{Value: "", TokenType: compiler_objects.INNER_EOF}}
	default:
		parser.messager.InsertMessage("ERROR: grammar rule id or string word expected, but "+parser.currentToken().TokenType+" was found", parser.currentStatement.Line)
		parser.isParsedSuccessfully = false
		return compiler_objects.GrammarWord{}
	}

}

func (parser *Parser) parseGrammarAddRule() []compiler_objects.SyntaxOption {
	options := []compiler_objects.SyntaxOption{}
	options = append(options, parser.parseGrammarRule())

	for parser.currentToken().TokenType == compiler_objects.COMMA {
		parser.eatToken(compiler_objects.COMMA)
		options = append(options, parser.parseGrammarRule())
		if !parser.isParsedSuccessfully {
			return []compiler_objects.SyntaxOption{}
		}
	}
	return options
}

func (parser *Parser) eatToken(token_type string) {
	if parser.currentStatement.IsEmpty() {
		msg := "invalid syntaxis: " + token_type + " was expected, but end of expression was found."
		parser.messager.InsertMessage(msg, parser.currentStatement.Line)
		parser.isParsedSuccessfully = false
		return
	}
	if parser.currentStatement.Face().TokenType == token_type {
		parser.currentStatement.Dequeue()

	} else {
		msg := "invalid syntaxis: " + token_type + " was expected, but " + parser.currentStatement.Face().TokenType + " was found."
		parser.messager.InsertMessage(msg, parser.currentStatement.Line)
		parser.isParsedSuccessfully = false
	}
}

func (parser *Parser) have(statement *compiler_objects.Formula) compiler_objects.Proposition {
	parser.Clear()
	parser.currentStatement = statement
	parser.eatToken(compiler_objects.HAVE)

	proposition := compiler_objects.Proposition{Line: parser.currentStatement.Line, Formula: compiler_objects.CreateStatement(statement.Line)}
	for parser.currentToken().TokenType != compiler_objects.FROM && !parser.currentStatement.IsEmpty() && parser.isParsedSuccessfully {
		proposition.Formula.Enqueue(parser.currentToken())
		parser.eatToken(compiler_objects.FORMULA_TOKEN)
	}
	return proposition
}

func (parser *Parser) axiom(statement *compiler_objects.Formula) compiler_objects.Axiom {
	parser.Clear()
	parser.currentStatement = statement
	parser.eatToken(compiler_objects.AXIOM)

	axiom := compiler_objects.Axiom{Line: parser.currentStatement.Line, Formula: compiler_objects.CreateStatement(statement.Line)}
	for !parser.currentStatement.IsEmpty() && parser.isParsedSuccessfully {
		//fmt.Println("reading axiom token...")
		//fmt.Println(parser.currentToken())
		axiom.Formula.Enqueue(parser.currentToken())
		parser.eatToken(compiler_objects.FORMULA_TOKEN)
	}
	return axiom
}

func (parser *Parser) substitution(statement *compiler_objects.Formula) compiler_objects.Substitution {
	parser.Clear()
	parser.currentStatement = statement
	parser.eatToken(compiler_objects.SUBSTITUTION)
	var init compiler_objects.Param
	var sub compiler_objects.Param
	init, sub = parser.parseSubstitutionParams()
	substitution := compiler_objects.CreateSubstitution(init, sub)
	if parser.currentToken().TokenType == compiler_objects.COLON {
		substitution.Consditions = parser.parseSubstitutionConditions()
	}
	return substitution
}

func (parser *Parser) parseSubstitutionParams() (compiler_objects.Param, compiler_objects.Param) {
	var init compiler_objects.Param
	var sub compiler_objects.Param
	// init
	parser.eatToken(compiler_objects.BRACKETS_L)
	initParamType := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)
	initId := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)
	init = compiler_objects.Param{Id: initId, GrammarType: initParamType}
	parser.eatToken(compiler_objects.COMMA)
	//sub

	subParamType := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)
	subId := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)
	sub = compiler_objects.Param{Id: subId, GrammarType: subParamType}

	parser.eatToken(compiler_objects.BRACKETS_R)
	return init, sub
}

func (parser *Parser) parseSubstitutionConditions() []compiler_objects.Formula {
	var consditions []compiler_objects.Formula
	parser.eatToken(compiler_objects.COLON)

	for parser.currentToken().TokenType == compiler_objects.FORMULA_TOKEN && parser.isParsedSuccessfully {
		condition := compiler_objects.CreateStatement(parser.currentStatement.Line)
		fmt.Println("parsing parsing substitution conditions...")

		if parser.currentToken().Value == "," {
			consditions = append(consditions, condition)
			condition = compiler_objects.CreateStatement(parser.currentStatement.Line)
			parser.eatToken(compiler_objects.FORMULA_TOKEN)
		}
		if parser.currentToken().Value == "," {
			parser.messager.InsertMessage("unexpected symbol ,", parser.currentStatement.Line)
			parser.isParsedSuccessfully = false
			return consditions
		}

		condition.Enqueue(parser.currentToken())
		parser.eatToken(compiler_objects.FORMULA_TOKEN)

	}
	return consditions
}

func (parser *Parser) parseRuleParams() []compiler_objects.Param {
	var params []compiler_objects.Param

	parser.eatToken(compiler_objects.BRACKETS_L)

	fmt.Println("first param type: " + parser.currentToken().Value)

	paramType := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)
	paramId := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)
	params = append(params, compiler_objects.Param{Id: paramId, GrammarType: paramType})

	for parser.currentToken().TokenType != compiler_objects.BRACKETS_R && parser.isParsedSuccessfully {
		parser.eatToken(compiler_objects.COMMA)
		paramType := parser.currentToken().Value
		parser.eatToken(compiler_objects.ID)
		paramId := parser.currentToken().Value
		parser.eatToken(compiler_objects.ID)
		params = append(params, compiler_objects.Param{Id: paramId, GrammarType: paramType})
	}
	// check if each parameter has a unique id
	for i := 0; i < len(params); i++ {
		for j := i + 1; j < len(params); j++ {
			if params[i].Id == params[j].Id {
				parser.isParsedSuccessfully = false
				parser.messager.InsertMessage("several params have the same identifier.", parser.currentStatement.Line)
				return params
			}
		}
	}
	parser.eatToken(compiler_objects.BRACKETS_R)
	return params
}

func (parser *Parser) parseRulePremises() []compiler_objects.Formula {
	var premises []compiler_objects.Formula
	if parser.currentToken().TokenType != compiler_objects.ARROW {
		fmt.Println("parsing rule premises...")
		parser.eatToken(compiler_objects.COLON)

		if parser.currentToken().Value == "," {
			parser.messager.InsertMessage("unexpected symbol ,", parser.currentStatement.Line)
			parser.isParsedSuccessfully = false
			return premises
		}
		for parser.currentToken().TokenType != compiler_objects.ARROW && parser.isParsedSuccessfully {
			premise := compiler_objects.CreateStatement(parser.currentStatement.Line)
			fmt.Println("parsing rule premise...")

			for parser.currentToken().Value != "," && parser.currentToken().TokenType == compiler_objects.FORMULA_TOKEN {
				premise.Enqueue(parser.currentToken())
				parser.eatToken(compiler_objects.FORMULA_TOKEN)
			}
			premises = append(premises, premise)
			if parser.currentToken().TokenType == compiler_objects.ARROW {
				break
			}
			parser.eatToken(compiler_objects.FORMULA_TOKEN)
			if parser.currentToken().Value == "," {
				parser.isParsedSuccessfully = false
				parser.messager.InsertMessage("unexpected symbol ,", parser.currentStatement.Line)
				return premises
			}

		}
	}
	return premises
}

func (parser *Parser) parseRuleConclusion() compiler_objects.Formula {
	fmt.Println("Parsing conclsusion...")
	fmt.Println(parser.isParsedSuccessfully)
	var conclusion compiler_objects.Formula
	for parser.currentToken().TokenType == compiler_objects.FORMULA_TOKEN && parser.isParsedSuccessfully {
		fmt.Println(parser.currentToken())
		conclusion.Enqueue(parser.currentToken())
		parser.eatToken(compiler_objects.FORMULA_TOKEN)
	}
	return conclusion
}

func (parser *Parser) rule(statement *compiler_objects.Formula) compiler_objects.Rule {
	parser.Clear()
	parser.currentStatement = statement

	parser.eatToken(compiler_objects.RULE)
	// get the name
	rule_name := parser.currentToken().Value
	parser.eatToken(compiler_objects.ID)

	//get params
	params := parser.parseRuleParams()
	if !parser.isParsedSuccessfully {
		return compiler_objects.Rule{}
	}

	// get premises
	premises := parser.parseRulePremises()
	parser.eatToken(compiler_objects.ARROW)
	if !parser.isParsedSuccessfully {
		return compiler_objects.Rule{}
	}
	fmt.Println("num of premises: " + strconv.Itoa(len(premises)))

	// get conclusions
	conclusions := parser.parseRuleConclusion()
	if !parser.isParsedSuccessfully {
		return compiler_objects.Rule{}
	}

	return compiler_objects.CreateRule(rule_name, params, premises, conclusions, parser.currentStatement.Line)

}
