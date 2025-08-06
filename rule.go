package main

type Rule struct {
	name        string
	premises    []Argument
	params      []string
	conclusions []Argument
	are_any_premisses bool
	are_any_params bool
	line        int
}

func createRule(name string, params []string, premises []Argument, conclusions []Argument, line int, any_params bool, any_premisses bool) Rule {
	res := Rule{
		name:        name,
		premises:    premises,
		params:      params,
		conclusions: conclusions,
		line:        line,
		are_any_premisses: any_premisses,
		are_any_params: any_params,
	}
	return res
}

func deepCopyRule(oldRule Rule) Rule {
	var newRule Rule
	newRule.name = oldRule.name
	newRule.line = oldRule.line
	newRule.conclusions = append(newRule.conclusions, oldRule.conclusions...)

	newRule.params = append(newRule.params, oldRule.params...)
	newRule.premises = append(newRule.premises, oldRule.premises...)

	newRule.are_any_premisses = oldRule.are_any_premisses
	newRule.are_any_params = oldRule.are_any_params

	return newRule
}

func compareRule(rule1 Rule, rule2 Rule) bool {
	if !rule1.are_any_params && !rule2.are_any_params{
		if len(rule1.params) != len(rule2.params){
			return false
		}
	}
	
	if !rule1.are_any_premisses && !rule2.are_any_premisses{
		for i:=0; i< len(rule1.premises); i++{
			var match bool = false
			for j:=0; j<len(rule2.premises); j++{
				if compareArguments(rule1.premises[i], rule2.premises[j]){
					match = true
					break
				}
			}
			if !match {return false}
		} 
	}

	for i:=0; i< len(rule1.conclusions); i++{
			var match bool = false
			for j:=0; j<len(rule2.conclusions); j++{
				if compareArguments(rule1.conclusions[i], rule2.conclusions[j]){
					match = true
					break
				}
			}
			if !match {return false}
		} 
	return true
}


func (rule Rule) ToString() string {
	res := rule.name
	res += convertRuleParamsToString(rule)
	res += " :"
	res += convertRulePremisesToString(rule)
	res += " -> "
	res += convertRuleConclusionToString(rule)
	return res
}

func convertRuleParamsToString(rule Rule) string{
	res := "("
	if len(rule.params) > 0{
	for i, param := range rule.params {
			res += param
			if i < len(rule.params)-1 {
				res += ", "
			}
		}
	}
	res += ")"
	return res
}

func convertRulePremisesToString(rule Rule) string{
	res := ""
	if len(rule.premises) > 0 {
		for i, premise := range rule.premises {
			if premise.argument_type == PROPOSITIONAL_ARGUMENT_TYPE{
				res += premise.propositional_value
			}else{
				res += premise.rule_value.ToString()
			}
			if i < len(rule.premises)-1 {
				res += ", "
			}
		}
	}
	return res
}

func convertRuleConclusionToString(rule Rule) string{
	res := ""
	if len(rule.conclusions) > 0 {
		for i, conclusion := range rule.conclusions {
			if conclusion.argument_type == PROPOSITIONAL_ARGUMENT_TYPE{
				res += conclusion.propositional_value
			}else{
				res += conclusion.rule_value.ToString()
			}
			if i < len(rule.conclusions)-1 {
				res += ", "
			}
		}
	}
	return res
}