package compiler_objects

import "fmt"

type Formula struct {
	Line              int
	HeadToken         *TokenNode
	tailToken         *TokenNode
	previousStatement *Formula
}

func (statement Formula) IsEmpty() bool {
	return statement.HeadToken == nil
}

func CreateStatement(line int) Formula {
	res := Formula{
		Line:      line,
		HeadToken: nil,
		tailToken: nil,
	}
	return res
}

type StatementQueue struct {
	headStatement *Formula
	tailStatement *Formula
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

func (queue *StatementQueue) Enqueue(statement Formula) {
	if queue.IsEmpty {
		queue.headStatement = &statement
		queue.tailStatement = &statement
	} else {
		queue.tailStatement.previousStatement = &statement
		queue.tailStatement = &statement
	}
	queue.IsEmpty = false
}

func (queue *StatementQueue) Face() *Formula {
	return queue.headStatement
}

type TokenNode struct {
	Content             Token
	PreviousNode        *TokenNode
	IsPreviousNodeEmpty bool
}

func (statement *Formula) Enqueue(token Token) {

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

func (statement *Formula) Dequeue() {
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

func (statement *Formula) Face() Token {
	if statement.IsEmpty() {
		return Token{}
	}
	return statement.HeadToken.Content
}
func PrintStatement(statement Formula) {
	statement_copy := statement
	ref := &statement_copy
	for !ref.IsEmpty() {
		fmt.Println(ref.Face())
		ref.Dequeue()
	}
}

func (statement1 Formula) Compare(statement2 Formula) bool {
	for !statement1.IsEmpty() && !statement2.IsEmpty() {
		if statement1.Face().Value != statement2.Face().Value {
			return false
		}
		statement1.Dequeue()
		statement2.Dequeue()
	}
	if statement1.IsEmpty() && statement2.IsEmpty() {
		return true
	}
	return false
}

func (stmt Formula) Exchange(value Formula, tokenSequence Formula) Formula {
	result := CreateStatement(stmt.Line)
	for !stmt.IsEmpty() {
		if !stmt.containsStatementFace(value) {
			result.Enqueue(stmt.Face())
			stmt.Dequeue()
		} else {
			sequence := tokenSequence
			for !sequence.IsEmpty() {
				result.Enqueue(sequence.Face())
				sequence.Dequeue()
			}
			stmt.Dequeue()
		}
	}
	return result

}

func (stmt Formula) containsStatementFace(contains Formula) bool {
	con := contains
	for !con.IsEmpty() {
		if con.Face() != stmt.Face() {
			return false
		}
		con.Dequeue()
		stmt.Dequeue()
	}
	return true
}

func mergeStatements(head Formula, tail Formula) Formula {
	res := head
	add := tail
	for !add.IsEmpty() {
		res.Enqueue(add.Face())
		add.Dequeue()
	}
	return res
}
