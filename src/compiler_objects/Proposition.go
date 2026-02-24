package compiler_objects

type Proposition struct {
	Conclusion Statement
	Line       int
}

type Proposition2 struct {
	ruleName    string
	Conclusions []string
	params      []string
	premises    []string
	Line        int
}

func createProposition(rule_name string, concusions []string, params []string, premises []string, line int) Proposition2 {
	res := Proposition2{
		ruleName:    rule_name,
		Conclusions: concusions,
		params:      params,
		premises:    premises,
		Line:        line,
	}
	return res
}

func deepCopyProposition(statement Proposition2) Proposition2 {
	var new_statement Proposition2
	new_statement.ruleName = statement.ruleName
	new_statement.Conclusions = statement.Conclusions
	new_statement.params = append(new_statement.params, statement.params...)
	new_statement.premises = append(new_statement.premises, statement.premises...)
	new_statement.Line = statement.Line
	return new_statement
}

/*
func (proposition Proposition2) verify_proposition(project *Project, containing_area *PropArea) bool {
	present_statement := deepCopyProposition(proposition)
	//we look for a rule in project.all_rules reference to which must be contained in statement.rule_name
	var applied_rule Rule
	for i := 0; i < len(project.rules); i++ {
		if project.rules[i].name == present_statement.ruleName {
			applied_rule = deep_copy_rule(project.rules[i])
		}
	}

	// applied rule.name being empty indicates that there is no such rule in project.all_rules. In this case we message an error and return false
	if applied_rule.name == "" {
		project.messanger.InsertMessage("Error: no rule "+present_statement.ruleName+" was found.", present_statement.Line)
		return false
	}

	// if the number of params in applied_rule is not equal to that in present_statement, we message an error and return false
	// However, if applied_rule.are_any_params is true, we do not check the number of params
	if len(applied_rule.params) != len(present_statement.params) && !applied_rule.are_any_params {
		project.messanger.InsertMessage("Error: derriving a statement, there must be as many parameters as there defined in the applied rule.", present_statement.Line)
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

func compareConclusions(rule Rule, statement Proposition2, project *Project) bool {
	for _, conclusion := range statement.Conclusions {
		if !slices.Contains(rule.conclusions, conclusion) {
			msg_line := "Error: conclusion " + conclusion + " does not correspond to any conclusion of the rule " + rule.name + "\n	See:"
			project.messanger.InsertMessage(msg_line, statement.Line)
			for _, element := range rule.conclusions {
				project.messanger.InsertMessage("	"+element, -1)
			}
			return false
		}
	}

	return true
}
func verifyPremisses(rule Rule, statement Proposition2, project *Project, containing_area *PropArea) bool {
	if rule.are_any_premisses {
		return true
	}
	for index, statement_premise := range statement.premises {
		if rule.premises[index] != statement_premise {
			project.messanger.InsertMessage("Error: Premises passed in for verification of a proposition must satisfy the condtion required.", statement.Line)
			project.messanger.InsertMessage("	got \""+statement_premise+"\" where", -1)
			project.messanger.InsertMessage("	\""+rule.premises[index]+"\" expected", -1)
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
			project.messanger.InsertMessage("Error: not all premises are verified. See: ", -1)
			project.messanger.InsertMessage(premise, statement.Line)
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
*/
