package main

import (
	"github.com/FarPlanetSystems/lang_construct/compiler_objects"
)

// since we can import other projects, we need to be sure that none of them point to the initiate one (we want to prevent cycles)
// for this purpose we create an array containing the initiate projects, and then add all imported projects, that import other projects
// if at the moment of importation the project we import is already contained in the array, then a cycle is present
var importing_projects []*Project

type Project struct {
	allSyntaxRules  []compiler_objects.SyntaxRule
	allPropositions []compiler_objects.Proposition
}

func createProject(syntax []compiler_objects.SyntaxRule, propositions []compiler_objects.Proposition) *Project {

	res := Project{
		allSyntaxRules:  syntax,
		allPropositions: propositions,
	}
	return &res
}

// gets a string representing an id and a LC_project where we want to find it
// it checks all rules and compares the given id with their names
// if any matches are present, returns true; otherwise - false
func (project Project) findIdInProject(id string) bool {
	for i := 0; i < len(project.allSyntaxRules); i++ {
		if id == project.allSyntaxRules[i].Name {
			return true
		}
	}
	return false
}

/*
	func main() {
		run()
	}
		func run() {
	// we get the file name of a .txt from the terminal
	file_path, err := getFilePath()
	// check that everything was alright with getting file name
	if err != nil {
		fmt.Println(err)
		return
	}
	// get the code inside the file
	code, err := readCode(file_path)
	// check that everything was alright opening and reading the file
	if err != nil {
		fmt.Println(err)
		return
	}
	// check whether the file is empty
	if len(code) < 1 {
		fmt.Println("cannot run an empty file " + file_path)
		return
	}
	project := createProject(code, file_path)

	correct := interpret(project)
	// if both interpretation and verification appear successful, we send a corresponding message
	if correct {
		project.messanger.insertMessage("Coherence verified!", -1)
	}
	//we send all saved messages to the file
	report(*project)

}
*/
/*
func interpret(project *Project) bool {

	//we let the lexer and parser do their work
	interpretProject(project)

	// we import all needed projects
	for i := 0; i < len(project.importedProjectsPaths); i++ {
		// get the code inside the file
		//code, err := readCode(project.importedProjectsPaths[i])
		code := ""
		// check whether the file is empty
		if len(code) < 1 {
			fmt.Println("cannot run an empty file " + project.importedProjectsPaths[i])
			return false
		}
		imported_project := createProject(code, project.importedProjectsPaths[i])

		if !importProject(imported_project, project) {
			return false
		}
	}
	//if there is no errors in the code, we start to verify each given statement (lines with HAVE)
	if project.is_interpreted_succesfully {
		project.isCoherent = project.verify()
	}
	return project.is_interpreted_succesfully && project.isCoherent
}

func interpretProject(project *Project) {
	parser := createParser(project.messanger)
	project.is_interpreted_succesfully = true //parser.Language(project)
	project.isThereReportSection = parser.isThereReportSection
}

// gets a string representing the path of a txt working file in "projects" folder.
// if there is such a file, returns its content before "@" symbol converted to string and nil.
// otherwise, an empty string and an error

func (project *Project) verify() bool {
	queue := &project.unverifiedExpressions
	result := true
	for len(queue.elements) > 0 {
		expression := queue.dequeue()
		isVerified := expression.verify(project)
		if isVerified {
			project.addToVerified(expression)
		} else {
			result = false
		}
	}
	return result
}

func (project *Project) addToVerified(element UnverifiedElement) {
	fmt.Println(element.unverifiedType)
	switch element.unverifiedType {
	case UNVERIFIED_PROPAREA:
		project.legalExpressions = append(project.legalExpressions, element.propArea.confirmedPropositions...)
		fmt.Println(element.propArea.confirmedPropositions)
	case UNVERIFIED_SPECIFICATION:
		{
			project.rules = append(project.rules, element.specification)
		}
	case UNVERIFIED_PROPOSITION:
		project.legalExpressions = append(project.legalExpressions, element.proposition.Conclusions...)
	}
}

func createDefinition(definition string, project *Project) {
	project.definitions = append(project.definitions, definition)
	project.legalExpressions = append(project.legalExpressions, definition)
}

// reads and interpret the code in file named project_file.
// if both reading and interpretation run successfully, ads all rules and legal expressions from read project to the one given in params and returns true
// otherwise returns false
func importProject(project_from *Project, project_to *Project) bool {
	importing_projects = append(importing_projects, project_to)

	// check importation cylcing
	for i := 0; i < len(importing_projects); i++ {
		if project_from.projectFilePath == (*importing_projects[i]).projectFilePath {
			fmt.Println("imported projects cycle. See " + project_from.projectFilePath)
			return false
		}
	}
	// we read and interpret the imported project as we do with the initiate one
	correct := interpret(project_from)
	// we add all rules and legal expressions from one to another
	if correct {
		project_to.legalExpressions = append(project_to.legalExpressions, project_from.legalExpressions...)
		project_to.rules = append(project_to.rules, project_from.rules...)
		return true
	} else {
		fmt.Println("could not interpret project in file " + project_from.projectFilePath + ". Run it separatly to fix all occured errors")
		return false
	}
}

func report(project Project) {
	os.OpenFile(project.projectFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	new_doc_line := project.code
	if !project.isThereReportSection {
		new_doc_line = project.code + "@\n"
	}
	for i := 0; i < len(project.messanger.messages); i++ {
		new_doc_line += project.messanger.messages[i]
		new_doc_line += "\n"
	}
	os.Truncate(project.projectFilePath, 0)
	i := os.WriteFile(project.projectFilePath, []byte(new_doc_line), 0644)
	if i != nil {
		fmt.Println(i)
	}
}
*/
