package compiler_objects

type Rule struct {
	Name       string
	Premises   []Formula
	Params     []Param
	Conclusion Formula
	Line       int
}

type Param struct {
	Id          string
	GrammarType string
}

type VariableSubstitution struct {
	Variable Param
	Value    Formula
}

func AddToSubstitutions(substitutions []VariableSubstitution, substitution VariableSubstitution) (bool, []VariableSubstitution) {
	//fmt.Println("adding to substitution")
	substitutions = append(substitutions, substitution)
	for _, sub := range substitutions {
		if sub.Variable.Id == substitution.Variable.Id && !substitution.Value.Compare(sub.Value) {
			//fmt.Println(sub.Variable.Id)
			//PrintStatement(substitution.Value)
			//fmt.Println("was found")
			PrintStatement(substitution.Value)
			//return false, substitutions
		}
	}
	return true, substitutions
}

func AreSubstitutionsConsistent(substitutions []VariableSubstitution) bool {
	for _, sub1 := range substitutions {
		for _, sub2 := range substitutions {
			if sub1.Variable.Id == sub2.Variable.Id && !sub1.Value.Compare(sub2.Value) {
				return false
			}
		}
	}
	return true
}

func (rule Rule) IsParamId(id string) bool {
	for _, param := range rule.Params {
		if param.Id == id {
			return true
		}
	}
	return false
}

func (rule Rule) GetParamGrammarType(id string) string {
	for _, param := range rule.Params {
		if param.Id == id {
			return param.GrammarType
		}
	}
	return ""
}

func CreateRule(name string, params []Param, premises []Formula, conclusion Formula, line int) Rule {
	res := Rule{
		Name:       name,
		Premises:   premises,
		Params:     params,
		Conclusion: conclusion,
		Line:       line,
	}
	return res
}

func Deep_copy_rule(old_rule Rule) Rule {
	var new_rule Rule
	new_rule.Name = old_rule.Name
	new_rule.Line = old_rule.Line
	new_rule.Conclusion = old_rule.Conclusion

	new_rule.Params = old_rule.Params
	new_rule.Premises = old_rule.Premises

	return new_rule
}
