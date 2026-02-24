package compiler_objects

const UNVERIFIED_PROPOSITION = "UNVERIFIED_PROPOSITION"
const UNVERIFIED_SPECIFICATION = "UNVERIFIED_SPECIFICATION"
const UNVERIFIED_PROPAREA = "UNVERIFIED_PROPAREA"

type UnverifiedElement struct {
	proposition    Proposition2
	propArea       *PropArea
	specification  Rule
	unverifiedType string
}

func createUnverifiedProposition(proposition Proposition2) UnverifiedElement {
	return UnverifiedElement{
		proposition:    proposition,
		propArea:       &PropArea{},
		specification:  Rule{},
		unverifiedType: UNVERIFIED_PROPOSITION,
	}
}

func createUnverifiedSpecification(spec Rule) UnverifiedElement {
	return UnverifiedElement{
		proposition:    Proposition2{},
		propArea:       &PropArea{},
		specification:  spec,
		unverifiedType: UNVERIFIED_SPECIFICATION,
	}
}

func createUnverifiedPropArea(propArea *PropArea) UnverifiedElement {
	return UnverifiedElement{
		proposition:    Proposition2{},
		propArea:       propArea,
		specification:  Rule{},
		unverifiedType: UNVERIFIED_PROPAREA,
	}
}

/*
func (element UnverifiedElement) verify(project *Project) bool {
	switch element.unverifiedType {
	case UNVERIFIED_PROPAREA:
		{
			result := element.propArea.verify(project)
			return result

		}
	case UNVERIFIED_PROPOSITION:
		{
			result := element.proposition.verify_proposition(project, nil)
			return result
		}
	case UNVERIFIED_SPECIFICATION:
		{
			result := findSpecifiedProposition(*project, element.specification)
			if result == "success" {
				return true
			}
			project.messanger.InsertMessage(result, element.specification.line)
			return false
		}
	default:
		return false
	}

}
*/
