package main

import (
	"strconv"
	"unicode"
)

type RESERVED_KEYWORD struct
{
	keyword string
	token Token
}

var RESERVED_KEYWORDS = [4] RESERVED_KEYWORD{
	{
		keyword: "have",
		token: Token{
			token_type: HAVE,
			value: "have",
		},
	},
	{
		keyword: "def",
		token: Token{
			token_type: DEF,
			value: "def",
		},
	},
	{
		keyword: "rule",
		token: Token{
			token_type: RULE,
			value: "rule",
		},
	},
	{
		keyword: "from",
		token: Token{
			token_type: FROM,
			value: "from",
		},
	},
}
// gets a string which is supposed to be a reserved key word
// if it is, returns the index of it in RESERVED_WORDS array
// otherwise returns -1
func find_reserved_word(word string)int{
	for i := 0; i < len(RESERVED_KEYWORDS); i++{
		if RESERVED_KEYWORDS[i].keyword == word{
			return i
		}
	}
	return -1
}

type Lexer struct {
	text         string
	pos          int
	current_char byte
	current_line int
	project *LC_project
}

func create_Lexer(text string, project *LC_project) *Lexer {
	res := Lexer{
		text:         text,
		pos:          0,
		current_char: byte(text[0]),
		current_line: 1,
		project: project,
	}
	return &res
}

func advance(lexer_ptr *Lexer) {
	lexer_ptr.pos += 1
	if lexer_ptr.pos > len(lexer_ptr.text)-1 {
		lexer_ptr.current_char = 0
	} else {
		lexer_ptr.current_char = byte(lexer_ptr.text[lexer_ptr.pos])
	}
}

func peek(lexer *Lexer) byte {
	peek_pos := lexer.pos + 1
	if peek_pos > len(lexer.text)-1 {
		return 0
	} else {
		return byte(lexer.text[peek_pos])
	}
}

func read_string(lexer *Lexer) string {

	res := ""
	advance(lexer)
	for lexer.current_char != 0 && lexer.current_char != byte('"') {
		res += string(lexer.current_char)
		advance(lexer)
	}

	if lexer.current_char == byte('"') {
		advance(lexer)
	}
	return res
}
func read_comment(lexer *Lexer) string {
	res := ""
	for lexer.current_char != 0 && lexer.current_char != '\n' {
		res += string(lexer.current_char)
		advance(lexer)
	}
	if lexer.current_char == '\n' {
		advance(lexer)
	}
	return res
}

func skip_spaces(lexer *Lexer){
	for lexer.current_char != 0 && lexer.current_char == byte(' '){
		advance(lexer)
	}
}

func read_id(lexer *Lexer) Token {
	result := ""
	for lexer.current_char != 0 && (unicode.IsLetter(rune(lexer.current_char)) || unicode.IsDigit(rune(lexer.current_char))){
		//fmt.Println(string(lexer.current_char))
		result += string(lexer.current_char)
		advance(lexer)
	}
	index := find_reserved_word(result)

	if index != -1{
		
		return RESERVED_KEYWORDS[index].token
	}else{
		return create_Token(ID, result)
	}
}

func get_next_token(lexer *Lexer) Token {
	for lexer.current_char != 0 {
		//fmt.Println(string(lexer.current_char))
		switch lexer.current_char {
		case byte(' '): // skipping empty spaces
			skip_spaces(lexer)
		case byte('"'): // strings
			return create_Token(STRING, read_string(lexer))

		case byte('#'): //comments
			return create_Token(COMMENT, read_comment(lexer))

		case byte('.'): // end
			advance(lexer)
			return create_Token(DOT, ".")
		case byte(','):
			advance(lexer)
			return create_Token(COMMA, ",")
		case byte(';'):
			advance(lexer)
			return create_Token(SEMI, ";")
		case byte('\r'):
			advance(lexer)
		case byte('\n'):
			advance(lexer)
			lexer.current_line += 1
			return create_Token(NEW_LINE, "\n")
		case byte('('):
			advance(lexer)
			return create_Token(BRACKETS_L, "(")
		case byte(')'):
			advance(lexer)
			return create_Token(BRACKETS_R, ")")
		case byte (':'): // premises intro
			advance(lexer)
			return create_Token(COLON, ":")
		case byte('-'): // conclusion intro
			if peek(lexer) == byte('>'){
				advance(lexer)
				advance(lexer)
				return create_Token(ARROW, "->")
			}else{
				message("Unexpected symbol: -> was expected. line" + strconv.Itoa(lexer.current_line), lexer.project)
				return create_Token(UNEXPECTED_SYMBOL, string(lexer.current_char))
			}
		case byte('@'): // report sectioin
			advance(lexer)
			return create_Token(REPORT_SECTION, "@")
		
		default:
			
			if unicode.IsLetter(rune(lexer.current_char)){ // reading keywords or names of rules
				return read_id(lexer)
		}else{
			message("Unexpected symbol. line" + strconv.Itoa(lexer.current_line), lexer.project)
			return create_Token(UNEXPECTED_SYMBOL, string(lexer.current_char))
		}
		}
		
	}
	res := Token{
		token_type: EOF,
		value:      "EOF",
	}
	return res
}