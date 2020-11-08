package main

import (
	"reflect"
)

// NewPromise creates a Promise with default values. It has the implementation of resolve and reject methods
func NewPromise(callback func(rs func(interface{}), rj func(error))) *Promise {

	// Define a empty Promise with default values
	p := Promise{}
	p.channel = make(chan int)
	p.status = PENDING

	// Resolves the promise
	resolve := func(v interface{}) {
		p.status = RESOLVED
		p.res = v
	}

	// Rejects the promise
	reject := func(e error) {
		p.status = REJECTED
		p.err = e
	}

	go func() {
		callback(resolve, reject)
		p.channel <- 1
	}()
	<-p.channel
	return &p
}

// Then method returns a Promise. It takes up to two arguments: callback functions for the success and failure cases of the Promise.
func (p *Promise) Then(f ...interface{}) *Promise {

	// Check if previous promise is settled
	if p.status == PENDING {
		return p
	}

	var r func(r interface{}) (*Promise, interface{}, error)
	var e func(err error) (*Promise, interface{}, error)

	// Assign value to result callback function
	if len(f) >= 1 && f[0] != nil {
		r = reflect.ValueOf(f[0]).Interface().(func(r interface{}) (*Promise, interface{}, error))
	}

	// Assign value to errored callback function
	if len(f) >= 2 && f[1] != nil {
		e = reflect.ValueOf(f[1]).Interface().(func(err error) (*Promise, interface{}, error))
	}

	var p1 *Promise

	go func() {

		// If error found
		if p.err != nil {

			// If error function does not exist return
			if e == nil {
				p1 = p
				p.channel <- 1
				return
			}
			// If error function does exist
			// Get the Promise to be returned, by calling the error callback
			p1 = populatePromise(e(p.err))
			p.channel <- 1
			return
		}
		// If no error found
		// If result callback does not exist return
		if r == nil {
			p1 = p
			p.channel <- 1
			return
		}
		// Get the Promise to be returned, by calling the result callback
		p1 = populatePromise(r(p.res))
		p.channel <- 1
	}()

	// Wait for channel and return
	<-p.channel
	return p1
}

// Catch method returns a Promise and deals with rejected cases only
func (p *Promise) Catch(f ...interface{}) *Promise {

	// Check if callback function provided, if not then do nothing
	if len(f) == 0 {
		return p
	}

	// Call the Then method with success callback as nil and failure callback as the one provided in this function
	return p.Then(
		nil,
		f[0],
	)
}

// Finally Method returns a Promise. When the promise is settled, i.e either fulfilled or rejected, the specified callback function is executed. This provides a way for code to be run whether the promise was fulfilled successfully or rejected once the Promise has been dealt with.
func (p *Promise) Finally(f func()) *Promise {

	// Check if previous promise is settled
	if p.status == PENDING {
		return p
	}

	go func() {
		f()
		p.channel <- 1
	}()

	// Wait for channel and return
	<-p.channel
	return p
}
