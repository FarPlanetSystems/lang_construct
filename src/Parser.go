package main

import (
	"os"
)

type Parser struct {
	lexer                *Lexer
	currentToken         Token
	isThereReportSection bool
	isParsedSuccessfully bool
	messager             *Messanger
}

func createParser(lexer *Lexer, messager *Messanger) *Parser {
	res := Parser{
		lexer:                lexer,
		currentToken:         getNextToken(lexer, messager),
		isThereReportSection: false,
		isParsedSuccessfully: true,
	}
	return &res
}

func eat(parser *Parser, token_type string) {
	if parser.currentToken.token_type == token_type {
		parser.currentToken = getNextToken(parser.lexer, parser.messager)
		if parser.currentToken.token_type == UNEXPECTED_SYMBOL {
			parser.isParsedSuccessfully = false
		}
	} else {
		msg := "invalid syntaxis: " + token_type + " was expected, but " + parser.currentToken.token_type + " was found."
		parser.messager.message(msg, parser.lexer.current_line)
		parser.isParsedSuccessfully = false
		parser.currentToken = getNextToken(parser.lexer, parser.messager)
	}
}

// rule : RULE ID ((BRACKET_L BRACKET_R) | (BRACKET_L ID BRACKET_R) | (BRACKET_L ID (COMMA ID)* BRACKET_R)) | | (BRACKET_L $any BRACKET_R) COLON (STRING | STRING (COMMA STRING)*) ARROW (STRING | (STRING (COMMA STRING)*)) SEMI
// Example: rule sum_1 (x, y) : "x belong Natural" , "y belong Natural" -> "x + y belong Natural";
func (parser *Parser) rule(project *Project) {
	eat(parser, RULE)
	// get the name
	rule_name := parser.currentToken.value
	eat(parser, ID)
	if project.findIdInProject(rule_name) {
		parser.isParsedSuccessfully = false
		parser.messager.message("id "+rule_name+" already used.", parser.lexer.current_line)
		return
	}
	//get params
	var params []string
	var any_params bool = false
	eat(parser, BRACKETS_L)
	for parser.currentToken.token_type != BRACKETS_R && parser.isParsedSuccessfully {
		switch parser.currentToken.token_type {
		case COMMA:
			eat(parser, COMMA)
			param := parser.currentToken.value
			eat(parser, ID)
			params = append(params, param)
		case ID:
			param := parser.currentToken.value
			eat(parser, ID)
			params = append(params, param)
		case ANY:
			any_params = true
			eat(parser, ANY)
		default:
			parser.isParsedSuccessfully = false
			parser.messager.message("invalid syntaxis: "+parser.currentToken.token_type+" was found, but identificator or $any was expected.", parser.lexer.current_line)
			return
		}
	}
	if len(params) != 0 && any_params == true {
		parser.isParsedSuccessfully = false
		parser.messager.message("cannot use $any and parameters at the same time.", parser.lexer.current_line)
		return
	}
	eat(parser, BRACKETS_R)

	// check if each parameter has a unique id
	for i := 0; i < len(params); i++ {
		for j := i + 1; j < len(params); j++ {
			if params[i] == params[j] {
				parser.isParsedSuccessfully = false
				parser.messager.message("several params have the same identifier.", parser.lexer.current_line)
				return
			}
		}
	}
	// get premises
	var any_premises bool = false
	var premises []string
	if parser.currentToken.token_type != ARROW {
		eat(parser, COLON)
		for parser.currentToken.token_type != ARROW && parser.isParsedSuccessfully {
			switch parser.currentToken.token_type {
			case COMMA:
				eat(parser, COMMA)
				premise := parser.currentToken.value
				eat(parser, STRING)
				premises = append(premises, premise)
			case STRING:
				premise := parser.currentToken.value
				eat(parser, STRING)
				premises = append(premises, premise)
			case ANY:
				any_premises = true
				eat(parser, ANY)
			default:
				parser.isParsedSuccessfully = false
				parser.messager.message("invalid syntaxis: "+parser.currentToken.token_type+" was found, but string or $any was expected.", parser.lexer.current_line)
				return
			}
		}
	}
	if len(premises) != 0 && any_premises == true {
		parser.isParsedSuccessfully = false
		parser.messager.message("cannot use $any and premises at the same time.", parser.lexer.current_line)
		return
	}
	eat(parser, ARROW)
	// get conclusions
	var conclusions []string
	var conclusion string = ""
	for parser.currentToken.token_type != SEMI_COLON && parser.isParsedSuccessfully {
		if parser.currentToken.token_type == COMMA {
			eat(parser, COMMA)
			conclusion = parser.currentToken.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)

		} else {
			conclusion = parser.currentToken.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)
		}
	}

	eat(parser, SEMI_COLON)
	create_rule(rule_name, params, premises, conclusions, parser.lexer.current_line, any_params, any_premises, project)

}

