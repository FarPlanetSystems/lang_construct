package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


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
	file_name               string
}

func create_project(raw_text string, file_name string) *LC_project {
	curdir, err := os.Getwd()
	if err != nil{
		fmt.Println(err)
	}
	res := LC_project{
		doc_code:                raw_text,
		is_there_report_section: false,
		file_name:               curdir + "\\projects\\" + file_name,
		is_interpreted_succesfully: true,
		is_coherent: false,
	}
	return &res
}

func deep_copy_project(old_project LC_project) LC_project {
	var new_project LC_project
	new_project.doc_code = old_project.doc_code
	new_project.file_name = old_project.file_name
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

func main() {
	os.Mkdir("projects", 0644)
	for{
		// we get the file name of a .txt from the terminal
		file_name := get_file_name() 
		interpretation_cycle(file_name)
		
	}
}

func interpretation_cycle(file_name string){
	//we get the content of the file of the given name and convert it to string
	code, is_succesfull := read_code(file_name)
	if len(code) < 1{
		is_succesfull = false
		fmt.Println("cannot run an empty file")
	}
	if is_succesfull{
		project := create_project(code, file_name)
		//we let the lexer and parser do their work
		interpret_project(project)
		//if there is no errors in the code, we start to verify each given statement (lines with HAVE)
		if project.is_interpreted_succesfully{
			verify_all_included_statements(project)
		}
		// if both interpretation and verification appear successful, we send a corresponding message
		if project.is_interpreted_succesfully && project.is_coherent{
			message("Coherence verified!", project)
		}
		//we send all saved messages to the file
		report(*project)
	}

}


func interpret_project(project *LC_project) {
	parser := create_Parser(create_Lexer(project.doc_code, project), project)
	project.is_interpreted_succesfully = Language(parser)
	project.is_there_report_section = parser.is_there_report_section
}

// gets a string representing the name of a txt working file in "projects" folder. 
//if there is such a file, returns its content before "@" symbol converted to string and true.
//otherwise, an empty string and false
func read_code(file_name string) (string, bool){
	curdir, err := os.Getwd()
	if err != nil{
		fmt.Println(err)
	}
	file_path := curdir + "\\projects\\" + file_name
	_, err = os.Open(file_path)
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

func get_file_name()string{
	fmt.Println("please enter the name of the lang construct file in the projects folder: ")
	reader := bufio.NewReader(os.Stdin)
	file_name, _ := reader.ReadString('\n')
	file_name = strings.ReplaceAll(file_name, " ", "")
	file_name = strings.ReplaceAll(file_name, "\n", "")
	file_name = strings.ReplaceAll(file_name, "\r", "")
	return file_name
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
	os.OpenFile(project.file_name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	new_doc_line := project.doc_code
	if !project.is_there_report_section {
		new_doc_line = project.doc_code + "@\n"
	}
	for i := 0; i <len(project.reports); i++{
		new_doc_line += project.reports[i]
		new_doc_line += "\n"
	}
	os.Truncate(project.file_name, 0)
	i := os.WriteFile(project.file_name, []byte(new_doc_line), 0644)
	if i != nil{
		fmt.Println(i)
	}
}

func create_definition(definition string, project *LC_project) {
	project.all_definitions = append(project.all_definitions, definition)
	project.all_legal_expressions = append(project.all_legal_expressions, definition)
}