package main

import "strconv"

type Messanger struct {
	messages []string
}

func (messager *Messanger) message(message_line string, line int) {
	if line > 0 {
		message_line = message_line + " Line " + strconv.Itoa(line)
	}
	//message_line = message_line + "\n"
	messager.messages = append(messager.messages, message_line)
}
