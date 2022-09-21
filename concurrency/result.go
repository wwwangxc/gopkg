package concurrency

import (
	"bytes"
	"fmt"
	"sync"
)

// Result set
type Result struct {
	resultSet []*singleResult
	failed    bool

	rw sync.RWMutex
}

// Succeed return true when no failed result
func (r *Result) Succeed() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return !r.failed
}

// Failed return true when there are failed result
func (r *Result) Failed() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.failed
}

// MergedError return merged error
//
// Format like:
// 	2 errors occurred:
// 		* error message ...
// 		* [PANIC]panic message ...
func (r *Result) MergedError() error {
	if r.Succeed() {
		return nil
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	errNum := 0
	errMsg := bytes.NewBufferString("")
	for _, v := range r.resultSet {
		if v.err == nil {
			continue
		}

		errNum++
		fmt.Fprintf(errMsg, "\n    * %v", v.err)
	}

	return fmt.Errorf("%d errors occurred:%s", errNum, errMsg.String())
}

// Errors return all error
//
// return nil when no failed result
func (r *Result) Errors() []error {
	if r.Succeed() {
		return nil
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	errs := make([]error, 0, len(r.resultSet))
	for _, v := range r.resultSet {
		if v.err == nil {
			continue
		}

		errs = append(errs, v.err)
	}

	return errs
}

// Results return all result
//
// return nil when there are failed result
func (r *Result) Results() []interface{} {
	if r.Failed() {
		return nil
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	rets := make([]interface{}, 0, len(r.resultSet))
	for _, v := range r.resultSet {
		if v.err != nil || v.result == nil {
			continue
		}

		rets = append(rets, v.result)
	}

	return rets
}

func (r *Result) append(result interface{}, err error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if err != nil {
		r.failed = true
	}

	r.resultSet = append(r.resultSet, &singleResult{
		result: result,
		err:    err,
	})
}

type singleResult struct {
	result interface{}
	err    error
}
