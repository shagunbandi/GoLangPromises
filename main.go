package main

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

// Promise Struct
type Promise struct {
	wg      sync.WaitGroup
	channel chan int
	res     interface{}
	err     error
}

// NewPromise Returns a New Promise
func NewPromise(
	callback func(
		resolve func(interface{}) (interface{}, error),
		reject func(error) (interface{}, error),
	) (interface{}, error)) *Promise {

	resolveFunc := func(v interface{}) (interface{}, error) {
		fmt.Println("Resolved")
		return v, nil
	}

	rejectFunc := func(e error) (interface{}, error) {
		fmt.Println("Rejected")
		return "", e
	}

	p := Promise{}
	p.channel = make(chan int)
	p.wg.Add(1)
	go func() {
		p.res, p.err = callback(resolveFunc, rejectFunc)
		p.wg.Done()
	}()
	return &p
}

func getPromiseOrEmptyPromise(p *Promise) *Promise {
	if p != nil {
		return p
	}
	p1 := &Promise{}
	p1.channel = make(chan int)
	p1.res = nil
	p1.err = nil
	return p1
}

// Then Method
func (p *Promise) Then(funcs ...interface{}) *Promise {

	var r func(r interface{}) (*Promise, interface{}, error)
	var e func(err error) (*Promise, interface{}, error)

	if len(funcs) == 1 {
		r = reflect.ValueOf(funcs[0]).Interface().(func(r interface{}) (*Promise, interface{}, error))
	}

	if len(funcs) == 2 {
		e = reflect.ValueOf(funcs[1]).Interface().(func(err error) (*Promise, interface{}, error))
	}

	var p1 *Promise

	go func() {
		p.wg.Wait()
		if p.err != nil {
			fmt.Println("Found Error")

			if e == nil {
				p1 = p
				p.channel <- 1
				return
			}
			prom1, val1, err1 := e(p.err)
			if prom1 == nil {
				prom1 = getPromiseOrEmptyPromise(nil)
				if val1 == nil {
					prom1.err = err1
				}
				if val1 != nil {
					prom1.res = val1
				}
			}
			p1 = prom1
			p.channel <- 1
			return
		}
		fmt.Println("No Error")
		if r == nil {
			p1 = p
			p.channel <- 1
			return
		}

		prom1, val1, err1 := r(p.res)
		if prom1 == nil {
			prom1 = getPromiseOrEmptyPromise(nil)
			if val1 == nil {
				prom1.err = err1
			}
			if val1 != nil {
				prom1.res = val1
			}
		}
		p1 = prom1
		p.channel <- 1
	}()
	<-p.channel
	return getPromiseOrEmptyPromise(p1)
}

// Catch Method
func (p *Promise) Catch(funcs ...interface{}) *Promise {

	var e func(err error) (*Promise, interface{}, error)

	if len(funcs) == 1 {
		e = reflect.ValueOf(funcs[0]).Interface().(func(err error) (*Promise, interface{}, error))
	}

	var p1 *Promise

	go func() {
		p.wg.Wait()
		if p.err != nil {

			if e == nil {
				p1 = p
				p.channel <- 1
				return
			}

			prom1, val1, err1 := e(p.err)
			if prom1 == nil {
				prom1 = getPromiseOrEmptyPromise(nil)
				if val1 == nil {
					prom1.err = err1
				}
				if val1 != nil {
					prom1.res = val1
				}
			}
			p1 = prom1
			p.channel <- 1
			return
		}
		p1 = p
		p.channel <- 1
	}()
	<-p.channel
	return getPromiseOrEmptyPromise(p1)
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

func main() {

	links := []string{
		"http://google.com",
		"http://youtube.com",
		"http://facebook.com",
	}

	link := links[0]
	link2 := links[1]

	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			fmt.Println("Calling Now")
			_, err := http.Get(link)
			fmt.Println("Got the Result")
			if err != nil {
				return reject(fmt.Errorf("%v is down :(", link))
			}
			return resolve(link + " is up :)")
		},
	).Finally(
		func() {
			fmt.Println("FinallyFinallyFinallyFinally")
		},
	).Catch().Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("On Success", r)
			return NewPromise(
				func(
					resolve func(v interface{}) (interface{}, error),
					reject func(e error) (interface{}, error),
				) (interface{}, error) {

					_, err := http.Get(link2)
					if err != nil {
						return reject(fmt.Errorf("%v is down :(", link2))
					}
					return resolve(link2 + " is up :)")
				},
			), nil, nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success <<<<<", r)
			return nil, nil, fmt.Errorf("Just a value")
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success >>>>", r)
			return nil, "Just a value", nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("On Fail2", err)
			return nil, nil, nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("On Fail3", err)
			return nil, nil, nil
		},
	).Finally(
		func() {
			fmt.Println("FinallyFinallyFinallyFinally")
		},
	)

	fmt.Println("I'm Here")

}
