package main

func getPromiseOrEmptyPromise(p *Promise) *Promise {
	if p != nil {
		return p
	}
	p1 := &Promise{}
	p1.channel = make(chan int)
	p1.res = nil
	p1.err = nil
	p1.status = 0
	return p1
}

func populatePromise(p *Promise, r interface{}, e error) *Promise {
	if p == nil {
		p = getPromiseOrEmptyPromise(nil)
		if r != nil {
			p.status = 1
			p.res = r
			p.err = nil
		}
		if e != nil {
			p.status = 2
			p.err = e
			p.res = nil
		}
		if r == nil && e == nil {
			p.status = 1
		}
	}
	return p
}
