package main

import (
	"fmt"
	"testing"
)

func TestResolvedThenAndNoCatch(t *testing.T) {
	p := NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	)
	p.Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Resolved" {
				t.Errorf("Expected 'Resolved', got %v", res)
			}
			return nil, nil, nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			t.Errorf("Should not enter catch")
			return nil, nil, nil
		},
	)

}

func TestRejectedNoThenAndCatch(t *testing.T) {
	p := NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return reject(fmt.Errorf("Rejected"))
		},
	)
	p.Then(
		func(r interface{}) (*Promise, interface{}, error) {
			t.Errorf("Should not enter then")
			return nil, nil, nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			if err.Error() != "Rejected" {
				t.Errorf("Expected 'Rejected', got %v", err.Error())
			}
			return nil, nil, nil
		},
	)
}

func TestEmptyThen(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	).Then().Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Resolved" {
				t.Errorf("Expected 'Resolved', got %v", res)
			}
			return nil, nil, nil
		},
	)
}

func TestEmptyCatch(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return reject(fmt.Errorf("Rejected"))
		},
	).Catch().Catch(
		func(err error) (*Promise, interface{}, error) {
			if err.Error() != "Rejected" {
				t.Errorf("Expected 'Rejected', got %v", err.Error())
			}
			return nil, nil, nil
		},
	)
}
