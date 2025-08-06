package main

import (
	"fmt"
	"strconv"
)

type Statement struct {
	rule_name   string
	conclusion  Argument
	params      []string
	premises    []Argument
	line int
}

func createStatement(ruleName string, concusion Argument, params []string, premises []Argument, line int, project *LC_project) Statement {
	res := Statement{
		rule_name:  ruleName,
		conclusion: concusion,
		params:     params,
		premises:   premises,
		line: line,
	}
	project.all_statements = append(project.all_statements, res)
	return res
}

func deepCopyStatement(statement Statement) Statement{
	var newStatement Statement
	newStatement.rule_name = statement.rule_name
	newStatement.conclusion = statement.conclusion
	newStatement.params = append(newStatement.params, statement.params...)
	newStatement.premises = append(newStatement.premises, statement.premises...)
	newStatement.line = statement.line
	return newStatement
}
// func gets a statement which coherence we want to check in project
// it returns true if it is coherent, otherwise - false
func verifyStatement(statement Statement, project *LC_project) bool{
	presentStatement := deepCopyStatement(statement)
	//we look for a rule in project.all_rules reference to which must be contained in statement.rule_name
	var appliedRule Rule
	for i := 0; i < len(project.all_rules); i++{
		if project.all_rules[i].name == presentStatement.rule_name{
			appliedRule = deepCopyRule(project.all_rules[i])
			
		}
	}
	// applied rule.name being empty indicates that there is no such rule in project.all_rules. In this case we message an error and return false
	if appliedRule.name == ""{
		message("no rule " + presentStatement.rule_name + " was found. Line " + strconv.Itoa(presentStatement.line), project)
		return false
	}
	// if the number of params in applied_rule is not equal to that in present_statement, we message an error and return false
	// However, if applied_rule.are_any_params is true, we do not check the number of params
	if len(appliedRule.params) != len(presentStatement.params) && !appliedRule.are_any_params {
		message("derriving a statement, there must be as many parameters as there defined in the applied rule. Line "  + strconv.Itoa(presentStatement.line), project)
		return false
	}
	// if the number of premises given in applied_rule is not equal to that in present_statement, we try to find any sound premises that could replace the initiate ones
	if len(appliedRule.premises) != len(presentStatement.premises) && !appliedRule.are_any_premisses {
		soundPremises := findSoundPremises(appliedRule, presentStatement, project)
		if len(soundPremises) > 0{
			presentStatement.premises = append(presentStatement.premises, soundPremises...)
		}else{
		// if there no such premises we message an error and return false
		message("derriving a statement, there must be as many premises as there defined in the applied rule. Line "  + strconv.Itoa(presentStatement.line), project)
		return false
		}
	}
	
	if !checkRuleApplicability(presentStatement, appliedRule, project){
		return false
	}
	if !arePremisesVerified(presentStatement.premises, *project){
		message("not all premises are verified. Line " + strconv.Itoa(presentStatement.line), project)
		return false
	}
	project.all_legal_expressions = append(project.all_legal_expressions, presentStatement.conclusion)
	return true
}
// func gets a rule and a statement we have applied the rule on
// it returns another rule being a copy of the initiate rule which params contained in the conclusion and premises are exchanged with arguments given in the statement
func substitudeRuleWithParams(statement Statement, rule Rule) Rule {
	substitutedRule := deepCopyRule(rule)
	for i := 0; i<len(substitutedRule.params); i++{
		sequence := "[" + substitutedRule.params[i] + "]"
		// replacing params signs in premises of the rule (rule.params) with expressions in statements as arguments (statement.params)

		for j := 0; j < len(substitutedRule.premises); j++{
			
			
			substitutedRule.premises[j] = substitutedRule.premises[j].Replace(sequence, statement.params[i])
		}
		for j := 0; j < len(substitutedRule.conclusions); j++{
			
			substitutedRule.conclusions[j] = substitutedRule.conclusions[j].Replace(sequence, statement.params[i])
		}
	}

	return substitutedRule
}
// func checks if a given statement can be correctly infered from a given statement considering all legal expressions from the current project
// if it can, returns true. Otherwise - false
func checkRuleApplicability(statement Statement, rule Rule, project *LC_project) bool{
	substitutedRule := substitudeRuleWithParams(statement, rule)
	// checking if there is correspondece with the statement's conclusion with one of the rule's conclusion
	correspondeceFound := false


	for i := 0; i < len(substitutedRule.conclusions); i++{

		if compareArguments(substitutedRule.conclusions[i], statement.conclusion){
			correspondeceFound = true
		}
	}
	if !correspondeceFound{
		msgLine := "conclusion " + statement.conclusion.ToString() + " does not correspond to any conclusion of the rule " + substitutedRule.name + ". Line " + strconv.Itoa(statement.line) + "\n See:"
		message(msgLine, project)
		for i := 0; i<len(substitutedRule.conclusions); i++{
			message(substitutedRule.conclusions[i].ToString(), project)
		}
		return false
	}
	// checking the correspondence among premises
	for i := 0; i < len(substitutedRule.premises); i++{
		fmt.Println("statement premises " + statement.premises[0].ToString())
		fmt.Println("rule premises " + substitutedRule.premises[0].ToString())
		if !compareArguments( substitutedRule.premises[i], statement.premises[i]) {
			msgLine := "a premise "+ statement.premises[i].ToString() + " does not correspond to the required one " + substitutedRule.premises[i].ToString() + ". Line " + strconv.Itoa(statement.line) + "\n See:"
			message(msgLine, project)
			message(substitutedRule.premises[i].ToString() + " was expected, but " + statement.premises[i].ToString() + " was found", project)
			return false
		}
	}
	return true
}

