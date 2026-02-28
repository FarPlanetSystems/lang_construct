package main

import (
	"fmt"
	"slices"

	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

type InnerCompiler struct {
	IsParsedSuccessfully bool
	// it would be more wise if there were an internal messenger istead of importing from external
	Messenger       *Messenger
	allPropositions []compiler_objects.Proposition
	allSyntaxRules  []compiler_objects.SyntaxRule
}

func CreateInnerCompiler(project Project) *InnerCompiler {
	return &InnerCompiler{
		IsParsedSuccessfully: true,
		Messenger:            &Messenger{},
		allPropositions:      project.allPropositions,
		allSyntaxRules:       project.allSyntaxRules,
	}
}

func (compiler *InnerCompiler) InnerParse() {
	root := compiler.findSyntaxRule("proposition")
	fmt.Println("root: ")
	root.PrintSyntaxRule()
	for _, proposition := range compiler.allPropositions {
		compiler_objects.PrintStatement(proposition.Conclusion)
		result := compiler.innerParseProposition(root, proposition)
		if !result.Conclusion.IsEmpty() {
			compiler.Messenger.InsertMessage("inner parsing error: found "+result.Conclusion.Face().Value+
				" while EOF was expected", result.Line)

			compiler.IsParsedSuccessfully = false
		}
	}
}

// next step: find out how to simplify the algorithm for inner parse. Extend it for more complicate syntaxes
func (compiler *InnerCompiler) innerParseProposition(root compiler_objects.SyntaxRule, proposition compiler_objects.Proposition) *compiler_objects.Proposition {
	result := &proposition

	option := compiler.decideOption(root, proposition)

	// inner eat
	for !option.IsEmpty {
		fmt.Println("remaining tokens to parse")
		compiler_objects.PrintStatement(result.Conclusion)
		fmt.Println("remaining tokens to expect")
		option.PrintSyntaxOption()
		switch option.HeadWord.Content.TokenType {
		case compiler_objects.ID:
			result = compiler.innerParsePropositionHandleIDToken(&option, result)
		case compiler_objects.STRING:
			result = compiler.innerParsePropositionHandleSTRINGToken(&option, result)
		case compiler_objects.INNER_EOF:
			option.Dequeue()
		}
	}

	return result
}

func (compiler *InnerCompiler) innerParsePropositionHandleIDToken(option *compiler_objects.SyntaxOption, result *compiler_objects.Proposition) *compiler_objects.Proposition {

	//fmt.Println("handle ID:" + option.HeadWord.Content.Value)

	nextRuleName := option.HeadWord.Content.Value
	nextRule := compiler.findSyntaxRule(nextRuleName)
	if nextRule.Name != "" {
		result = compiler.innerParseProposition(nextRule, *(result))
	} else {
		compiler.Messenger.InsertMessage("inner parsing error: unknown grammar rule name found: "+nextRuleName, result.Line)
		compiler.IsParsedSuccessfully = false
	}
	option.Dequeue()
	return result
}

func (compiler *InnerCompiler) innerParsePropositionHandleSTRINGToken(option *compiler_objects.SyntaxOption, result *compiler_objects.Proposition) *compiler_objects.Proposition {
	//fmt.Println("handle string:" + option.HeadWord.Content.Value)
	if result.Conclusion.IsEmpty() {
		compiler.Messenger.InsertMessage("inner parsing error: no correct error for the expression was found. ", result.Line)
		compiler.IsParsedSuccessfully = false
	}
	if result.Conclusion.HeadToken.Content.Value != option.HeadWord.Content.Value {

		compiler.Messenger.InsertMessage("inner parsing error: found "+result.Conclusion.Face().Value+
			" while "+option.HeadWord.Content.Value+" was expected", result.Line)

		compiler.IsParsedSuccessfully = false
	}
	result.Conclusion.Dequeue()
	option.Dequeue()
	return result

}

// remember that every grammar rule option has its one unique first set
func (compiler *InnerCompiler) decideOption(rule compiler_objects.SyntaxRule, proposition compiler_objects.Proposition) compiler_objects.SyntaxOption {
	fmt.Println("determining rule")
	fmt.Println("initial rule:")
	rule.PrintSyntaxRule()
	first := proposition.Conclusion.Face()
	for _, option := range rule.Options {
		//fmt.Println("deciding option: " + option.HeadWord.Content.TokenType)

		switch option.HeadWord.Content.TokenType {
		case compiler_objects.STRING:
			if option.HeadWord.Content.Value == first.Value {
				//fmt.Println("String")
				return option
			}

		case compiler_objects.ID:
			nextRule := compiler.findSyntaxRule(option.HeadWord.Content.Value)
			fmt.Println("first set: " + nextRule.Name)
			fmt.Println(compiler.getFirstSet(nextRule))
			if slices.Contains(compiler.getFirstSet(nextRule), first.Value) {
				//fmt.Println("ID")
				return option
			}

		case compiler_objects.INNER_EOF:
			if proposition.Conclusion.IsEmpty() {
				//fmt.Println("EOF")
				return compiler_objects.CreateEOFSyntaxOption()
			}
		default:
			//fmt.Println("default")
		}
	}

	if !proposition.Conclusion.IsEmpty() && rule.ContainsEOFOption() {
		return compiler_objects.CreateEOFSyntaxOption()
	}

	compiler.Messenger.InsertMessage("no correct syntax rule for the schema was found", proposition.Line)
	compiler.IsParsedSuccessfully = false
	return compiler_objects.CreateSyntaxOption()
}

func (compiler *InnerCompiler) getFirstSet(rule compiler_objects.SyntaxRule) []string {
	res := []string{}
	for _, option := range rule.Options {
		switch option.HeadWord.Content.TokenType {
		case compiler_objects.STRING:
			res = append(res, option.HeadWord.Content.Value)
		case compiler_objects.ID:
			optionRule := compiler.findSyntaxRule(option.HeadWord.Content.Value)
			optionTokens := compiler.getFirstSet(optionRule)
			res = append(res, optionTokens...)
		case compiler_objects.INNER_EOF:
			res = append(res, "")
		}
	}
	return res
}

func (compiler *InnerCompiler) findSyntaxRule(name string) compiler_objects.SyntaxRule {
	for _, rule := range compiler.allSyntaxRules {
		if rule.Name == name {
			return rule
		}
	}
	return compiler_objects.SyntaxRule{}
}
