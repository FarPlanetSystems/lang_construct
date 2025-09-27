package main

import (
	"slices"
	"strings"
)

type Proposition struct {
	rule_name   string
	conclusions []string
	params      []string
	premises    []string
	line        int
}

func createProposition(rule_name string, concusions []string, params []string, premises []string, line int, project *Project) Proposition {
	res := Proposition{
		rule_name:   rule_name,
		conclusions: concusions,
		params:      params,
		premises:    premises,
		line:        line,
	}
	project.statements = append(project.statements, res)
	return res
}

func deepCopyProposition(statement Proposition) Proposition {
	var new_statement Proposition
	new_statement.rule_name = statement.rule_name
	new_statement.conclusions = statement.conclusions
	new_statement.params = append(new_statement.params, statement.params...)
	new_statement.premises = append(new_statement.premises, statement.premises...)
	new_statement.line = statement.line
	return new_statement
}

func (proposition Proposition) verify_proposition(project *Project, containing_area *PropArea) bool {
	present_statement := deepCopyProposition(proposition)
	//we look for a rule in project.all_rules reference to which must be contained in statement.rule_name
	var applied_rule Rule
	for i := 0; i < len(project.rules); i++ {
		if project.rules[i].name == present_statement.rule_name {
			applied_rule = deep_copy_rule(project.rules[i])
		}
	}

	// applied rule.name being empty indicates that there is no such rule in project.all_rules. In this case we message an error and return false
	if applied_rule.name == "" {
		project.messanger.message("Error: no rule "+present_statement.rule_name+" was found.", present_statement.line)
		return false
	}

	// if the number of params in applied_rule is not equal to that in present_statement, we message an error and return false
	// However, if applied_rule.are_any_params is true, we do not check the number of params
	if len(applied_rule.params) != len(present_statement.params) && !applied_rule.are_any_params {
		project.messanger.message("Error: derriving a statement, there must be as many parameters as there defined in the applied rule.", present_statement.line)
		return false
	}

	applied_rule = substitudeRuleWithParams(proposition.params, applied_rule)
	if !compareConclusions(applied_rule, proposition, project) {
		return false
	}
	if !verifyPremisses(applied_rule, proposition, project, containing_area) {
		return false
	}
	return true
}

func compareConclusions(rule Rule, statement Proposition, project *Project) bool {
	for _, conclusion := range statement.conclusions {
		if !slices.Contains(rule.conclusions, conclusion) {
			msg_line := "Error: conclusion " + conclusion + " does not correspond to any conclusion of the rule " + rule.name + "\n	See:"
			project.messanger.message(msg_line, statement.line)
			for _, element := range rule.conclusions {
				project.messanger.message("	"+element, -1)
			}
			return false
		}
	}

	return true
}
func verifyPremisses(rule Rule, statement Proposition, project *Project, containing_area *PropArea) bool {
	if rule.are_any_premisses {
		return true
	}
	for index, statement_premise := range statement.premises {
		if rule.premises[index] != statement_premise {
			project.messanger.message("Error: Premises passed in for verification of a proposition must satisfy the condtion required.", statement.line)
			project.messanger.message("	got \""+statement_premise+"\" where", -1)
			project.messanger.message("	\""+rule.premises[index]+"\" expected", -1)
			return false
		}
	}
	for _, premise := range rule.premises {
		is_premise_verified := false
		if slices.Contains(project.legalExpressions, premise) {
			is_premise_verified = true
		}
		if containing_area != nil && slices.Contains(containing_area.confirmedPropositions, premise) {

			is_premise_verified = true
		}
		if containing_area != nil && containing_area.condition == premise {
			is_premise_verified = true
		}
		if !is_premise_verified {
			project.messanger.message("Error: not all premises are verified. See: ", -1)
			project.messanger.message(premise, statement.line)
			return false
		}
		is_premise_verified = false
	}

	return true
}

// func gets a rule and a statement we have applied the rule on
// it returns another rule being a copy of the initiate rule which params contained in the conclusion and premises are exchanged with arguments given in the statement
func substitudeRuleWithParams(params []string, rule Rule) Rule {
	substituted_rule := rule
	for i := 0; i < len(substituted_rule.params); i++ {
		consequence := "[" + substituted_rule.params[i] + "]"
		// replacing params signs in premises of the rule (rule.params) with expressions in statements as arguments (statement.params)
		for j := 0; j < len(substituted_rule.premises); j++ {
			substituted_rule.premises[j] = strings.Replace(substituted_rule.premises[j], consequence, params[i], -1)
		}
		for j := 0; j < len(substituted_rule.conclusions); j++ {
			substituted_rule.conclusions[j] = strings.Replace(substituted_rule.conclusions[j], consequence, params[i], -1)
		}
	}
	// if we have %any params, we replace all occurances of [%any] string inside of each rule premise and conclusion with the given params strings
	if substituted_rule.are_any_params {
		new_rule_premises := []string{}
		new_rule_conclusions := []string{}
		for _, param := range params {
			for _, premise := range rule.premises {
				new_premise := strings.Replace(premise, "[$any]", param, -1)
				new_rule_premises = append(new_rule_premises, new_premise)
			}
			for _, conclusion := range substituted_rule.conclusions {
				new_conclusion := strings.Replace(conclusion, "[$any]", param, -1)
				new_rule_conclusions = append(new_rule_conclusions, new_conclusion)
			}
		}
		substituted_rule.premises = new_rule_premises
		substituted_rule.conclusions = new_rule_conclusions
	}

	return substituted_rule
}
