package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//since we can import other projects, we need to be sure that none of them point to the initiate one (we want to prevent cycles)
//for this purpose we ...
var importing_projects []string

type LC_project struct {
	all_rules               []Rule
	all_statements          []Statement
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

func create_project(raw_text string, file_path string) *LC_project {
	res := LC_project{
		doc_code:                raw_text,
		is_there_report_section: false,
		project_file_path:               file_path,
		is_interpreted_succesfully: true,
		is_coherent: false,
	}
	return &res
}

func deep_copy_project(old_project LC_project) LC_project {
	var new_project LC_project
	new_project.doc_code = old_project.doc_code
	new_project.project_file_path = old_project.project_file_path
	new_project.is_coherent = old_project.is_coherent
	new_project.is_interpreted_succesfully = old_project.is_interpreted_succesfully
	new_project.is_there_report_section = old_project.is_there_report_section

	new_project.all_definitions = append(new_project.all_definitions, old_project.all_definitions...)
	new_project.all_legal_expressions = append(new_project.all_legal_expressions, old_project.all_legal_expressions...)
	new_project.all_rules = append(new_project.all_rules, old_project.all_rules...)
	new_project.all_statements = append(new_project.all_statements, old_project.all_statements...)
	new_project.reports = append(new_project.reports, old_project.reports...)
	return new_project
}
// gets a string representing an id and a LC_project where we want to find it
// it checks all rules and compares the given id with their names
// if any matches are present, returns true; otherwise - false
func find_id_in_project(id string, project LC_project) bool{
	for i := 0; i < len(project.all_rules); i++{
		if id == project.all_rules[i].name{
			return true
		}
	}
	return false
}

func main() {
	for{
		// we get the file name of a .txt from the terminal
		file_path := get_file_path()


		project := interpretation_cycle(file_path)
		// if both interpretation and verification appear successful, we send a corresponding message
		if project.is_interpreted_succesfully && project.is_coherent{
			message("Coherence verified!", project)
		}
		//we send all saved messages to the file
		report(*project)
	}
}

func interpretation_cycle(file_name string) *LC_project{
	//we get the content of the file of the given name and convert it to string
	var project *LC_project
	
	code, is_succesfull := read_code(file_name)
	if len(code) < 1{
		is_succesfull = false
		fmt.Println("cannot run an empty file " + file_name)
	}
	if is_succesfull{
		project = create_project(code, file_name)
		//we let the lexer and parser do their work
		interpret_project(project)
		// we import all needed projects
		for i := 0; i < len(project.imported_projects_file_paths); i++{
			if !import_project_to(project.imported_projects_file_paths[i], project){
				var empty_project *LC_project = create_project("", "")
				return empty_project
			}
		}
		//if there is no errors in the code, we start to verify each given statement (lines with HAVE)
		if project.is_interpreted_succesfully{
			verify_all_included_statements(project)
		}
		
	}
	return project
}


func interpret_project(project *LC_project) {
	parser := create_Parser(create_Lexer(project.doc_code, project), project)
	project.is_interpreted_succesfully = Language(parser)
	project.is_there_report_section = parser.is_there_report_section
}

// gets a string representing the path of a txt working file in "projects" folder. 
//if there is such a file, returns its content before "@" symbol converted to string and true.
//otherwise, an empty string and false
func read_code(file_path string) (string, bool){
	_, err := os.Open(file_path)
	if err != nil{
		fmt.Println(err)
		return "", false
	}
	bytes, _:= os.ReadFile(file_path)
	code := ""
	for i := 0; i < len(bytes); i++{
		if string(bytes[i]) != "@"{
		code += string(bytes[i])
	}else{
		return code, true
	}
	}
	return code, true
}

func get_file_path()string{
	// get file name
	fmt.Println("please enter the name of the lang_construct file in the current directory: ")
	reader := bufio.NewReader(os.Stdin)
	file_name, _ := reader.ReadString('\n')
	// parsing
	file_name = strings.ReplaceAll(file_name, " ", "")
	file_name = strings.ReplaceAll(file_name, "\n", "")
	file_name = strings.ReplaceAll(file_name, "\r", "")
	// creating the path
	curdir, err := os.Getwd()
	if err != nil{
		fmt.Println(err)
	}
	file_path:= curdir + "\\" + file_name
	return file_path
}

func verify_all_included_statements(project *LC_project){
	for i := 0; i < len(project.all_statements); i++ {
		is_verified := verify_statement(project.all_statements[i], project)
		if !is_verified {
			project.is_coherent = false
			return
		}
	}
	project.is_coherent = true
}

func message(message_line string, project *LC_project) {
	message_line = message_line + "\n"
	project.reports = append(project.reports, message_line)
}

func report(project LC_project){
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

func create_definition(definition string, project *LC_project) {
	project.all_definitions = append(project.all_definitions, definition)
	project.all_legal_expressions = append(project.all_legal_expressions, definition)
}
// reads and interpret the code in file named project_file.
// if both reading and interpretation run successfully, ads all rules and legal expressions from read project to the one given in params and returns true
// otherwise returns false
func import_project_to(project_file string, project *LC_project) bool{
	importing_projects = append(importing_projects, project.project_file_path)
	
	for i := 0; i<len(importing_projects); i++{
		if project_file == importing_projects[i]{
			fmt.Println("imported projects cycle. See " + project_file)
			return false
		}
	}
	// we read and interpret the imported project as we do with the initiate one
	new_project := interpretation_cycle(project_file)
	// we add all rules and legal expressions from one to another
	if new_project.is_interpreted_succesfully && new_project.is_coherent{
		project.all_legal_expressions = append(project.all_legal_expressions, new_project.all_legal_expressions...)
		project.all_rules = append(project.all_rules, new_project.all_rules...)
		return true
	}else{
		fmt.Println("could not interpret project in file "+ project_file+". Run it separatly to fix all occured errors")
		return false
	}
}