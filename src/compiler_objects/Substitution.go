package compiler_objects

type Substitution struct {
	init        Param
	sub         Param
	Consditions []Formula
}

type Substitute struct {
	Proposition Proposition
	init        Formula
	sub         Formula
}

func CreateSubstitution(init Param, sub Param) Substitution {
	return Substitution{init: init, sub: sub}
}

func CreateSubstitute(prop Proposition, init Formula, sub Formula) Substitute {
	return Substitute{Proposition: prop, init: init, sub: sub}
}
