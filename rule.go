package main

type Rule struct {
	name        string
	premises    []string
	params      []string
	conclusions []string
	line        int
}

func create_rule(name string, params []string, premises []string, conclusions []string, line int, project *LC_project) Rule {
	res := Rule{
		name:        name,
		premises:    premises,
		params:      params,
		conclusions: conclusions,
		line:        line,
	}
	project.all_rules = append(project.all_rules, res)
	return res
}

func deep_copy_rule(old_rule Rule) Rule {
	var new_rule Rule
	new_rule.name = old_rule.name
	new_rule.line = old_rule.line
	new_rule.conclusions = append(new_rule.conclusions, old_rule.conclusions...)

	new_rule.params = append(new_rule.params, old_rule.params...)
	new_rule.premises = append(new_rule.premises, old_rule.premises...)

	return new_rule
}