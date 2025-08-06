package main

import (
	"strconv"
	"strings"
)

type Proposition struct {
	rule_name   string
	conclusion  string
	params      []string
	premises    []string
	line int
}

func create_proposition(rule_name string, concusion string, params []string, premises []string, line int, project *Project) Proposition {
	res := Proposition{
		rule_name:  rule_name,
		conclusion: concusion,
		params:     params,
		premises:   premises,
		line: line,
	}
	project.all_statements = append(project.all_statements, res)
	return res
}

func deep_copy_proposition(statement Proposition) Proposition{
	var new_statement Proposition
	new_statement.rule_name = statement.rule_name
	new_statement.conclusion = statement.conclusion
	new_statement.params = append(new_statement.params, statement.params...)
	new_statement.premises = append(new_statement.premises, statement.premises...)
	new_statement.line = statement.line
	return new_statement
}

func verify_proposition(statement Proposition, project *Project)bool{
	present_statement := deep_copy_proposition(statement)
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
	// However, if applied_rule.are_any_params is true, we do not check the number of params
	if len(applied_rule.params) != len(present_statement.params) && !applied_rule.are_any_params {
		message("derriving a statement, there must be as many parameters as there defined in the applied rule. Line "  + strconv.Itoa(present_statement.line), project)
		return false
	}

	applied_rule = substitude_rule_with_params(statement.params, applied_rule)
	if(!compare_conclusions(applied_rule, statement, project)){
		return false
	}
	if (!verify_premisses(applied_rule, statement, project)){
		return false
	}
	project.all_legal_expressions = append(project.all_legal_expressions, statement.conclusion)
	return true
}

func compare_conclusions(rule Rule, statement Proposition, project *Project)bool{
	for _, element:= range rule.conclusions{
		if (element == statement.conclusion){
			return true;
		}
	}
	msg_line := "conclusion "+ statement.conclusion + " does not correspond to any conclusion of the rule " + rule.name + ". Line " + strconv.Itoa(statement.line) + "\n See:"
	message(msg_line, project)
	for _, element := range rule.conclusions{
	message(element, project)
	}
	return false
}
func verify_premisses(rule Rule, statement Proposition, project *Project)bool{
	if(rule.are_any_premisses){
		return true
	}
	for index, req_premise := range rule.premises{
		if (req_premise!= statement.premises[index]){
			message("Premises passed in for verification of a proposition must satisfy the condtion required:", project)
			message("got \"" + statement.premises[index] + "\" where", project)
			message("\""+ req_premise + "\" expected", project)
			return false;
		}
	}
	for _, premise := range rule.premises{
		is_premise_verified := false
		for _, expression := range project.all_legal_expressions{
			if (premise == expression){
				is_premise_verified = true
			}
		}
		if (!is_premise_verified){
			message("not all premises are verified. Line " + strconv.Itoa(statement.line), project)
			return false
		}
		is_premise_verified = false
	}
	
	return true
}

// func gets a rule and a statement we have applied the rule on
// it returns another rule being a copy of the initiate rule which params contained in the conclusion and premises are exchanged with arguments given in the statement
func substitude_rule_with_params(params []string, rule Rule) Rule {
	substituted_rule := rule
	for i := 0; i<len(substituted_rule.params); i++{
		consequence := "[" + substituted_rule.params[i] + "]"
		// replacing params signs in premises of the rule (rule.params) with expressions in statements as arguments (statement.params)
		for j := 0; j < len(substituted_rule.premises); j++{
			substituted_rule.premises[j] = strings.Replace(substituted_rule.premises[j], consequence, params[i], -1)
		}
		for j := 0; j < len(substituted_rule.conclusions); j++{
			substituted_rule.conclusions[j] = strings.Replace(substituted_rule.conclusions[j], consequence, params[i], -1)
		}
	}

	return substituted_rule
}
