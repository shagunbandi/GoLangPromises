package main

import (
	"fmt"
	"testing"
)

func TestResolvedThenAndNoCatch(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			t.Errorf("Should not enter catch")
			return nil, nil, nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Resolved" {
				t.Errorf("Expected 'Resolved', got '%v'", res)
			}
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
				t.Errorf("Expected 'Rejected', got '%v'", err.Error())
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
				t.Errorf("Expected 'Resolved', got '%v'", res)
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
				t.Errorf("Expected 'Rejected', got '%v'", err.Error())
			}
			return nil, nil, nil
		},
	)
}

func TestUnresolvedPromiseNoThenNoCatch(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return nil, nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			t.Errorf("Should not enter 'then', if promise unresolved")
			return nil, nil, nil
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			t.Errorf("Should not enter 'catch', if promise unresolved")
			return nil, nil, nil
		},
	)
}

func TestReturnValueFromBlock(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Resolved" {
				t.Errorf("Expected 'Resolved', got '%v'", res)
			}
			return nil, "Return from then", nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Return from then" {
				t.Errorf("Expected 'Return from then', got '%v'", res)
			}
			return nil, nil, nil
		},
	)
}

func TestReturnErrorFromBlock(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Resolved" {
				t.Errorf("Expected 'Resolved', got '%v'", res)
			}
			return nil, nil, fmt.Errorf("Error from then")
		},
	).Catch(
		func(err error) (*Promise, interface{}, error) {
			if err.Error() != "Error from then" {
				t.Errorf("Expected 'Error from then', got '%v'", err.Error())
			}
			return nil, nil, nil
		},
	)
}

func TestPromiseReturn(t *testing.T) {
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved1")
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			return NewPromise(
				func(
					resolve func(v interface{}) (interface{}, error),
					reject func(e error) (interface{}, error),
				) (interface{}, error) {
					return resolve("Resolved2")
				},
			), nil, nil
		},
	).Then(
		func(r interface{}) (*Promise, interface{}, error) {
			res := fmt.Sprintf("%v", r)
			if res != "Resolved2" {
				t.Errorf("Expected 'Resolved2', got '%v'", res)
			}
			return nil, nil, nil
		},
	)
}

func TestFinally(t *testing.T) {
	v := 0
	NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	).Finally(
		func() {
			v = 1
		},
	)
	if v == 0 {
		t.Errorf("Expected 'v=0', got 'v=%v'", v)
	}
}

func TestUnresolvedPromise(t *testing.T) {
	p := NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return nil, nil
		},
	)
	if p.status != 0 {
		t.Errorf("Expected Status '0', got '%v'", p.status)
	}
}

func TestResolvedPromise(t *testing.T) {
	p := NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return resolve("Resolved")
		},
	)
	if p.status != 1 {
		t.Errorf("Expected Status '1', got '%v'", p.status)
	}
}

func TestRejectedPromise(t *testing.T) {
	p := NewPromise(
		func(
			resolve func(v interface{}) (interface{}, error),
			reject func(e error) (interface{}, error),
		) (interface{}, error) {
			return reject(fmt.Errorf("Rejected"))
		},
	)
	if p.status != 2 {
		t.Errorf("Expected Status '2', got '%v'", p.status)
	}
}
