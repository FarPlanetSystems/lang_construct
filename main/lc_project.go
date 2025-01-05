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

func Create_project(raw_text string, file_name string) *LC_project {
	res := LC_project{
		doc_code:                raw_text,
		is_there_report_section: false,
		file_name:               "projects\\" + file_name,
		is_interpreted_succesfully: true,
		is_coherent: false,
	}
	return &res
}

func run() {
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
	project := Create_project(code, file_name)
	if is_succesfull{
		project.doc_code = code
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
	_, err := os.Open("projects/"+file_name)
	if err != nil{
		fmt.Println(err)
		return "", false
	}
	bytes, _:= os.ReadFile("projects/"+file_name)
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
	curdir, _ := os.Getwd()
	os.Truncate(project.file_name, 0)
	i := os.WriteFile(curdir +"\\" + project.file_name, []byte(new_doc_line), 0644)
	if i != nil{
		fmt.Println(i)
	}
}

func create_definition(definition string, project *LC_project) {
	project.all_definitions = append(project.all_definitions, definition)
	project.all_legal_expressions = append(project.all_legal_expressions, definition)
}

func create_rule(name string, params []string, premises []string, conclusion string, project *LC_project) Rule {
	res := Rule{
		name:       name,
		premises:   premises,
		params:     params,
		conclusion: conclusion,
		message:    default_message,
	}
	project.all_rules = append(project.all_rules, res)
	return res
}

func create_statement(rule_name string, concusion string, params []string, premises []string, project *LC_project) Statement {
	res := Statement{
		rule_name:  rule_name,
		conclusion: concusion,
		params:     params,
		premises:   premises,
	}
	project.all_statements = append(project.all_statements, res)
	return res
}
/*
func reset_project(project *LC_project) {
	clear(project.all_rules)
	clear(project.all_statements)
	clear(project.all_definitions)
	clear(project.all_legal_expressions)
	clear(project.reports)
	project.is_coherent = false
	project.is_interpreted_succesfully = true
	project.is_there_report_section = false
}
	*/


