package main

import "fmt"

type PropArea struct {
	param                 string
	condition             string
	containedPropositions []Proposition
	confirmedPropositions []string
}

func (area *PropArea) verify(project *Project) bool {
	is_verified := false
	for _, proposition := range area.containedPropositions {
		is_verified = proposition.verify_proposition(project, area)
		if is_verified {
			area.confirmedPropositions = append(area.confirmedPropositions, proposition.conclusion)
		}
	}
	return is_verified

}

func findSpecifiedProposition(project Project, rule Rule) string {
	// in the future we will allow our specifications to have more than one single parameter and condition
	if len(rule.premises) != 1 {
		message := "specification must contain only one premise"
		return message
	}
	if len(rule.params) != 1 {
		message := "specification must contain only one parameter"
		return message
	}
	if len(rule.conclusions) != 1 {
		message := "specification must contain only one conclusion"
		return message
	}
	for _, area := range project.propositionalAreas {
		if area.condition == rule.premises[0] && area.param == rule.params[0] {
			fmt.Println("rule conclusion")
			fmt.Println(rule.conclusions[0])
			fmt.Println("area conclusions")
			for _, proposition := range area.confirmedPropositions {
				fmt.Println(proposition)
				if proposition == rule.conclusions[0] {
					return "success"
				}
			}
		}

	}
	message := "no specified proposition with the given parameter and condition was found among propositional areas"
	return message
}
