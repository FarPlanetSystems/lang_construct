package lang_construct

import "fmt"
type Rule struct {
	name                string
	premises           []string
	params              []string
	conclusion          string
	message             func(string)
}

func default_message(msg_line string){
	fmt.Print(msg_line)
}