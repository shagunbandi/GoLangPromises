package main

import (
	"fmt"
	"net/http"
)

func main() {

	links := []string{
		"http://google.com",
		"http://youtube.com",
		"http://facebook.com",
	}

	link := links[0]
	link2 := links[1]
	link3 := links[2]
	var p Promiser
	p = NewPromise(
		func(
			resolve func(v interface{}),
			reject func(e error),
		) {
			_, err := http.Get(link)
			if err != nil {
				reject(fmt.Errorf("%v is down :(", link))
			}
			resolve(link + " is up :)")
		},
	)
	p.Finally(
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
					resolve func(v interface{}),
					reject func(e error),
				) {

					_, err := http.Get(link2)
					if err != nil {
						reject(fmt.Errorf("%v is down :(", link2))
					}
					resolve(link2 + " is up :)")
				},
			), nil, nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			fmt.Println("Success >>>>", r)
			return nil, "Just a value", nil
		},
	).Finally(
		func() {
			fmt.Println("Finally 2")
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
					resolve func(v interface{}),
					reject func(e error),
				) {

					_, err := http.Get(link3)
					if err != nil {
						reject(fmt.Errorf("%v is down :(", link3))
					}
					resolve(link3 + " is up :)")
				},
			), nil, nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			fmt.Println("On Fail3", err)
			return nil, nil, nil
		},
	)
}
