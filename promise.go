package main

import (
	"reflect"
)

// NewPromise Returns a New Promise
func NewPromise(
	callback func(
		rs func(interface{}) (interface{}, error),
		rj func(error) (interface{}, error),
	) (interface{}, error)) *Promise {

	p := Promise{}
	p.channel = make(chan int)
	p.status = 0
	resolve := func(v interface{}) (interface{}, error) {
		p.status = 1
		return v, nil
	}

	reject := func(e error) (interface{}, error) {
		p.status = 2
		return "", e
	}

	p.wg.Add(1)
	go func() {
		p.res, p.err = callback(resolve, reject)
		p.wg.Done()
	}()
	return &p
}

// Then Method
func (p *Promise) Then(f ...interface{}) *Promise {

	p.wg.Wait()
	if p.status == 0 {
		return p
	}

	var r func(r interface{}) (*Promise, interface{}, error)
	var e func(err error) (*Promise, interface{}, error)

	if len(f) >= 1 {
		r = reflect.ValueOf(f[0]).Interface().(func(r interface{}) (*Promise, interface{}, error))
	}

	if len(f) >= 2 {
		e = reflect.ValueOf(f[1]).Interface().(func(err error) (*Promise, interface{}, error))
	}

	var p1 *Promise

	go func() {
		if p.err != nil {
			if e == nil {
				p1 = p
				p.channel <- 1
				return
			}
			p1 = populatePromise(e(p.err))
			p.channel <- 1
			return
		}
		if r == nil {
			p1 = p
			p.channel <- 1
			return
		}
		p1 = populatePromise(r(p.res))
		p.channel <- 1
	}()
	<-p.channel
	return p1
}

// Catch Method
func (p *Promise) Catch(f ...interface{}) *Promise {

	if len(f) == 0 {
		return p
	}

	return p.Then(
		func(r interface{}) (*Promise, interface{}, error) {
			return nil, r, nil
		},
		f[0],
	)
}

// Finally Method
func (p *Promise) Finally(f func()) *Promise {
	p.wg.Wait()
	go func() {
		f()
		p.channel <- 1
	}()
	<-p.channel
	return p
}
