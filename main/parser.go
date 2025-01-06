package lang_construct

import (
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

// rule : RULE ID ((BRACKET_L BRACKET_R) | (BRACKET_L ID BRACKET_R) | (BRACKET_L ID (COMMA ID)* BRACKET_R)) COLON (STRING | STRING (COMMA STRING)*) ARROW STRING SEMI
// rule sum_1 (x, y) : "x belong Natural" , "y belong Natural" -> "x + y belong Natural";
func rule(parser *Parser) {
	eat(parser, RULE)
	// get the name
	rule_name := parser.current_token.value
	eat(parser, ID)
	//get params
	var params []string
	eat(parser, BRACKETS_L)
	for parser.current_token.token_type != BRACKETS_R {
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
	

	// get premises
	var premises []string
	if parser.current_token.token_type != ARROW{
		eat(parser, COLON)
		for parser.current_token.token_type != ARROW {
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
	// get conclusion
	conclusion := parser.current_token.value
	eat(parser, STRING)
	eat(parser, SEMI)
	create_rule(rule_name, params, premises, conclusion, parser.project)
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

// statement : HAVE STRING FROM ID ((BRACKET_L BRACKET_R) | (BRACKET_L ID BRACKET_R) | (BRACKET_L ID (COMMA ID)* BRACKET_R)) (SEMI | STRING SEMI | STRING (COMMA STRING)* SEMI)
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
	for parser.current_token.token_type != BRACKETS_R {
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

	// get premises
	var premises []string
	for parser.current_token.token_type != SEMI {
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
	create_statement(rule_name, conclusion, params, premises, parser.project)
}

// we parce the whole code that must represent a formal language
// see the specification of language construct interpretator
// language: (RULE | DEF | HAVE | COMMENT | NEW_LINE)* (EOF | REPORT_SECTION)
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