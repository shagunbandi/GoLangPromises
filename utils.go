package main

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

func populatePromise(p *Promise, val1 interface{}, err1 error) *Promise {
	if p == nil {
		p = getPromiseOrEmptyPromise(nil)
		if val1 == nil {
			p.err = err1
		}
		if val1 != nil {
			p.res = val1
		}
	}
	return p
}
