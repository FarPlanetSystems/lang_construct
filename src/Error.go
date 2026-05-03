package main

type Error struct {
	Message string
	isError bool
	line    int
}

func createError() Error {
	return Error{isError: false}
}
