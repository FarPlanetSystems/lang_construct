package main

import (
	"unicode"

	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

type RESERVED_KEYWORD struct {
	keyword string
	token   compiler_objects.Token
}

var RESERVED_KEYWORDS = [8]RESERVED_KEYWORD{
	{
		keyword: "rule",
		token: compiler_objects.Token{
			TokenType: compiler_objects.RULE,
			Value:     "rule",
		},
	},

	{
		keyword: "EOF",
		token: compiler_objects.Token{
			TokenType: compiler_objects.INNER_EOF,
			Value:     "EOF",
		},
	},

	// remove this
	{
		keyword: "have",
		token: compiler_objects.Token{
			TokenType: compiler_objects.HAVE,
			Value:     "have",
		},
	},
	{
		keyword: "def",
		token: compiler_objects.Token{
			TokenType: compiler_objects.DEF,
			Value:     "def",
		},
	},

	{
		keyword: "from",
		token: compiler_objects.Token{
			TokenType: compiler_objects.FROM,
			Value:     "from",
		},
	},
	{
		keyword: "import",
		token: compiler_objects.Token{
			TokenType: compiler_objects.IMPORT,
			Value:     "import",
		},
	},
	{
		keyword: "if",
		token: compiler_objects.Token{
			TokenType: compiler_objects.IF,
			Value:     "if",
		},
	},
	{
		keyword: "spec",
		token: compiler_objects.Token{
			TokenType: compiler_objects.SPEC,
			Value:     "spec",
		},
	},
}

// gets a string which is supposed to be a reserved key word
// if it is, returns the index of it in RESERVED_WORDS array
// otherwise returns -1
func findReservedWord(word string) int {
	for i := 0; i < len(RESERVED_KEYWORDS); i++ {
		if RESERVED_KEYWORDS[i].keyword == word {
			return i
		}
	}
	return -1
}

const LEXER_MODE_DEFAULT = "LEXER_MODE_DEFAULT"
const LEXER_MODE_HAVE = "LEXER_MODE_HAVE"

type Lexer struct {
	text         string
	pos          int
	current_char byte
	current_line int
	mode         string
}

func CreateLexer(text string) *Lexer {
	res := Lexer{
		text:         text,
		pos:          0,
		current_char: byte(text[0]),
		current_line: 1,
		mode:         LEXER_MODE_DEFAULT,
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

func peekString(lexer *Lexer, steps int) string {
	peek_pos := lexer.pos
	res := ""
	for i := 0; i < steps; i++ {
		peek_pos += 1
		if peek_pos > len(lexer.text)-1 {
			return ""
		} else {
			res += string(lexer.text[peek_pos])
		}
	}
	return res

}

func readString(lexer *Lexer) string {

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
func readComment(lexer *Lexer) string {
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

func skipSpaces(lexer *Lexer) {
	for lexer.current_char != 0 && lexer.current_char == byte(' ') {
		advance(lexer)
	}
}

func readId(lexer *Lexer) compiler_objects.Token {
	result := ""
	for isIdCharCorrect(lexer.current_char) {
		//fmt.Println(string(lexer.current_char))
		result += string(lexer.current_char)
		advance(lexer)
	}
	index := findReservedWord(result)

	if index != -1 {
		return RESERVED_KEYWORDS[index].token
	}
	return compiler_objects.CreateToken(compiler_objects.ID, result)
}

func readConclusionToken(lexer *Lexer) compiler_objects.Token {
	result := ""
	for isIdCharCorrect(lexer.current_char) {
		result += string(lexer.current_char)
		advance(lexer)
	}
	return compiler_objects.CreateToken(compiler_objects.CONCLUSION_TOKEN, result)
}

// the function checks if a char can be represented in a id string
// it gets a char and returns true when it is whether a letter or a digit or an underscore
// otherwise false
func isIdCharCorrect(id_char byte) bool {
	if id_char == 0 {
		return false
	}
	switch string(id_char) {
	case " ":
		return false
	case ";":
		return false
	}
	/*
		if unicode.IsLetter(rune(id_char)) || unicode.IsDigit(rune(id_char)) || string(id_char) == "_" {
			return true
		}
	*/
	return true
}

func (lexer *Lexer) getNextToken(messager *Messenger) compiler_objects.Token {
	skipSpaces(lexer)
	for lexer.current_char != 0 {
		//fmt.Println(string(lexer.current_char))
		//fmt.Println("mode " + lexer.mode)
		if lexer.mode == LEXER_MODE_HAVE {
			switch lexer.current_char {
			case byte(';'):
				advance(lexer)
				//fmt.Println("changing mode def")
				(*lexer).mode = LEXER_MODE_DEFAULT
				return compiler_objects.CreateToken(compiler_objects.SEMI_COLON, ";")

			default:
				return readConclusionToken(lexer)
			}
		}
		switch lexer.current_char {
		case byte('"'): // strings
			return compiler_objects.CreateToken(compiler_objects.STRING, readString(lexer))

		case byte('#'): //comments
			lexer.current_line += 1
			return compiler_objects.CreateToken(compiler_objects.COMMENT, readComment(lexer))
		case byte('.'): // end
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.DOT, ".")
		case byte(','):
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.COMMA, ",")
		case byte(';'):
			advance(lexer)
			(*lexer).mode = LEXER_MODE_DEFAULT
			//fmt.Println("changing mode def")
			return compiler_objects.CreateToken(compiler_objects.SEMI_COLON, ";")
		case byte('\r'):
			advance(lexer)
		case byte('\n'):
			advance(lexer)
			lexer.current_line += 1
		case byte('('):
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.BRACKETS_L, "(")
		case byte(')'):
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.BRACKETS_R, ")")
		case byte('{'):
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.CURL_BRACKETS_L, "{")
		case byte('}'):
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.CURL_BRACKETS_R, "}")
		case byte('$'):
			if peekString(lexer, 3) == "any" {
				advance(lexer)
				advance(lexer)
				advance(lexer)
				advance(lexer)
				return compiler_objects.CreateToken(compiler_objects.ANY, "$any")
			} else {
				messager.InsertMessage("Unexpected symbol: $any was expected", lexer.current_line)
				return compiler_objects.CreateToken(compiler_objects.UNEXPECTED_SYMBOL, string(lexer.current_char))
			}
		case byte(':'): // premises intro grammar definition

			if peekString(lexer, 1) == "=" {
				advance(lexer)
				advance(lexer)
				return compiler_objects.CreateToken(compiler_objects.COLON_EQUAL, ":=")
			}
			return compiler_objects.CreateToken(compiler_objects.COLON, ":")
		case byte('-'): // conclusion intro
			if peek(lexer) == byte('>') {
				advance(lexer)
				advance(lexer)
				return compiler_objects.CreateToken(compiler_objects.ARROW, "->")
			} else {
				messager.InsertMessage("Unexpected symbol: -> was expected.", lexer.current_line)
				return compiler_objects.CreateToken(compiler_objects.UNEXPECTED_SYMBOL, string(lexer.current_char))
			}
		case byte('@'): // report sectioin
			advance(lexer)
			return compiler_objects.CreateToken(compiler_objects.REPORT_SECTION, "@")

		default:
			if unicode.IsLetter(rune(lexer.current_char)) || unicode.IsDigit(rune(lexer.current_char)) { // reading keywords or names of rules
				idToken := readId(lexer)
				if idToken.TokenType == compiler_objects.HAVE {
					//fmt.Println("changing mode have")
					(*lexer).mode = LEXER_MODE_HAVE
				}
				return idToken
			} else {
				messager.InsertMessage("Unexpected symbol: "+string(lexer.current_char), lexer.current_line)
				return compiler_objects.CreateToken(compiler_objects.UNEXPECTED_SYMBOL, string(lexer.current_char))
			}

		}

	}
	res := compiler_objects.Token{
		TokenType: compiler_objects.EOF,
		Value:     "EOF",
	}
	return res
}
