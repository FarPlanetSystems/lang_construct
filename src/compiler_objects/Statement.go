package compiler_objects

import "fmt"

type Statement struct {
	Line              int
	HeadToken         *TokenNode
	tailToken         *TokenNode
	previousStatement *Statement
}

func (statement Statement) IsEmpty() bool {
	return statement.HeadToken == nil
}

func CreateStatement() Statement {
	res := Statement{
		Line:      0,
		HeadToken: nil,
		tailToken: nil,
	}
	return res
}

type StatementQueue struct {
	headStatement *Statement
	tailStatement *Statement
	IsEmpty       bool
}

func CreateStatementQueue() *StatementQueue {
	res := StatementQueue{
		headStatement: nil,
		tailStatement: nil,
		IsEmpty:       true,
	}
	return &res
}

func (queue *StatementQueue) Dequeue() {

	if queue.IsEmpty {
		return
	}
	if queue.headStatement.previousStatement == nil {
		queue.headStatement = nil
		queue.tailStatement = nil
		queue.IsEmpty = true
	} else {
		queue.headStatement = queue.headStatement.previousStatement
	}
	//fmt.Println("dequeuing: previous statement:")
	//fmt.Println(queue.headStatement.previousStatement)
}

func (queue *StatementQueue) Enqueue(statement Statement) {
	if queue.IsEmpty {
		queue.headStatement = &statement
		queue.tailStatement = &statement
	} else {
		queue.tailStatement.previousStatement = &statement
		queue.tailStatement = &statement
	}
	queue.IsEmpty = false
}

func (queue *StatementQueue) Face() *Statement {
	return queue.headStatement
}

type TokenNode struct {
	Content             Token
	PreviousNode        *TokenNode
	IsPreviousNodeEmpty bool
}

func (statement *Statement) Enqueue(token Token) {

	tokenNode := TokenNode{
		Content:             token,
		PreviousNode:        nil,
		IsPreviousNodeEmpty: true,
	}

	if statement.IsEmpty() {
		statement.HeadToken = &tokenNode
		statement.tailToken = &tokenNode
	} else {

		statement.tailToken.PreviousNode = &tokenNode
		statement.tailToken.IsPreviousNodeEmpty = false
		statement.tailToken = &tokenNode
		//fmt.Println("enqueueing " + statement.tailToken.content.token_type)
	}
}

func (statement *Statement) Dequeue() {
	if statement.IsEmpty() {
		return
	}
	if statement.HeadToken.IsPreviousNodeEmpty {
		statement.HeadToken = nil
		statement.tailToken = nil
	} else {
		statement.HeadToken = statement.HeadToken.PreviousNode
	}
}

func (statement *Statement) Face() Token {
	if statement.IsEmpty() {
		return Token{}
	}
	return statement.HeadToken.Content
}
func PrintStatement(statement Statement) {
	statement_copy := statement
	ref := &statement_copy
	for !ref.IsEmpty() {
		fmt.Println(ref.Face())
		ref.Dequeue()
	}
}

func mergeStatements(head Statement, tail Statement) Statement {
	res := head
	add := tail
	for !add.IsEmpty() {
		res.Enqueue(add.Face())
		add.Dequeue()
	}
	return res
}
