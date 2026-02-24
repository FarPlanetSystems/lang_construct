package compiler_objects

// Token types
const STRING = "STRING"

const ARROW = "ARROW"
const COLON = "COLON"
const COLON_EQUAL = "COLON_EQUAL"
const COMMA = "COMMA"
const DOT = "DOT"
const NEW_LINE = "NEW_LINE"
const SEMI_COLON = "SEMI_COLON"
const FROM = "FROM"
const UNEXPECTED_SYMBOL = "UNEXPECTED_SYMBOL"
const COMMENT = "COMMENT"
const ID = "ID"
const HAVE = "HAVE"
const RULE = "RULE"
const DEF = "DEF"
const IF = "IF"
const SPEC = "SPEC"
const IMPORT = "IMPORT"
const REPORT_SECTION = "REPORT_SECTION"
const BRACKETS_R = "BRACKETS_R"
const BRACKETS_L = "BRACKETS_L"
const CURL_BRACKETS_R = "CURL_BRACKETS_R"
const CURL_BRACKETS_L = "CURL_BRACKETS_L"
const SQ_BRACKEtS_R = "SQ_BRACKETS_R"
const SQ_BRACKETS_L = "SQ_BRACKETS_R"
const ANY = "ANY"
const EOF = "EOF"
const CONCLUSION_TOKEN = "CONCLUSION_TOKEN"

type Token struct {
	TokenType string
	Value     string
}

func CreateToken(tokenType string, value string) Token {
	res := Token{
		TokenType: tokenType,
		Value:     value,
	}
	return res
}
