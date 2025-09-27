package main

import (
	"os"
)

type Parser struct {
	lexer                   *Lexer
	current_token           Token
	project                 *Project
	is_there_report_section bool
	is_parsed_successfully  bool
}

func create_Parser(lexer *Lexer, project *Project) *Parser {
	res := Parser{
		lexer:                   lexer,
		current_token:           get_next_token(lexer, project),
		project:                 project,
		is_there_report_section: false,
		is_parsed_successfully:  true,
	}
	return &res
}

func eat(parser *Parser, token_type string) {
	if parser.current_token.token_type == token_type {
		parser.current_token = get_next_token(parser.lexer, parser.project)
		if parser.current_token.token_type == UNEXPECTED_SYMBOL {
			parser.is_parsed_successfully = false
		}
	} else {
		msg := "invalid syntaxis: " + token_type + " was expected, but " + parser.current_token.token_type + " was found."
		parser.project.message(msg, parser.lexer.current_line)
		parser.is_parsed_successfully = false
		parser.current_token = get_next_token(parser.lexer, parser.project)
	}
}

// rule : RULE ID ((BRACKET_L BRACKET_R) | (BRACKET_L ID BRACKET_R) | (BRACKET_L ID (COMMA ID)* BRACKET_R)) | | (BRACKET_L $any BRACKET_R) COLON (STRING | STRING (COMMA STRING)*) ARROW (STRING | (STRING (COMMA STRING)*)) SEMI
// Example: rule sum_1 (x, y) : "x belong Natural" , "y belong Natural" -> "x + y belong Natural";
func rule(parser *Parser) {
	eat(parser, RULE)
	// get the name
	rule_name := parser.current_token.value
	eat(parser, ID)
	if parser.project.findIdInProject(rule_name) {
		parser.is_parsed_successfully = false
		parser.project.message("id "+rule_name+" already used.", parser.lexer.current_line)
		return
	}
	//get params
	var params []string
	var any_params bool = false
	eat(parser, BRACKETS_L)
	for parser.current_token.token_type != BRACKETS_R && parser.is_parsed_successfully {
		switch parser.current_token.token_type {
		case COMMA:
			eat(parser, COMMA)
			param := parser.current_token.value
			eat(parser, ID)
			params = append(params, param)
		case ID:
			param := parser.current_token.value
			eat(parser, ID)
			params = append(params, param)
		case ANY:
			any_params = true
			eat(parser, ANY)
		default:
			parser.is_parsed_successfully = false
			parser.project.message("invalid syntaxis: "+parser.current_token.token_type+" was found, but identificator or $any was expected.", parser.lexer.current_line)
			return
		}
	}
	if len(params) != 0 && any_params == true {
		parser.is_parsed_successfully = false
		parser.project.message("cannot use $any and parameters at the same time.", parser.lexer.current_line)
		return
	}
	eat(parser, BRACKETS_R)

	// check if each parameter has a unique id
	for i := 0; i < len(params); i++ {
		for j := i + 1; j < len(params); j++ {
			if params[i] == params[j] {
				parser.is_parsed_successfully = false
				parser.project.message("several params have the same identifier.", parser.lexer.current_line)
				return
			}
		}
	}
	// get premises
	var any_premises bool = false
	var premises []string
	if parser.current_token.token_type != ARROW {
		eat(parser, COLON)
		for parser.current_token.token_type != ARROW && parser.is_parsed_successfully {
			switch parser.current_token.token_type {
			case COMMA:
				eat(parser, COMMA)
				premise := parser.current_token.value
				eat(parser, STRING)
				premises = append(premises, premise)
			case STRING:
				premise := parser.current_token.value
				eat(parser, STRING)
				premises = append(premises, premise)
			case ANY:
				any_premises = true
				eat(parser, ANY)
			default:
				parser.is_parsed_successfully = false
				parser.project.message("invalid syntaxis: "+parser.current_token.token_type+" was found, but string or $any was expected.", parser.lexer.current_line)
				return
			}
		}
	}
	if len(premises) != 0 && any_premises == true {
		parser.is_parsed_successfully = false
		parser.project.message("cannot use $any and premises at the same time.", parser.lexer.current_line)
		return
	}
	eat(parser, ARROW)
	// get conclusions
	var conclusions []string
	var conclusion string = ""
	for parser.current_token.token_type != SEMI_COLON && parser.is_parsed_successfully {
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

	eat(parser, SEMI_COLON)
	create_rule(rule_name, params, premises, conclusions, parser.lexer.current_line, any_params, any_premises, parser.project)

}

// We parce a definition line
// def : DEF STRING SEMI
// Example: def "one belong Natural";
func def(parser *Parser) {
	eat(parser, DEF)
	def_line := parser.current_token.value
	eat(parser, STRING)
	create_definition(def_line, parser.project)
	eat(parser, SEMI_COLON)

}

// have : HAVE STRING FROM ID ((BRACKET_L BRACKET_R) | (BRACKET_L STRING BRACKET_R) | (BRACKET_L ID (COMMA STRING)* BRACKET_R)) (SEMI | STRING SEMI | STRING (COMMA STRING)* SEMI)
func have(parser *Parser) UnverifiedElement {

	eat(parser, HAVE)

	conclusions := []string{}
	for parser.current_token.token_type != FROM && parser.is_parsed_successfully {
		if parser.current_token.token_type == COMMA {
			eat(parser, COMMA)
			conclusion := parser.current_token.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)
		} else {
			conclusion := parser.current_token.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)
		}
	}
	// get rule's name
	eat(parser, FROM)
	rule_name := parser.current_token.value
	eat(parser, ID)
	// get params
	var params []string
	eat(parser, BRACKETS_L)
	for parser.current_token.token_type != BRACKETS_R && parser.is_parsed_successfully {

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
	for parser.current_token.token_type != SEMI_COLON && parser.is_parsed_successfully {
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
	eat(parser, SEMI_COLON)
	//create statement
	result := Proposition{
		rule_name:   rule_name,
		conclusions: conclusions,
		params:      params,
		premises:    premises,
		line:        parser.lexer.current_line - 1,
	}
	return createUnverifiedProposition(result)
}

func (parser *Parser) ifarea(project *Project) UnverifiedElement {
	eat(parser, IF)
	eat(parser, BRACKETS_L)
	param_name := parser.current_token.value
	if parser.current_token.token_type != ID {
		parser.is_parsed_successfully = false
		project.message("invalid syntax: identificator of a parameter was expected, but "+parser.current_token.token_type+" was found.", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, ID)
	eat(parser, BRACKETS_R)
	eat(parser, COLON)

	condition := parser.current_token.value
	if parser.current_token.token_type != STRING {
		parser.is_parsed_successfully = false
		project.message("invalid syntax: a conditional string was expected, but "+parser.current_token.token_type+" was found.", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, STRING)
	eat(parser, ARROW)
	eat(parser, CURL_BRACKETS_L)
	//here we parse the content of the propositional area body
	propositions := []Proposition{}
	token := parser.current_token.token_type
	for token != REPORT_SECTION && token != EOF && token != CURL_BRACKETS_R && parser.is_parsed_successfully {
		switch token {
		case HAVE:
			statement := have(parser)
			propositions = append(propositions, statement.proposition)
		case NEW_LINE:
			eat(parser, NEW_LINE)
		case COMMENT:
			eat(parser, COMMENT)
		}
		token = parser.current_token.token_type
	}
	eat(parser, CURL_BRACKETS_R)
	result := PropArea{param: param_name, condition: condition, containedPropositions: propositions, confirmedPropositions: []string{}}
	return createUnverifiedPropArea(&result)
}

func (parser *Parser) spec(project *Project) UnverifiedElement {
	eat(parser, SPEC)
	name := parser.current_token.value
	if parser.current_token.token_type != ID {
		parser.is_parsed_successfully = false
		project.message("Invalid syntax: an identificator of a specification expected, but "+parser.current_token.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, ID)
	eat(parser, BRACKETS_L)
	param := []string{parser.current_token.value}
	if parser.current_token.token_type != ID {
		parser.is_parsed_successfully = false
		project.message("Invalid syntax: an identificator of a parameter expected, but "+parser.current_token.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, ID)
	eat(parser, BRACKETS_R)
	eat(parser, COLON)
	condition := []string{parser.current_token.value}
	if parser.current_token.token_type != STRING {
		parser.is_parsed_successfully = false
		project.message("Invalid syntax: conditional string was expected, but "+parser.current_token.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, STRING)
	eat(parser, ARROW)
	conclusion := []string{parser.current_token.value}
	if parser.current_token.token_type != STRING {
		parser.is_parsed_successfully = false
		project.message("Invalid syntax: conclusion string was expected, but "+parser.current_token.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, STRING)
	eat(parser, SEMI_COLON)
	specification := Rule{
		name:              name,
		params:            param,
		conclusions:       conclusion,
		premises:          condition,
		are_any_premisses: false,
		are_any_params:    false,
		line:              parser.lexer.current_line - 1,
	}
	return createUnverifiedSpecification(specification)

}

func read_import(parser *Parser) {
	eat(parser, IMPORT)
	if parser.current_token.token_type == STRING {
		imported_file := parser.current_token.value
		curdir, err := os.Getwd()
		if err != nil {
			parser.project.message(err.Error(), -1)
		}
		imported_file_path := curdir + "\\" + imported_file
		parser.project.importedProjectsPaths = append(parser.project.importedProjectsPaths, imported_file_path)
	}
	eat(parser, STRING)
	eat(parser, SEMI_COLON)
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
			expression := have(parser)
			parser.project.statements = append(parser.project.statements, expression.proposition)
			queue := &parser.project.unverifiedExpressions
			queue.enqueue(expression)
		case NEW_LINE:
			eat(parser, NEW_LINE)
		case COMMENT:
			eat(parser, COMMENT)
		case IF:
			expression := parser.ifarea(parser.project)
			parser.project.propositionalAreas = append(parser.project.propositionalAreas, expression.propArea)
			queue := &parser.project.unverifiedExpressions
			queue.enqueue(expression)
		case SPEC:
			expression := parser.spec(parser.project)
			parser.project.specifications = append(parser.project.rules, expression.specification)
			queue := &parser.project.unverifiedExpressions
			queue.enqueue(expression)
		case IMPORT:
			read_import(parser)
		default:
			parser.is_parsed_successfully = false
			parser.project.message("unexpected expression "+parser.current_token.value, 1)
		}
		token = parser.current_token.token_type
	}
	if token == REPORT_SECTION {
		parser.is_there_report_section = true
	}
	return parser.is_parsed_successfully
}
