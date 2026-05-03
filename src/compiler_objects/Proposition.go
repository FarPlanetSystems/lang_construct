package compiler_objects

type Proposition struct {
	Formula Formula
	Line    int
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
