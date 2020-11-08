package main

// Promiser Interface has the basic methods for a promise
type Promiser interface {
	Then(funcs ...interface{}) *Promise
	Catch(funcs ...interface{}) *Promise
	Finally(f func()) *Promise
}

// Promise has information about the the state and results if any
// If status is PENDING, res and err will be nil
// If status is Resolved, res will have a value
// If status is REJECTED, err will have a value
type Promise struct {
	channel chan int
	res     interface{}
	err     error
	status  int
}

const (

	// PENDING State - Promise is not resolved or rejected
	PENDING int = 0

	// RESOLVED State - Promise is resolved
	RESOLVED int = 1

	// REJECTED State - Promise is rejected
	REJECTED int = 2
)