func arePremisesVerified(premises []Argument, project LC_project) bool{
	for i := 0; i < len(premises); i++{
		isPremiseFound := false
		for j:=0; j < len(project.all_legal_expressions); j++{
			if compareArguments(project.all_legal_expressions[j], premises[i]){
				isPremiseFound = true
				break
			}
		}
		if !isPremiseFound{
			return false
		}
	}
	return true
}
// func gets a rule we want to apply in order to verify an expression, the expression itself, which containS no premises, and the current project
// it looks for a set of premises in project.all_legal_expressions which complete the given statement to a one that can be verified with the rule
// if there is such a set of legal expressions, it returns an array of strings representing it
// otherwise returns an empty array
func findSoundPremises(rule Rule, statement Statement, project *LC_project)[]Argument{
	soundStatement := deepCopyStatement(statement)
	if len(rule.premises) != 0 && len(statement.premises) == 0{
		possiblePremises := get_all_k_elements_premises(len(rule.premises), []Argument{}, project)
		for i := 0; i < len(possiblePremises); i++{
			match := true
			soundStatement.premises = possiblePremises[i]
			substitutedRule := substitudeRuleWithParams(soundStatement, rule)
				for i := 0; i < len(substitutedRule.premises); i++{
					fmt.Println("rule premise " + substitutedRule.premises[i].ToString())
					fmt.Println("statement premise" + soundStatement.premises[i].ToString())
					if !compareArguments(substitutedRule.premises[i], soundStatement.premises[i]){
						match =  false
					}
			}
			if match == true{
				return possiblePremises[i]
			}
		}
	}
	return []Argument{}
}
// recursive func gets an integer k represnting the size of permutation array, k_element, that must be an empty string array, and the current project
// returns an array of arrays res representing all possible k-sized permutations of all_legal_expressions array in project
func get_all_k_elements_premises(k int, element []Argument, project *LC_project)[][]Argument{
	res:= [][]Argument{}
	for i := 0; i<len(project.all_legal_expressions); i++{
		// k_element contains one single permutation
		// in each iteration we copy k_element in order to add the following element in each step of recursion
		k_element_clone := []Argument{}
		copy(k_element_clone, element)
		k_element_clone = append(element, project.all_legal_expressions[i])
		// as we achieve the final step of recurcion (k), we add the collected permutation k_element to res
		if k == 1{
			res = append(res, k_element_clone)
		// in all other steps we go further in order to get all permutations after i-th element
		}else{
			res = append(res, get_all_k_elements_premises(k-1, k_element_clone, project)...)
		}
	}
	return res
}