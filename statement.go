package main

import (
	"strconv"
	"strings"
)

type Statement struct {
	rule_name   string
	conclusion  string
	params      []string
	premises    []string
	line int
}

func create_statement(rule_name string, concusion string, params []string, premises []string, line int, project *LC_project) Statement {
	res := Statement{
		rule_name:  rule_name,
		conclusion: concusion,
		params:     params,
		premises:   premises,
		line: line,
	}
	project.all_statements = append(project.all_statements, res)
	return res
}

func deep_copy_statement(statement Statement) Statement{
	var new_statement Statement
	new_statement.rule_name = statement.rule_name
	new_statement.conclusion = statement.conclusion
	new_statement.params = append(new_statement.params, statement.params...)
	new_statement.premises = append(new_statement.premises, statement.premises...)
	new_statement.line = statement.line
	return new_statement
}
// func gets a statement which coherence we want to check in project
// it returns true if it is coherent, otherwise - false
func verify_statement(statement Statement, project *LC_project) bool{
	present_statement := deep_copy_statement(statement)
	//we look for a rule in project.all_rules reference to which must be contained in statement.rule_name
	var applied_rule Rule
	for i := 0; i < len(project.all_rules); i++{
		if project.all_rules[i].name == present_statement.rule_name{
			applied_rule = deep_copy_rule(project.all_rules[i])
		}
	}
	// applied rule.name being empty indicates that there is no such rule in project.all_rules. In this case we message an error and return false
	if applied_rule.name == ""{
		message("no rule " + present_statement.rule_name + " was found. Line " + strconv.Itoa(present_statement.line), project)
		return false
	}
	// if the number of params in applied_rule is not equal to that in present_statement, we message an error and return false
	if len(applied_rule.params) != len(present_statement.params) {
		message("derriving a statement, there must be as many parameters as there defined in the applied rule. Line "  + strconv.Itoa(present_statement.line), project)
		return false
	}
	// if the number of premises given in applied_rule is not equal to that in present_statement, we try to find any sound premises that could replace the initiate ones
	if len(applied_rule.premises) != len(present_statement.premises) {
		
		sound_premises := find_sound_premises(applied_rule, present_statement, project)
		if len(sound_premises) > 0{
			present_statement.premises = append(present_statement.premises, sound_premises...)
		}else{
		// if there no such premises we message an error and return false
		message("derriving a statement, there must be as many premises as there defined in the applied rule. Line "  + strconv.Itoa(present_statement.line), project)
		return false
	}
	}
	if !check_rule_applicability(present_statement, applied_rule, project){
		return false
	}
	if !are_premises_verified(present_statement.premises, *project){
		message("not all premises are verified. Line " + strconv.Itoa(present_statement.line), project)
		return false
	}
	project.all_legal_expressions = append(project.all_legal_expressions, present_statement.conclusion)
	return true
}
// func gets a rule and a statement we have applied the rule on
// it returns another rule being a copy of the initiate rule which params contained in the conclusion and premises are exchanged with arguments given in the statement
func substitude_rule_with_params(statement Statement, rule Rule) Rule {
	substituted_rule := rule
	for i := 0; i<len(substituted_rule.params); i++{
		consequence := "[" + substituted_rule.params[i] + "]"
		// replacing params signs in premises of the rule (rule.params) with expressions in statements as arguments (statement.params)
		for j := 0; j < len(substituted_rule.premises); j++{
			substituted_rule.premises[j] = strings.Replace(substituted_rule.premises[j], consequence, statement.params[i], -1)
		}
		for j := 0; j < len(substituted_rule.conclusions); j++{
			substituted_rule.conclusions[j] = strings.Replace(substituted_rule.conclusions[j], consequence, statement.params[i], -1)
		}
	}

	return substituted_rule
}
// func checks if a given statement can be correctly infered from a given statement considering all legal expressions from the current project
// if it can, returns true. Otherwise - false
func check_rule_applicability(statement Statement, rule Rule, project *LC_project) bool{
	substituted_rule := substitude_rule_with_params(statement, rule)
	// checking if there is correspondece with the statement's conclusion with one of the rule's conclusion
	correspondece_found := false

	for i := 0; i < len(substituted_rule.conclusions); i++{

		if substituted_rule.conclusions[i] == statement.conclusion{
			correspondece_found = true
		}
	}
	if !correspondece_found{
		msg_line := "conclusion "+ statement.conclusion + " does not correspond to any conclusion of the rule " + substituted_rule.name + ". Line " + strconv.Itoa(statement.line) + "\n See:"
		message(msg_line, project)
		for i := 0; i<len(substituted_rule.conclusions); i++{
			message(substituted_rule.conclusions[i], project)
		}
		return false
	}
	// checking the correspondence among premises
	for i := 0; i < len(substituted_rule.premises); i++{
		if substituted_rule.premises[i] != statement.premises[i]{
			msg_line := "a premise "+ statement.premises[i] + " does not correspond to the required one " + substituted_rule.premises[i] + ". Line " + strconv.Itoa(statement.line) + "\n See:"
			message(msg_line, project)
			message(substituted_rule.premises[i] + " was expected, but " + statement.premises[i] + " was found", project)
			return false
		}
	}
	return true
}

func are_premises_verified(premises []string, project LC_project) bool{
	for i := 0; i < len(premises); i++{
		is_premise_found := false
		for j:=0; j < len(project.all_legal_expressions); j++{
			if project.all_legal_expressions[j] == premises[i]{
				is_premise_found = true
				break
			}
		}
		if !is_premise_found{
			return false
		}
	}
	return true
}
// func gets a rule we want to apply in order to verify an expression and the expression, which must contain no premises, and the current project
// it looks for a set of premises in project.all_legal_expressions which complete the given statement to a one that can be verified with the rule
// if there is such a set of legal expressions, it returns an array of strings representing it
// otherwise returns an empty array
func find_sound_premises(rule Rule, statement Statement, project *LC_project)[]string{
	sound_statement := deep_copy_statement(statement)

	if len(rule.premises) != 0 && len(statement.premises) == 0{
		possible_premises := get_all_k_elements_premises(len(rule.premises), []string{}, project)
		for i := 0; i < len(possible_premises); i++{
			match := true
			sound_statement.premises = possible_premises[i]
			substituted_rule := substitude_rule_with_params(sound_statement, rule)
				for i := 0; i < len(substituted_rule.premises); i++{
				if substituted_rule.premises[i] != sound_statement.premises[i]{
					match =  false
				}
			}
			if match == true{
				return possible_premises[i]
			}
		}
	}
	return []string{}
}
// recursive func gets an integer k represnting the size of permutation array, k_element, that must be an empty string array, and the current project
// returns an array of arrays res representing all possible k-sized permutations of all_legal_expressions array in project
func get_all_k_elements_premises(k int, k_element []string, project *LC_project)[][]string{
	res:= [][]string{}
	for i := 0; i<len(project.all_legal_expressions); i++{
		// k_element contains one single permutation
		// in each iteration we copy k_element in order to add the following element in each step of recursion
		k_element_clone := []string{}
		copy(k_element_clone, k_element)
		k_element_clone = append(k_element, project.all_legal_expressions[i])
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