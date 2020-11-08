package main

import (
	"fmt"
	"reflect"
)

// NewPromise Returns a New Promise
func NewPromise(
	callback func(
		rs func(interface{}) (interface{}, error),
		rj func(error) (interface{}, error),
	) (interface{}, error)) *Promise {

	resolve := func(v interface{}) (interface{}, error) {
		return v, nil
	}

	reject := func(e error) (interface{}, error) {
		return "", e
	}

	p := Promise{}
	p.channel = make(chan int)
	p.wg.Add(1)
	go func() {
		p.res, p.err = callback(resolve, reject)
		p.wg.Done()
	}()
	return &p
}

// Then Method
func (p *Promise) Then(f ...interface{}) *Promise {

	fmt.Println("Then Called")

	var r func(r interface{}) (*Promise, interface{}, error)
	var e func(err error) (*Promise, interface{}, error)

	if len(f) >= 1 {
		r = reflect.ValueOf(f[0]).Interface().(func(r interface{}) (*Promise, interface{}, error))
	}

	if len(f) == 2 {
		e = reflect.ValueOf(f[1]).Interface().(func(err error) (*Promise, interface{}, error))
	}

	var p1 *Promise

	go func() {
		p.wg.Wait()
		if p.err != nil {

			if e == nil {
				fmt.Println("e is null")
				p1 = p
				p.channel <- 1
				return
			}
			// prom1, val1, err1 :=
			p1 = populatePromise(e(p.err))
			p.channel <- 1
			return
		}
		if r == nil {
			fmt.Println("r is null")
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

	fmt.Println(len(f))

	if len(f) == 0 {
		return p
	}

	return p.Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Callin then from Catch", r)
			return nil, r, nil
		},
		f[0],
	)
}

// Finally Method
func (p *Promise) Finally(f func()) *Promise {

	go func() {
		p.wg.Wait()
		f()
		p.channel <- 1
	}()
	<-p.channel
	return p
}
