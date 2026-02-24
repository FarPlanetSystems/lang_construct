package main

import (
	"fmt"
	"strconv"
)

type Messenger struct {
	messages []string
}

func (messager *Messenger) InsertMessage(message_line string, line int) {
	if line > 0 {
		message_line = message_line + " Line " + strconv.Itoa(line)
	}
	//message_line = message_line + "\n"
	messager.messages = append(messager.messages, message_line)
}

func (messager *Messenger) Report() {
	for _, message := range messager.messages {
		fmt.Println(message)
	}
}
