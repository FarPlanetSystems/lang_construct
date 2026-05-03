package main

import (
	"fmt"
	"strings"

	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

// since we can import other projects, we need to be sure that none of them point to the initiate one (we want to prevent cycles)
// for this purpose we create an array containing the initiate projects, and then add all imported projects, that import other projects
// if at the moment of importation the project we import is already contained in the array, then a cycle is present
var importing_projects []*Project

type Project struct {
	AllSyntaxRules      []compiler_objects.SyntaxRule
	AllPropositions     []compiler_objects.Proposition
	AllAxioms           []compiler_objects.Axiom
	AllRules            []compiler_objects.Rule
	AllVerifiedFormulas []compiler_objects.Formula
	AllSubstitutions    []compiler_objects.Substitution
	IsVerified          bool
}

func CreateProject(syntax []compiler_objects.SyntaxRule, propositions []compiler_objects.Proposition, axioms []compiler_objects.Axiom, rules []compiler_objects.Rule, substitutions []compiler_objects.Substitution) *Project {

	res := Project{
		AllSyntaxRules:   syntax,
		AllPropositions:  propositions,
		AllAxioms:        axioms,
		AllRules:         rules,
		AllSubstitutions: substitutions,
		IsVerified:       false,
	}
	return &res
}

func (project *Project) Verify(compiler InnerCompiler) bool {

	for _, axiom := range project.AllAxioms {
		project.AllVerifiedFormulas = append(project.AllVerifiedFormulas, axiom.Formula)
	}

	for _, proposition := range project.AllPropositions {
		if !project.verifyProposition(proposition, compiler) {
			return false
		} else {
			project.AllVerifiedFormulas = append(project.AllVerifiedFormulas, proposition.Formula)
		}
	}
	return true
}

// gets a string representing an id and a LC_project where we want to find it
// it checks all rules and compares the given id with their names
// if any matches are present, returns true; otherwise - false
func (project Project) findIdInProject(id string) bool {
	for i := 0; i < len(project.AllSyntaxRules); i++ {
		if id == project.AllSyntaxRules[i].Name {
			return true
		}
	}
	return false
}

func (project Project) verifyProposition(proposition compiler_objects.Proposition, compiler InnerCompiler) bool {
	result := false
	fmt.Println("verifying proposition...")
	for _, rule := range project.AllRules {
		if project.checkRule(rule, proposition, compiler, []compiler_objects.VariableSubstitution{}) {
			fmt.Println("proposition verified")
			return true
		}
	}

	return result

}

func (project Project) checkRule(rule compiler_objects.Rule, proposition compiler_objects.Proposition, compiler InnerCompiler, previousSubstitutions []compiler_objects.VariableSubstitution) bool {
	fmt.Println("check rule..." + rule.Name)

	iteration := 1
	for true {
		// if it is a wrong substitution, we exit. If the premises are not verified, we repeat
		appliedRule := rule
		propositionFormula := proposition
		substitutions := previousSubstitutions
		//fmt.Println("iteration" + strconv.Itoa(iteration))
		//compiler_objects.PrintStatement(propositionFormula.Formula)
		// make sure it handles the pointers corecttly
		//fmt.Println("template before mathcing pattern")
		//compiler_objects.PrintStatement(propositionFormula.Formula)
		substitution, isexit := compiler.matchPattern(appliedRule, &appliedRule.Conclusion, &propositionFormula.Formula, iteration)
		//fmt.Println("template after mathcing pattern")
		//compiler_objects.PrintStatement(propositionFormula.Formula)
		if !compiler.checkVariableType(substitution) {
			//fmt.Println("EXIT EXIT EXIT EXIT")
			break
		}

		if !isexit {
			substitutions = append(substitutions, substitution)

			if project.checkRule(appliedRule, propositionFormula, compiler, substitutions) {
				return true
			}
		}

		if appliedRule.Conclusion.Compare(propositionFormula.Formula) && compiler_objects.AreSubstitutionsConsistent(substitutions) {
			var newPremises []compiler_objects.Formula

			for _, premise := range rule.Premises {
				newPremise := premise
				for _, sub := range substitutions {
					fmt.Println("substitution:")
					newPremise = replaceVariable(newPremise, sub.Variable.Id, sub.Value)
				}
				newPremises = append(newPremises, newPremise)
			}

			if project.verifyPremises(newPremises, rule, compiler) {
				return true
			}
		}
		iteration++
		if propositionFormula.Formula.IsEmpty() {
			return false
		}
	}

	return false

}

func (compiler InnerCompiler) matchPattern(rule compiler_objects.Rule, template *compiler_objects.Formula, proposition *compiler_objects.Formula, iteration int) (compiler_objects.VariableSubstitution, bool) {
	var substitution compiler_objects.VariableSubstitution

	// clears up all simple literals
	for !rule.IsParamId(template.Face().Value) && !template.IsEmpty() {
		if !compiler.MatchConstantLiteral(*proposition, *template) {
			return substitution, true
		}
		template.Dequeue()
		proposition.Dequeue()
	}
	// if either rule conclusion or proposition formula is empty, looking for a substitution makes no sence, and we remove
	if template.IsEmpty() || proposition.IsEmpty() {
		//fmt.Println("Empty")
		return substitution, true
	}
	//fmt.Println(rule.Conclusion.Face().Value)
	// we create a substitution according to the variable syntax type and iteration
	if rule.IsParamId(template.Face().Value) {
		counter := 0
		//fmt.Println("check param " + template.Face().Value)
		freeVariable := &compiler_objects.Formula{}
		freeVariableId := template.Face().Value
		//fmt.Println("grammar type of free variable " + rule.GetParamGrammarType(freeVariableId))
		freeVariableSubstitution := compiler_objects.VariableSubstitution{Variable: compiler_objects.Param{Id: freeVariableId, GrammarType: rule.GetParamGrammarType(freeVariableId)}, Value: *freeVariable}
		template.Dequeue()
		for !proposition.IsEmpty() {
			if compiler.checkVariableType(freeVariableSubstitution) && counter == iteration {
				break
			}

			freeVariable.Enqueue(proposition.Face())
			freeVariableSubstitution.Value = *freeVariable
			proposition.Dequeue()
			if compiler.checkVariableType(freeVariableSubstitution) && counter != iteration {
				counter++
			}
		}
		//fmt.Println("substituition found:")
		//compiler_objects.PrintStatement(freeVariableSubstitution.Value)
		var parameter compiler_objects.Param
		for _, param := range rule.Params {
			if freeVariableId == param.Id {
				parameter = param
			}
		}
		if !compiler.checkVariableType(freeVariableSubstitution) {
			return substitution, true
		}
		substitution = compiler_objects.VariableSubstitution{Value: *freeVariable, Variable: parameter}
	}

	return substitution, false

}

func (project Project) verifyPremises(newPremises []compiler_objects.Formula, rule compiler_objects.Rule, compiler InnerCompiler) bool {
	fmt.Println("verifying premises...")
	var substitutions []compiler_objects.VariableSubstitution
	for _, premise := range newPremises {
		for _, sub := range substitutions {
			premise = replaceVariable(premise, sub.Variable.Id, sub.Value)
		}
		//fmt.Println("premise:")
		//compiler_objects.PrintStatement(premise)
		var success bool
		var newSubstitutions []compiler_objects.VariableSubstitution
		for _, formula := range project.AllVerifiedFormulas {
			success, newSubstitutions = project.matchVerifiedFormula(premise, formula, rule, []compiler_objects.VariableSubstitution{}, compiler)
			if success {
				substitutions = append(substitutions, newSubstitutions...)
				break
			}
		}
		if !success {
			fmt.Println("premise not verified")
			return false
		}
	}
	return true
}

func (project Project) matchVerifiedFormula(premise compiler_objects.Formula, verifiedProposition compiler_objects.Formula, rule compiler_objects.Rule, previousSubstitutions []compiler_objects.VariableSubstitution, compiler InnerCompiler) (bool, []compiler_objects.VariableSubstitution) {
	//fmt.Println("checking verified formula")
	//compiler_objects.PrintStatement(verifiedProposition)
	iteration := 1
	//verifiedFormula := verifiedProposition
	for true {
		//fmt.Println("iteration " + strconv.Itoa(iteration))
		//fmt.Println("premise:")
		//compiler_objects.PrintStatement(premise)
		substitutions := previousSubstitutions
		premiseFormula := premise
		propositionFormula := verifiedProposition
		// make sure it handles the pointers corecttly
		//fmt.Println("proposition before mathcing pattern")
		//compiler_objects.PrintStatement(propositionFormula)
		substitution, isexit := compiler.matchPattern(rule, &premiseFormula, &propositionFormula, iteration)
		premiseFormula = premise
		propositionFormula = verifiedProposition
		//fmt.Println("proposition after mathcing pattern")
		//compiler_objects.PrintStatement(propositionFormula)
		//fmt.Println(isexit)
		for _, sub := range substitutions {
			premiseFormula = replaceVariable(premiseFormula, sub.Variable.Id, sub.Value)
		}
		if isexit {

			//fmt.Println("comparing: " + premise.Face().Value + " " + verifiedProposition.Face().Value)
			if premiseFormula.Compare(verifiedProposition) && compiler_objects.AreSubstitutionsConsistent(substitutions) {
				fmt.Println("premise verified")
				return true, substitutions
			}
			return false, substitutions
		}
		substitutions = append(substitutions, substitution)
		success, newSubs := project.matchVerifiedFormula(premiseFormula, verifiedProposition, rule, substitutions, compiler)
		if success {
			return true, newSubs
		}
		iteration++
	}
	fmt.Println("fail")
	return false, previousSubstitutions
}

func replaceVariable(formula compiler_objects.Formula, variableId string, insertedTokens compiler_objects.Formula) compiler_objects.Formula {
	//fmt.Println("replacing variable: " + variableId)
	//compiler_objects.PrintStatement(insertedTokens)
	var result compiler_objects.Formula
	tokens := insertedTokens
	for !formula.IsEmpty() {
		if formula.Face().Value == variableId {
			//fmt.Println(variableId + " found")
			formula.Dequeue()
			for !tokens.IsEmpty() {
				result.Enqueue(tokens.Face())
				tokens.Dequeue()
			}
			tokens = insertedTokens
		} else {
			result.Enqueue(formula.Face())
			formula.Dequeue()
		}
	}
	return result
}

func (compiler InnerCompiler) checkVariableType(sub compiler_objects.VariableSubstitution) bool {
	syntax := compiler.findSyntaxRule(sub.Variable.GrammarType)
	//fmt.Println("checking variable type " + syntax.Name + " for variable " + sub.Value.Face().Value)
	innerCompiler := &compiler

	if innerCompiler.innerParseFormula(syntax, sub.Value).isError {

		return false
	}

	if !innerCompiler.IsParsedSuccessfully {
		return false
	}
	return true
}

/*
func (compiler InnerCompiler) MatchPattern(propositionFormula compiler_objects.Formula, template compiler_objects.Formula, rule compiler_objects.Rule, iteration int) (bool, []compiler_objects.VariableSubstitution) {
	var substitutions []compiler_objects.VariableSubstitution
	for !rule.Conclusion.IsEmpty() {
		success, substitution := compiler.MatchLiteral(propositionFormula, template, rule)

		if !success {
			fmt.Println("substitution failed: unable to match pattern")

			return false, substitutions
		}

		success, substitutions = compiler_objects.AddToSubstitutions(substitutions, substitution)

		if !success {
			fmt.Println("substitution failed: different patterns were matched to the same variable")

			return false, substitutions
		}

		fmt.Println("Free variable content:")
		compiler_objects.PrintStatement(substitution.Value)
		if !success {
			return false, substitutions
		}
		substitutions = append(substitutions, substitution)
	}

	return true, substitutions
}*/

func (compiler InnerCompiler) MatchLiteral(propositionFormula compiler_objects.Formula, template compiler_objects.Formula, rule compiler_objects.Rule) (bool, compiler_objects.VariableSubstitution) {
	var substitution compiler_objects.VariableSubstitution
	if rule.IsParamId(rule.Conclusion.Face().Value) {
		fmt.Println("check param " + rule.Conclusion.Face().Value)
		freeVariable := &compiler_objects.Formula{}
		freeVariableId := rule.Conclusion.Face().Value
		//fmt.Println("grammar type of free variable " + rule.GetParamGrammarType(freeVariableId))
		freeVariableSubstitution := compiler_objects.VariableSubstitution{Variable: compiler_objects.Param{Id: freeVariableId, GrammarType: rule.GetParamGrammarType(freeVariableId)}, Value: *freeVariable}
		rule.Conclusion.Dequeue()
		for (propositionFormula.Face().Value != rule.Conclusion.Face().Value || !compiler.checkVariableType(freeVariableSubstitution)) && !propositionFormula.IsEmpty() {
			freeVariable.Enqueue(propositionFormula.Face())
			freeVariableSubstitution.Value = *freeVariable
			propositionFormula.Dequeue()
		}
		var parameter compiler_objects.Param
		for _, param := range rule.Params {
			if freeVariableId == param.Id {
				parameter = param
			}
		}
		substitution = compiler_objects.VariableSubstitution{Value: *freeVariable, Variable: parameter}

	} else {
		if !compiler.MatchConstantLiteral(propositionFormula, template) {
			return false, substitution
		}
		rule.Conclusion.Dequeue()
		propositionFormula.Dequeue()
	}
	return true, substitution
}

func (compiler InnerCompiler) MatchVariableLiteral(propositionFormula compiler_objects.Formula, rule compiler_objects.Rule) {

}

func (compiler InnerCompiler) MatchConstantLiteral(propositionFormula compiler_objects.Formula, template compiler_objects.Formula) bool {
	var pattern compiler_objects.Token
	//fmt.Println("check pattern " + template.Face().Value)
	pattern = template.Face()

	if strings.Compare(propositionFormula.Face().Value, pattern.Value) != 0 {
		//fmt.Println("pattern not found: " + propositionFormula.Face().Value + " " + pattern.Value)
		return false
	}
	return true
}
