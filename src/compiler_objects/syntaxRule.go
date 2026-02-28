package compiler_objects

import (
	"fmt"
)

type SyntaxRule struct {
	Options []SyntaxOption
	Name    string
}

func CreateSyntaxRule(name string, options []SyntaxOption) SyntaxRule {
	res := SyntaxRule{
		Name:    name,
		Options: options,
	}
	return res
}

func (rule SyntaxRule) PrintSyntaxRule() {
	fmt.Println("printing syntax rule: " + rule.Name)
	fmt.Println("number of options " + fmt.Sprint(len(rule.Options)))
	for num, option := range rule.Options {
		fmt.Println("option " + fmt.Sprint(num))
		option.PrintSyntaxOption()
	}
}

func (rule SyntaxRule) ContainsEOFOption() bool {
	for _, option := range rule.Options {
		if option.HeadWord.Content.TokenType == INNER_EOF {
			return true
		}
	}
	return false
}

func (option SyntaxOption) PrintSyntaxOption() {
	opt := option
	for !opt.IsEmpty {
		fmt.Println(opt.HeadWord)
		opt.Dequeue()
	}
}

type SyntaxOption struct {
	HeadWord *GrammarWord
	tailWord *GrammarWord
	IsEmpty  bool
}

// syntax options are not parsed correctly : they are all empty. Just run the programm
func CreateSyntaxOption() SyntaxOption {
	return SyntaxOption{
		HeadWord: nil,
		tailWord: nil,
		IsEmpty:  true,
	}
}

func CreateEOFSyntaxOption() SyntaxOption {
	EOFOption := CreateSyntaxOption()
	EOFToken := CreateToken(INNER_EOF, "")
	EOFOption.Enqueue(GrammarWord{Content: EOFToken})
	return EOFOption
}

func (option *SyntaxOption) Enqueue(word GrammarWord) {
	if option.IsEmpty {
		option.HeadWord = &word
		option.tailWord = &word
		option.IsEmpty = false
	} else {
		option.tailWord.PreviousWord = &word
		option.tailWord = &word
	}
}

func (option *SyntaxOption) Dequeue() {
	if option.IsEmpty {
		return
	}
	if option.HeadWord.PreviousWord == nil {
		option.HeadWord = nil
		option.tailWord = nil
		option.IsEmpty = true
	} else {
		option.HeadWord = option.HeadWord.PreviousWord
	}
}

type GrammarWord struct {
	Content      Token
	PreviousWord *GrammarWord
}

// a very problematic recursion
// because of possible cycles in defining the grammar graph

/*
type GrammarTreeNode struct {
	isTerminate bool
	token       Token
	children    []*GrammarTreeNode
}

func createGrammarTreeNode(isTerminate bool, token Token) GrammarTreeNode {
	res := GrammarTreeNode{
		isTerminate: isTerminate,
		token:       token,
		children:    []*GrammarTreeNode{},
	}
	return res
}

func (node GrammarTreeNode) insertChild(child *GrammarTreeNode) {
	node.children = append(node.children, child)
}

func (node GrammarTreeNode) getFirstSet() []Token {
	res := []Token{}
	for _, child := range node.children {
		if child.isTerminate {
			res = append(res, child.token)
		} else {
			res = append(res, child.getFirstSet()...)
		}
	}
	return res
}

func (compiler *Compiler)createGrammarTree(rules []SyntaxRule) GrammarTreeNode {
	root := createGrammarTreeNode(false, Token{})
	stmt := findSyntaxRule(rules, "statement")
	if stmt.options == nil{
		compiler.messanger.insertMessage("no root syntax rule defined", -1)
		return GrammarTreeNode{}
	}

}

func convertRuleToTreeNode (rule SyntaxRule) GrammarTreeNode{

}

*/
