package main

import (
	"fmt"
	"os"
	"strings"
)

//since we can import other projects, we need to be sure that none of them point to the initiate one (we want to prevent cycles)
//for this purpose we create an array containing the initiate projects, and then add all imported projects, that import other projects
//if at the moment of importation the project we import is already contained in the array, then a cycle is present
var importing_projects []*Project

type Project struct {
	all_rules               []Rule
	all_statements          []Proposition
	all_definitions         []string
	all_legal_expressions   []string
	reports []string
	doc_code                string
	is_there_report_section bool
	is_interpreted_succesfully bool
	is_coherent bool
	project_file_path               string
	imported_projects_file_paths []string
}

func create_project(raw_text string, file_path string) *Project {
	
	res := Project{
		doc_code:                raw_text,
		is_there_report_section: false,
		project_file_path:               file_path,
		is_interpreted_succesfully: true,
		is_coherent: false,
	}
	return &res
}
// gets a string representing an id and a LC_project where we want to find it
// it checks all rules and compares the given id with their names
// if any matches are present, returns true; otherwise - false
func find_id_in_project(id string, project Project) bool{
	for i := 0; i < len(project.all_rules); i++{
		if id == project.all_rules[i].name{
			return true
		}
	}
	return false
}

func main() {
	run()	
}
func run(){
	// we get the file name of a .txt from the terminal
	file_path, err := get_file_path()
	// check that everything was alright with getting file name
	if err != nil{
		fmt.Println(err)
		return
	}
	// get the code inside the file
	code, err := read_code(file_path)
	// check that everything was alright opening and reading the file
	if err != nil{
		fmt.Println(err)
		return
	}
	// check whether the file is empty
	if len(code) < 1{
		fmt.Println("cannot run an empty file " + file_path)
		return
	}
	project := create_project(code, file_path)
	
	correct := interpretation_cycle(project)
	// if both interpretation and verification appear successful, we send a corresponding message
	if correct{
		message("Coherence verified!", project)
	}
	//we send all saved messages to the file
	report(*project)

}
// interprete and verify given project
// in case everything went fine returns true, otherwise false
func interpretation_cycle(project *Project) bool{

		//we let the lexer and parser do their work
		interpret_project(project)

		// we import all needed projects
		for i := 0; i < len(project.imported_projects_file_paths); i++{
			// get the code inside the file
			code, err := read_code(project.imported_projects_file_paths[i])
			// check that everything was alright opening and reading the file
			if err != nil{
				fmt.Println(err)
				return false
			}
			// check whether the file is empty
			if len(code) < 1{
				fmt.Println("cannot run an empty file " + project.imported_projects_file_paths[i])
				return false
			}
			imported_project := create_project(code, project.imported_projects_file_paths[i])
			
			if !import_project(imported_project, project){
				return false
			}
		}
		//if there is no errors in the code, we start to verify each given statement (lines with HAVE)
		if project.is_interpreted_succesfully{
			verify_all_included_statements(project)
		}
	return project.is_interpreted_succesfully && project.is_coherent
}


func interpret_project(project *Project) {
	parser := create_Parser(create_Lexer(project.doc_code, project), project)
	project.is_interpreted_succesfully = Language(parser)
	project.is_there_report_section = parser.is_there_report_section
}

// gets a string representing the path of a txt working file in "projects" folder. 
//if there is such a file, returns its content before "@" symbol converted to string and nil.
//otherwise, an empty string and an error
func read_code(file_path string) (string, error){
	_, err := os.Open(file_path)
	if err != nil{
		return "", err
	}
	bytes, _:= os.ReadFile(file_path)
	code := ""
	for i := 0; i < len(bytes); i++{
		if string(bytes[i]) != "@"{
		code += string(bytes[i])
	}else{
		return code, nil
	}
	}
	return code, nil
}

func get_file_path()(string, error){
	file_name := os.Args[1];
	// parsing
	file_name = strings.ReplaceAll(file_name, " ", "")
	file_name = strings.ReplaceAll(file_name, "\n", "")
	file_name = strings.ReplaceAll(file_name, "\r", "")
	// creating the path
	curdir, err := os.Getwd()
	file_path:= curdir + "\\" + file_name
	return file_path, err
}

func verify_all_included_statements(project *Project){
	for i := 0; i < len(project.all_statements); i++ {
		is_verified := verify_proposition(project.all_statements[i], project)
		if !is_verified {
			project.is_coherent = false
			return
		}
	}
	project.is_coherent = true
}


func create_definition(definition string, project *Project) {
	project.all_definitions = append(project.all_definitions, definition)
	project.all_legal_expressions = append(project.all_legal_expressions, definition)
}
// reads and interpret the code in file named project_file.
// if both reading and interpretation run successfully, ads all rules and legal expressions from read project to the one given in params and returns true
// otherwise returns false
func import_project(project_from *Project, project_to *Project) bool{
	importing_projects = append(importing_projects, project_to)
	
	// check importation cylcing
	for i := 0; i<len(importing_projects); i++{
		if project_from.project_file_path == (*importing_projects[i]).project_file_path{
			fmt.Println("imported projects cycle. See " + project_from.project_file_path)
			return false
		}
	}
	// we read and interpret the imported project as we do with the initiate one
	correct := interpretation_cycle(project_from)
	// we add all rules and legal expressions from one to another
	if correct{
		project_to.all_legal_expressions = append(project_to.all_legal_expressions, project_from.all_legal_expressions...)
		project_to.all_rules = append(project_to.all_rules, project_from.all_rules...)
		return true
	}else{
		fmt.Println("could not interpret project in file "+ project_from.project_file_path +". Run it separatly to fix all occured errors")
		return false
	}
}

func message(message_line string, project *Project) {
	message_line = message_line + "\n"
	project.reports = append(project.reports, message_line)
}

func report(project Project){
	os.OpenFile(project.project_file_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	new_doc_line := project.doc_code
	if !project.is_there_report_section {
		new_doc_line = project.doc_code + "@\n"
	}
	for i := 0; i <len(project.reports); i++{
		new_doc_line += project.reports[i]
		new_doc_line += "\n"
	}
	os.Truncate(project.project_file_path, 0)
	i := os.WriteFile(project.project_file_path, []byte(new_doc_line), 0644)
	if i != nil{
		fmt.Println(i)
	}
}