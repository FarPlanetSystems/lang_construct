package main

import (
	"os"
	"strconv"
)

type Parser struct {
	lexer                   *Lexer
	current_token           Token
	project                 *LC_project
	is_there_report_section bool
	is_parsed_successfully  bool
}

func createParser(lexer *Lexer, project *LC_project) *Parser {
	res := Parser{
		lexer:                   lexer,
		current_token:           get_next_token(lexer),
		project:                 project,
		is_there_report_section: false,
		is_parsed_successfully:  true,
	}
	return &res
}

func eat(parser *Parser, tokenType string) {
	//fmt.Println("token type: " + parser.current_token.token_type + " expected token: " + tokenType + " token value: " + parser.current_token.value)
	if parser.current_token.token_type == NEW_LINE && tokenType != NEW_LINE{
		parser.current_token = get_next_token(parser.lexer)
		eat(parser, tokenType)
	}else{
		if parser.current_token.token_type == tokenType {
			parser.current_token = get_next_token(parser.lexer)
			if parser.current_token.token_type == UNEXPECTED_SYMBOL {
				parser.is_parsed_successfully = false
			}
		} else {
			msg := "invalid syntaxis: " + tokenType + " was expected, but " + parser.current_token.token_type + " was found. line " + strconv.Itoa(parser.lexer.current_line)
			message(msg, parser.project)
			parser.is_parsed_successfully = false
			parser.current_token = get_next_token(parser.lexer)
		}
	}
}

// rule : (RULE | RULE (ORDER_SIGN)*) ID ((BRACKET_L BRACKET_R) | (BRACKET_L ID BRACKET_R) | (BRACKET_L ID (COMMA ID)* BRACKET_R)) | (BRACKET_L ANY BRACKET_R) COLON ((STRING | rule) | (STRING | rule (COMMA (STRING | rule))*) ARROW (STRING | (STRING (COMMA STRING)*)) SEMI
// Example: rule sum_1 (x, y) : "x belong Natural" , "y belong Natural" -> "x + y belong Natural";
func rule(parser *Parser) Rule {
	eat(parser, RULE)
	// get the name
	ruleName := parser.current_token.value

	
	eat(parser, ID)
	if findIdInProject(ruleName, *parser.project){
		parser.is_parsed_successfully = false
		message("id " + ruleName +" already used. Line"+ strconv.Itoa(parser.lexer.current_line), parser.project)
		return Rule{}
	}
	//get params
	var params []string
	var anyParams bool = false

	eat(parser, BRACKETS_L)
	for parser.current_token.token_type != BRACKETS_R && parser.is_parsed_successfully{
		switch parser.current_token.token_type {
		case COMMA:
			eat(parser, COMMA)
			param := parser.current_token.value
			eat(parser, ID)
			params = append(params, param)
			break
		case ID:
			param := parser.current_token.value
			eat(parser, ID)
			params = append(params, param)
		case ANY:
			anyParams = true
			eat(parser, ANY)
		default:
			parser.is_parsed_successfully = false
			message("invalid syntaxis: " + parser.current_token.token_type + " was found, but identificator or $any was expected. Line " + strconv.Itoa(parser.lexer.current_line), parser.project)
			return Rule{}
		}
	}
	if len(params) != 0 && anyParams == true{
		parser.is_parsed_successfully = false
		message("cannot use $any and parameters at the same time. Line " + strconv.Itoa(parser.lexer.current_line), parser.project)
		return Rule{}
	}
	eat(parser, BRACKETS_R)

	// check if each parameter has a unique id
	for i := 0; i< len(params); i++{
		for j:= i + 1; j < len(params); j++{
			if params[i] == params[j]{
				parser.is_parsed_successfully = false
				message("several params have the same identifier. Line " + strconv.Itoa(parser.lexer.current_line), parser.project)
				return Rule{}
			}
		} 
	}
	// get premises
	//fmt.Println("parsing premises")
	var anyPremises bool = false
	var premises []Argument
	if parser.current_token.token_type != ARROW{
		eat(parser, COLON)
		for parser.current_token.token_type != ARROW && parser.is_parsed_successfully{

			switch parser.current_token.token_type {
				case COMMA:
					eat(parser, COMMA)
					premises = append(premises, readArgument(parser))
				case ANY:
					anyPremises = true
					eat(parser, ANY)
					break
				default:
					premises = append(premises, readArgument(parser))
			}
		}
	}
	if len(premises) != 0 && anyPremises == true{
		parser.is_parsed_successfully = false
		message("cannot use $any and premises at the same time. Line " + strconv.Itoa(parser.lexer.current_line), parser.project)
		return Rule{}
	}
	eat(parser, ARROW)
	// get conclusions
	//fmt.Println("parsing conclusion")
	var conclusions []Argument

	for parser.current_token.token_type != SEMI && parser.current_token.token_type != CURL_BRACKETS_R && parser.is_parsed_successfully{
		switch parser.current_token.token_type {
				case COMMA:
					eat(parser, COMMA)
					conclusions = append(premises, readArgument(parser))
				case ANY:
					anyPremises = true
					eat(parser, ANY)
				default:
					conclusions = append(conclusions, readArgument(parser))
			}
	}
	res := createRule(ruleName, params, premises, conclusions, parser.lexer.current_line, anyParams, anyPremises)
	return res
	
}

