package main

import (
	"slices"

	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

type InnerCompiler struct {
	IsParsedSuccessfully bool
	// it would be more wise if there were an internal messenger istead of importing from external
	Messenger       *Messenger
	allPropositions []compiler_objects.Proposition
	allSyntaxRules  []compiler_objects.SyntaxRule
	allAxioms       []compiler_objects.Axiom
}

func CreateInnerCompiler(project Project) *InnerCompiler {
	return &InnerCompiler{
		IsParsedSuccessfully: true,
		Messenger:            &Messenger{},
		allPropositions:      project.AllPropositions,
		allSyntaxRules:       project.AllSyntaxRules,
		allAxioms:            project.AllAxioms,
	}
}

func (compiler *InnerCompiler) InnerParse() {
	root := compiler.findSyntaxRule("proposition")
	//fmt.Println("root: ")
	//root.PrintSyntaxRule()
	for _, proposition := range compiler.allPropositions {
		//compiler_objects.PrintStatement(proposition.Formula)
		err := compiler.innerParseFormula(root, proposition.Formula)
		if err.isError {
			compiler.Messenger.InsertMessage(err.Message, err.line)
			compiler.IsParsedSuccessfully = false
		}
	}

	for _, axiom := range compiler.allAxioms {
		err := compiler.innerParseFormula(root, axiom.Formula)
		if err.isError {
			compiler.Messenger.InsertMessage(err.Message, err.line)
			compiler.IsParsedSuccessfully = false
		}
	}
}

func (compiler *InnerCompiler) innerParseFormula(root compiler_objects.SyntaxRule, formula compiler_objects.Formula) Error {

	result, err := compiler.innerParseFormulaHandleToken(root, formula)

	if !result.IsEmpty() {
		err.isError = true
		err.Message = "inner parsing error: found " + result.Face().Value + " while end of expression was expected"
		err.line = result.Line
	}
	return err
}

func (compiler *InnerCompiler) innerParseFormulaHandleToken(root compiler_objects.SyntaxRule, formula compiler_objects.Formula) (*compiler_objects.Formula, Error) {
	result := &formula
	option, err := compiler.decideOption(root, formula)

	// inner eat
	for !option.IsEmpty && !err.isError {
		//fmt.Println("remaining tokens to parse")
		//compiler_objects.PrintStatement(result.Conclusion)
		//fmt.Println("remaining tokens to expect")
		//option.PrintSyntaxOption()
		switch option.HeadWord.Content.TokenType {
		case compiler_objects.ID:
			result, err = compiler.innerParseFormulaHandleIDToken(&option, result)
		case compiler_objects.STRING:
			result, err = compiler.innerParseFormulaHandleSTRINGToken(&option, result)
		case compiler_objects.INNER_EOF:
			option.Dequeue()
		}
	}

	return result, err
}

func (compiler *InnerCompiler) innerParseFormulaHandleIDToken(option *compiler_objects.SyntaxOption, result *compiler_objects.Formula) (*compiler_objects.Formula, Error) {

	//fmt.Println("handle ID:" + option.HeadWord.Content.Value)
	err := createError()
	nextRuleName := option.HeadWord.Content.Value
	nextRule := compiler.findSyntaxRule(nextRuleName)
	if nextRule.Name != "" {
		result, err = compiler.innerParseFormulaHandleToken(nextRule, *(result))
	} else {
		err.isError = true
		err.Message = "inner parsing error: unknown grammar rule name found: " + nextRuleName
		err.line = result.Line
	}
	option.Dequeue()
	return result, err
}

func (compiler *InnerCompiler) innerParseFormulaHandleSTRINGToken(option *compiler_objects.SyntaxOption, result *compiler_objects.Formula) (*compiler_objects.Formula, Error) {
	//fmt.Println("handle string:" + option.HeadWord.Content.Value)
	err := createError()
	if result.IsEmpty() {
		err.Message = "inner parsing error: no correct syntax rule for the expression was found. "
		err.line = result.Line
		err.isError = true
		return result, err
	}
	if result.HeadToken.Content.Value != option.HeadWord.Content.Value {
		err.Message = "inner parsing error: found " + result.Face().Value + " while " + option.HeadWord.Content.Value + " was expected"
		err.isError = true
		err.line = result.Line
		return result, err
	}
	result.Dequeue()
	option.Dequeue()
	return result, err

}

// remember that every grammar rule option has its one unique first set
func (compiler *InnerCompiler) decideOption(rule compiler_objects.SyntaxRule, formula compiler_objects.Formula) (compiler_objects.SyntaxOption, Error) {
	//fmt.Println("determining rule")
	//fmt.Println("initial rule:")
	//rule.PrintSyntaxRule()
	err := createError()
	first := formula.Face()
	for _, option := range rule.Options {
		//fmt.Println("deciding option: " + option.HeadWord.Content.TokenType)

		switch option.HeadWord.Content.TokenType {
		case compiler_objects.STRING:
			if option.HeadWord.Content.Value == first.Value {
				//fmt.Println("String")
				return option, err
			}

		case compiler_objects.ID:
			nextRule := compiler.findSyntaxRule(option.HeadWord.Content.Value)
			//fmt.Println("first set: " + nextRule.Name)
			//fmt.Println(compiler.getFirstSet(nextRule))
			if slices.Contains(compiler.getFirstSet(nextRule), first.Value) {
				//fmt.Println("ID")
				return option, err
			}

		case compiler_objects.INNER_EOF:
			if formula.IsEmpty() {
				//fmt.Println("EOF")
				return compiler_objects.CreateEOFSyntaxOption(), err
			}
		default:
			//fmt.Println("default")
		}
	}

	if !formula.IsEmpty() && rule.ContainsEOFOption() {
		return compiler_objects.CreateEOFSyntaxOption(), err
	}
	err.isError = true
	err.Message = "no correct syntax rule for the schema was found"
	err.line = formula.Line
	return compiler_objects.CreateSyntaxOption(), err
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
