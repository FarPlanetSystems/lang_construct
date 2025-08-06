package main

import "strings"

const PROPOSITIONAL_ARGUMENT_TYPE = "propositional"
const RULE_ARGUMENT_TYPE = "rule"

type Argument struct {
	argument_type       string
	propositional_value string
	rule_value          Rule
}

func createArgument(argumentType string, propositionalValue string, ruleValue Rule) Argument {
	res := Argument{
		argument_type:       argumentType,
		propositional_value: propositionalValue,
		rule_value:          ruleValue,
	}
	return res
}

func compareArguments(arg1 Argument, arg2 Argument) bool {
	if arg1.argument_type != arg2.argument_type {
		return false
	}
	if arg1.argument_type == PROPOSITIONAL_ARGUMENT_TYPE {
		if arg1.propositional_value != arg2.propositional_value {
			return false
		}
	}
	if arg1.argument_type == RULE_ARGUMENT_TYPE {
		if !compareRule(arg1.rule_value, arg2.rule_value) {
			return false
		}
	}
	return true
}

func (arg Argument) Replace(old string, new string) Argument{
	var res Argument
	if arg.argument_type == PROPOSITIONAL_ARGUMENT_TYPE{
		newProposition := strings.Replace(arg.propositional_value, old, new, -1)
		res = createArgument(PROPOSITIONAL_ARGUMENT_TYPE, newProposition, Rule{})
	}else{
		var newPremises []Argument
		var newConclusions []Argument
		for _, premise := range arg.rule_value.premises{
			newPremises = append(newPremises, premise.Replace(old, new))
		}
		for _, conclusion := range arg.rule_value.conclusions{
			newConclusions = append(newConclusions, conclusion.Replace(old, new))
		}
		argRule := createRule(arg.rule_value.name, arg.rule_value.params, newPremises, newConclusions, arg.rule_value.line, arg.rule_value.are_any_params, arg.rule_value.are_any_premisses)
		res = createArgument(RULE_ARGUMENT_TYPE, "", argRule)
	}
	return res
}

func (arg Argument) ToString() string{
	if arg.argument_type == PROPOSITIONAL_ARGUMENT_TYPE{
		return arg.propositional_value
	}else{
		return arg.rule_value.ToString()
	}
}