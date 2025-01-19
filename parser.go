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

func create_Parser(lexer *Lexer, project *LC_project) *Parser {
	res := Parser{
		lexer:                   lexer,
		current_token:           get_next_token(lexer),
		project:                 project,
		is_there_report_section: false,
		is_parsed_successfully:  true,
	}
	return &res
}

func eat(parser *Parser, token_type string) {
	//fmt.Println("given: \n" + parser.current_token.token_type)
	//fmt.Println("expected: \n" + token_type)
	if parser.current_token.token_type == token_type {
		parser.current_token = get_next_token(parser.lexer)
		if parser.current_token.token_type == UNEXPECTED_SYMBOL {
			parser.is_parsed_successfully = false
		}
	} else {
		msg := "invalid syntaxis: " + token_type + " was expected, but " + parser.current_token.token_type + " was found. line " + strconv.Itoa(parser.lexer.current_line)
		message(msg, parser.project)
		parser.is_parsed_successfully = false
		parser.current_token = get_next_token(parser.lexer)
	}
}

// rule : RULE ID ((BRACKET_L BRACKET_R) | (BRACKET_L ID BRACKET_R) | (BRACKET_L ID (COMMA ID)* BRACKET_R)) COLON (STRING | STRING (COMMA STRING)*) ARROW (STRING | (STRING (COMMA STRING)*)) SEMI
// rule sum_1 (x, y) : "x belong Natural" , "y belong Natural" -> "x + y belong Natural";
func rule(parser *Parser) {
	eat(parser, RULE)
	// get the name
	rule_name := parser.current_token.value
	eat(parser, ID)
	if find_id_in_project(rule_name, *parser.project){
		parser.is_parsed_successfully = false
		message("id " + rule_name +" already used. Line"+ strconv.Itoa(parser.lexer.current_line), parser.project)
		return
	}
	//get params
	var params []string
	eat(parser, BRACKETS_L)
	for parser.current_token.token_type != BRACKETS_R && parser.is_parsed_successfully{
		if parser.current_token.token_type == COMMA {
			eat(parser, COMMA)
			param := parser.current_token.value
			eat(parser, ID)
			params = append(params, param)
		} else {
			param := parser.current_token.value
			eat(parser, ID)
			params = append(params, param)
		}
	}
	eat(parser, BRACKETS_R)

	// check if each parameter has a unique id
	for i := 0; i< len(params); i++{
		for j:= i + 1; j < len(params); j++{
			if params[i] == params[j]{
				parser.is_parsed_successfully = false
				message("several params have the same identifier. Line " + strconv.Itoa(parser.lexer.current_line), parser.project)
				return
			}
		} 
	}
	// get premises
	var premises []string
	if parser.current_token.token_type != ARROW{
		eat(parser, COLON)
		for parser.current_token.token_type != ARROW && parser.is_parsed_successfully{
			if parser.current_token.token_type == COMMA {
				eat(parser, COMMA)
				premise := parser.current_token.value
				eat(parser, STRING)
				premises = append(premises, premise)
			} else {
				premise := parser.current_token.value
				eat(parser, STRING)
				premises = append(premises, premise)
			}
		}
	}
	eat(parser, ARROW)
	// get conclusions
	var conclusions []string
	var conclusion string = ""
	for parser.current_token.token_type != SEMI && parser.is_parsed_successfully{
		if parser.current_token.token_type == COMMA {
			eat(parser, COMMA)
			conclusion = parser.current_token.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)

		} else {
			conclusion = parser.current_token.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)
		}
	}

	eat(parser, SEMI)
	create_rule(rule_name, params, premises, conclusions, parser.lexer.current_line, parser.project)
	
}

// We parce a definition line
// def : DEF STRING SEMICOLON
// Example: def "one belong Natural";
func def(parser *Parser) {
	eat(parser, DEF)
	def_line := parser.current_token.value
	eat(parser, STRING)
	create_definition(def_line, parser.project)
	eat(parser, SEMI)

}

// statement : HAVE STRING FROM ID ((BRACKET_L BRACKET_R) | (BRACKET_L STRING BRACKET_R) | (BRACKET_L ID (COMMA STRING)* BRACKET_R)) (SEMI | STRING SEMI | STRING (COMMA STRING)* SEMI)
func statement(parser *Parser) {

	eat(parser, HAVE)
	conclusion := parser.current_token.value
	//get statement
	eat(parser, STRING)
	// get rule's name
	eat(parser, FROM)
	rule_name := parser.current_token.value
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
	var premises []string
	for parser.current_token.token_type != SEMI && parser.is_parsed_successfully{
		if parser.current_token.token_type == COMMA {
			eat(parser, COMMA)
			premise := parser.current_token.value
			eat(parser, STRING)
			premises = append(premises, premise)
		} else {
			premise := parser.current_token.value
			eat(parser, STRING)
			premises = append(premises, premise)
		}
	}
	eat(parser, SEMI)
	//create statement
	create_statement(rule_name, conclusion, params, premises, parser.lexer.current_line, parser.project)
}

func read_import(parser *Parser){
	eat(parser, IMPORT)
	if parser.current_token.token_type == STRING{
	imported_file := parser.current_token.value
	curdir, err := os.Getwd()
	if err != nil{
		message(err.Error(), parser.project)
	}
	imported_file_path := curdir + "\\" + imported_file
	parser.project.imported_projects_file_paths = append(parser.project.imported_projects_file_paths, imported_file_path)
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
			rule(parser)
		case DEF:
			def(parser)
		case HAVE:
			statement(parser)
		case NEW_LINE:
			eat(parser, NEW_LINE)
		case COMMENT:
			eat(parser, COMMENT)
		case IMPORT:
			read_import(parser)
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