func readArgument(parser *Parser) Argument {
	switch parser.current_token.token_type {
				
				case STRING:
					premiseProposition := parser.current_token.value
					argument := createArgument(PROPOSITIONAL_ARGUMENT_TYPE, premiseProposition, Rule{})
					eat(parser, STRING)
					return argument
				case CURL_BRACKETS_L:
					eat(parser, CURL_BRACKETS_L)
					premiseRule := rule(parser)
					argument := createArgument(RULE_ARGUMENT_TYPE, "", premiseRule)
					eat(parser, CURL_BRACKETS_R)
					return argument
				default:
					parser.is_parsed_successfully = false
					message("invalid syntaxis: " + parser.current_token.token_type + " was found, but string or $any was expected. Line " + strconv.Itoa(parser.lexer.current_line), parser.project)
					return Argument{}
			}
}

// We parce a definition line
// def : DEF STRING SEMI
// Example: def "one belong Natural";
func def(parser *Parser) string {
	eat(parser, DEF)
	defLine := parser.current_token.value
	eat(parser, STRING)
	return defLine
}

// statement : HAVE STRING FROM ID ((BRACKET_L BRACKET_R) | (BRACKET_L STRING BRACKET_R) | (BRACKET_L ID (COMMA STRING)* BRACKET_R)) (SEMI | STRING SEMI | STRING (COMMA STRING)* SEMI)
func statement(parser *Parser) {

	eat(parser, HAVE)
	// get conclusion
	conclusion := readArgument(parser)
	// get rule's name
	eat(parser, FROM)
	ruleName := parser.current_token.value
	eat(parser, ID)
	// get params
	var params []string
	eat(parser, BRACKETS_L)
	for parser.current_token.token_type != BRACKETS_R && parser.is_parsed_successfully{

		if parser.current_token.token_type == COMMA {
			eat(parser, COMMA)
			param := parser.current_token.value
			eat(parser, STRING)
			params = append(params, param)
		} else {
			param := parser.current_token.value
			eat(parser, STRING)
			params = append(params, param)
		}
	}
	eat(parser, BRACKETS_R)

	// get premises
	var premises []Argument
	for parser.current_token.token_type != SEMI && parser.is_parsed_successfully{
		if parser.current_token.token_type == COMMA {
			eat(parser, COMMA)
			premises = append(premises, readArgument(parser))
		} else {
			premises = append(premises, readArgument(parser))
		}
	}
	eat(parser, SEMI)
	//create statement
	createStatement(ruleName, conclusion, params, premises, parser.lexer.current_line, parser.project)
}

func readImport(parser *Parser){
	eat(parser, IMPORT)
	if parser.current_token.token_type == STRING{
	importedFile := parser.current_token.value
	curdir, err := os.Getwd()
	if err != nil{
		message(err.Error(), parser.project)
	}
	importedFilePath := curdir + "\\" + importedFile
	parser.project.imported_projects_file_paths = append(parser.project.imported_projects_file_paths, importedFilePath)
}
	eat(parser, STRING)
	eat(parser, SEMI)
}

// we parce the whole code that must represent a formal language
// see the specification of language construct interpretator
// language: (RULE | DEF | HAVE | IMPORT | COMMENT | NEW_LINE)* (EOF | REPORT_SECTION)
// returns true if code is succefully parced
func Language(parser *Parser) bool {
	
	token := parser.current_token.token_type
	for token != REPORT_SECTION && token != EOF && parser.is_parsed_successfully {
		//fmt.Println("given: \n" + token)

		switch token {
		case RULE:
			rule := rule(parser)
			appendRule(rule, parser.project)
			eat(parser, SEMI)
		case DEF:
			defLine:= def(parser)
			appendDefinition(defLine, parser.project)
			eat(parser, SEMI)
		case HAVE:
			statement(parser)
		case NEW_LINE:
			eat(parser, NEW_LINE)
		case COMMENT:
			eat(parser, COMMENT)
		case IMPORT:
			readImport(parser)
		default:
			parser.is_parsed_successfully = false
			message("unexpected expression "+parser.current_token.value+" on the line 1", parser.project)
		}
		token = parser.current_token.token_type
	}
	if token == REPORT_SECTION {
		parser.is_there_report_section = true
	}
	return parser.is_parsed_successfully
}