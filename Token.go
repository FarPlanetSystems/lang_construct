package main

//Token types
const STRING = "STRING"

const ARROW = "ARROW"
const COLON = "COLON"
const COMMA = "COMMA"
const DOT = "DOT"
const NEW_LINE = "NEW_LINE"
const SEMI = "SEMI"
const FROM = "FROM"
const UNEXPECTED_SYMBOL = "UNEXPECTED_SYMBOL"
const COMMENT = "COMMENT"
const ID = "ID"
const HAVE = "HAVE"
const RULE = "RULE"
const DEF = "DEF"
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

type Token struct{
	token_type string
	value string 
}

func create_Token(token_type string, value string) Token{
	res := Token{
		token_type: token_type,
		value: value,
	}	
	return res
}