// We parce a definition line
// def : DEF STRING SEMI
// Example: def "one belong Natural";
func (parser *Parser) def(project *Project) {
	eat(parser, DEF)
	def_line := parser.currentToken.value
	eat(parser, STRING)
	createDefinition(def_line, project)
	eat(parser, SEMI_COLON)

}

// have : HAVE STRING FROM ID ((BRACKET_L BRACKET_R) | (BRACKET_L STRING BRACKET_R) | (BRACKET_L ID (COMMA STRING)* BRACKET_R)) (SEMI | STRING SEMI | STRING (COMMA STRING)* SEMI)
func have(parser *Parser) UnverifiedElement {

	eat(parser, HAVE)

	conclusions := []string{}
	for parser.currentToken.token_type != FROM && parser.isParsedSuccessfully {
		if parser.currentToken.token_type == COMMA {
			eat(parser, COMMA)
			conclusion := parser.currentToken.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)
		} else {
			conclusion := parser.currentToken.value
			eat(parser, STRING)
			conclusions = append(conclusions, conclusion)
		}
	}
	// get rule's name
	eat(parser, FROM)
	rule_name := parser.currentToken.value
	eat(parser, ID)
	// get params
	var params []string
	eat(parser, BRACKETS_L)
	for parser.currentToken.token_type != BRACKETS_R && parser.isParsedSuccessfully {

		if parser.currentToken.token_type == COMMA {
			eat(parser, COMMA)
			param := parser.currentToken.value
			eat(parser, STRING)
			params = append(params, param)
		} else {
			param := parser.currentToken.value
			eat(parser, STRING)
			params = append(params, param)
		}
	}
	eat(parser, BRACKETS_R)

	// get premises
	var premises []string
	for parser.currentToken.token_type != SEMI_COLON && parser.isParsedSuccessfully {
		if parser.currentToken.token_type == COMMA {
			eat(parser, COMMA)
			premise := parser.currentToken.value
			eat(parser, STRING)
			premises = append(premises, premise)
		} else {
			premise := parser.currentToken.value
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
	param_name := parser.currentToken.value
	if parser.currentToken.token_type != ID {
		parser.isParsedSuccessfully = false
		parser.messager.message("invalid syntax: identificator of a parameter was expected, but "+parser.currentToken.token_type+" was found.", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, ID)
	eat(parser, BRACKETS_R)
	eat(parser, COLON)

	condition := parser.currentToken.value
	if parser.currentToken.token_type != STRING {
		parser.isParsedSuccessfully = false
		parser.messager.message("invalid syntax: a conditional string was expected, but "+parser.currentToken.token_type+" was found.", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, STRING)
	eat(parser, ARROW)
	eat(parser, CURL_BRACKETS_L)
	//here we parse the content of the propositional area body
	propositions := []Proposition{}
	token := parser.currentToken.token_type
	for token != REPORT_SECTION && token != EOF && token != CURL_BRACKETS_R && parser.isParsedSuccessfully {
		switch token {
		case HAVE:
			statement := have(parser)
			propositions = append(propositions, statement.proposition)
		case NEW_LINE:
			eat(parser, NEW_LINE)
		case COMMENT:
			eat(parser, COMMENT)
		}
		token = parser.currentToken.token_type
	}
	eat(parser, CURL_BRACKETS_R)
	result := PropArea{param: param_name, condition: condition, containedPropositions: propositions, confirmedPropositions: []string{}}
	return createUnverifiedPropArea(&result)
}

func (parser *Parser) spec(project *Project) UnverifiedElement {
	eat(parser, SPEC)
	name := parser.currentToken.value
	if parser.currentToken.token_type != ID {
		parser.isParsedSuccessfully = false
		parser.messager.message("Invalid syntax: an identificator of a specification expected, but "+parser.currentToken.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, ID)
	eat(parser, BRACKETS_L)
	param := []string{parser.currentToken.value}
	if parser.currentToken.token_type != ID {
		parser.isParsedSuccessfully = false
		parser.messager.message("Invalid syntax: an identificator of a parameter expected, but "+parser.currentToken.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, ID)
	eat(parser, BRACKETS_R)
	eat(parser, COLON)
	condition := []string{parser.currentToken.value}
	if parser.currentToken.token_type != STRING {
		parser.isParsedSuccessfully = false
		parser.messager.message("Invalid syntax: conditional string was expected, but "+parser.currentToken.token_type+" was found. ", parser.lexer.current_line)
		return UnverifiedElement{}
	}
	eat(parser, STRING)
	eat(parser, ARROW)
	conclusion := []string{parser.currentToken.value}
	if parser.currentToken.token_type != STRING {
		parser.isParsedSuccessfully = false
		parser.messager.message("Invalid syntax: conclusion string was expected, but "+parser.currentToken.token_type+" was found. ", parser.lexer.current_line)
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

func (parser *Parser) readImport(project *Project) {
	eat(parser, IMPORT)
	if parser.currentToken.token_type == STRING {
		imported_file := parser.currentToken.value
		curdir, err := os.Getwd()
		if err != nil {
			parser.messager.message(err.Error(), -1)
		}
		imported_file_path := curdir + "\\" + imported_file
		project.importedProjectsPaths = append(project.importedProjectsPaths, imported_file_path)
	}
	eat(parser, STRING)
	eat(parser, SEMI_COLON)
}

// we parce the whole code that must represent a formal language
// see the specification of language construct interpretator
// language: (RULE | DEF | HAVE | IMPORT | COMMENT | NEW_LINE)* (EOF | REPORT_SECTION)
// returns true if code is succefully parced
func (parser *Parser) Language(project *Project) bool {

	token := parser.currentToken.token_type
	for token != REPORT_SECTION && token != EOF && parser.isParsedSuccessfully {
		//fmt.Println("given: \n" + token)

		switch token {
		case RULE:
			parser.rule(project)
		case DEF:
			parser.def(project)
		case HAVE:
			expression := have(parser)
			project.statements = append(project.statements, expression.proposition)
			queue := &project.unverifiedExpressions
			queue.enqueue(expression)
		case NEW_LINE:
			eat(parser, NEW_LINE)
		case COMMENT:
			eat(parser, COMMENT)
		case IF:
			expression := parser.ifarea(project)
			project.propositionalAreas = append(project.propositionalAreas, expression.propArea)
			queue := &project.unverifiedExpressions
			queue.enqueue(expression)
		case SPEC:
			expression := parser.spec(project)
			project.specifications = append(project.specifications, expression.specification)
			queue := &project.unverifiedExpressions
			queue.enqueue(expression)
		case IMPORT:
			parser.readImport(project)
		default:
			parser.isParsedSuccessfully = false
			parser.messager.message("unexpected expression "+parser.currentToken.value, 1)
		}
		token = parser.currentToken.token_type
	}
	if token == REPORT_SECTION {
		parser.isThereReportSection = true
	}
	return parser.isParsedSuccessfully
}
