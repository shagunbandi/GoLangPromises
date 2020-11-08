package main

import "sync"

// Promiser Interface
type Promiser interface {
	Catch(funcs ...interface{}) *Promise
	Then(funcs ...interface{}) *Promise
	Finally(f func()) *Promise
}

// Promise Struct
type Promise struct {
	wg      sync.WaitGroup
	channel chan int
	res     interface{}
	err     error
	status  int
}
