package lang_construct

import (
	"fmt"
	"strings"
)

type Statement struct {
	rule_name   string
	conclusion  string
	params      []string
	premises    []string
}

func verify_statement(statement Statement, project *LC_project) bool{
	var applied_rule Rule
	for i := 0; i < len(project.all_rules); i++{
		if project.all_rules[i].name == statement.rule_name{
			applied_rule = project.all_rules[i]
		}
	}
	if applied_rule.name == ""{
		fmt.Println("no rule " + statement.rule_name + " was found")
		return false
	}
	if len(applied_rule.params) != len(statement.params) {
		fmt.Println("derriving a statement, there must be as many parameters as there defined in the applied rule")
		return false
	}
	if len(applied_rule.premises) != len(statement.premises) {
		fmt.Println("derriving a statement, there must be as many premises as there defined in the applied rule")
		return false
	}
	if !check_rule_applicability(statement, applied_rule){
		fmt.Println("the rule is unapplicable")
		return false
	}
	if !are_premises_verified(statement.premises, *project){
		fmt.Println("not all premises are verified")
		return false
	}
	project.all_legal_expressions = append(project.all_legal_expressions, statement.conclusion)
	return true
}

func substitude_rule_with_params(statement Statement, rule Rule) Rule {
	substituted_rule := rule
	for i := 0; i<len(substituted_rule.params); i++{
		consequence := "[" + substituted_rule.params[i] + "]"
		for j := 0; j < len(substituted_rule.premises); j++{
			substituted_rule.premises[j] = strings.Replace(substituted_rule.premises[j], consequence, statement.params[i], -1)
		}
		substituted_rule.conclusion = strings.Replace(substituted_rule.conclusion, consequence, statement.params[i], -1)
	}

	return substituted_rule
}

func check_rule_applicability(statement Statement, rule Rule) bool{
	substituted_rule := substitude_rule_with_params(statement, rule)
	if substituted_rule.conclusion != statement.conclusion{
		fmt.Println("FUCK YOUR CONCLUSION!")
		return false
	}
	for i := 0; i < len(substituted_rule.premises); i++{
		if substituted_rule.premises[i] != statement.premises[i]{
			fmt.Println("FUCK YOUR PREMISES!")
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