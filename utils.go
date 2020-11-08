package main

// returns the same promise if exists, or creates a new one with default value if does not exist
func getPromiseOrEmptyPromise(p *Promise) *Promise {
	if p != nil {
		return p
	}
	p1 := &Promise{}
	p1.channel = make(chan int)
	p1.res = nil
	p1.err = nil
	p1.status = PENDING
	return p1
}

// Return the same promise if exists, else creates a new one with default value, and depending upon the value or error provided populates the promise
func populatePromise(p *Promise, r interface{}, e error) *Promise {

	// If promise is undefined
	if p == nil {

		// Create a new promise
		p = getPromiseOrEmptyPromise(nil)

		// If result provided, set status and result values
		if r != nil {
			p.status = RESOLVED
			p.res = r
			p.err = nil
		}

		// If error provided, set status and error values
		if e != nil {
			p.status = REJECTED
			p.err = e
			p.res = nil
		}

		// If result and error both not provided, set status as RESOLVED (Case where no more then or catch are required)
		if r == nil && e == nil {
			p.status = RESOLVED
		}
	}
	return p
}
