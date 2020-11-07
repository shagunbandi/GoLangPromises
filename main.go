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
		return v, nil
	}

	rejectFunc := func(e error) (interface{}, error) {
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

func populatePromise(prom1 *Promise, val1 interface{}, err1 error) *Promise {
	if prom1 == nil {
		prom1 = getPromiseOrEmptyPromise(nil)
		if val1 == nil {
			prom1.err = err1
		}
		if val1 != nil {
			prom1.res = val1
		}
	}
	return prom1
}

// Then Method
func (p *Promise) Then(funcs ...interface{}) *Promise {

	fmt.Println("Then Called")

	var r func(r interface{}) (*Promise, interface{}, error)
	var e func(err error) (*Promise, interface{}, error)

	if len(funcs) >= 1 {
		r = reflect.ValueOf(funcs[0]).Interface().(func(r interface{}) (*Promise, interface{}, error))
	}

	if len(funcs) == 2 {
		e = reflect.ValueOf(funcs[1]).Interface().(func(err error) (*Promise, interface{}, error))
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
			prom1, val1, err1 := e(p.err)
			p1 = populatePromise(prom1, val1, err1)
			p.channel <- 1
			return
		}
		if r == nil {
			fmt.Println("r is null")
			p1 = p
			p.channel <- 1
			return
		}

		prom1, val1, err1 := r(p.res)
		p1 = populatePromise(prom1, val1, err1)
		p.channel <- 1
	}()
	<-p.channel
	return p1
}

// Catch Method
func (p *Promise) Catch(funcs ...interface{}) *Promise {

	fmt.Println(len(funcs))

	if len(funcs) == 0 {
		return p
	}

	return p.Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Callin then from Catch", r)
			return nil, r, nil
		},
		funcs[0],
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

func main() {

	fmt.Println("Starting\n")

	links := []string{
		"http://google.com",
		"http://youtube.com",
		"http://facebook.com",
	}

	link := links[0]
	link2 := links[1]
	link3 := links[2]

	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			_, err := http.Get(link)
			if err != nil {
				return reject(fmt.Errorf("%v is down :(", link))
			}
			return resolve(link + " is up :)")
		},
	).Finally(
		func() {
			fmt.Println("Finally 1")
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("On Fail 0000000", err)
			return nil, nil, nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success >>>>", r)
			return nil, nil, fmt.Errorf("Throwing Exception")
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success >>>>", r)
			return nil, "Should Catch Exception, somethign is wrong", nil
		},
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("Failed >>>>", err)
			return nil, "Exception Caught Correctly 111", nil
		},
	).Catch().Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("On Success 1111 >> ", r)
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
			fmt.Println("Success >>>>", r)
			return nil, "Just a value", nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success <<<<<", r)
			return nil, nil, fmt.Errorf("Just a Fail")
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success >>>>", r)
			return nil, "Just a value", nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("On Fail2", err)
			return NewPromise(
				func(
					resolve func(v interface{}) (interface{}, error),
					reject func(e error) (interface{}, error),
				) (interface{}, error) {

					_, err := http.Get(link3)
					if err != nil {
						return reject(fmt.Errorf("%v is down :(", link3))
					}
					return resolve(link3 + " is up :)")
				},
			), nil, nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("On Fail3", err)
			return nil, nil, nil
		},
	).Finally(
		func() {
			fmt.Println("Finally 2")
		},
	)

	fmt.Println("All Done")
	// time.Sleep(10 * time.Second)

}
