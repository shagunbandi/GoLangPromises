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

	p := NewPromise(
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
	)

	fmt.Println(p)

	fmt.Println("All Done")
	// time.Sleep(10 * time.Second)

}
