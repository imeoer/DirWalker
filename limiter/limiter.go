package limiter

import (
	"sync"
)

const (
	MAX_OPEN_NUMBER = 500
)

type Limiter struct {
	wg    sync.WaitGroup
	count chan bool
}

func New() *Limiter {
	limiter := new(Limiter)
	limiter.count = make(chan bool, MAX_OPEN_NUMBER)
	return limiter
}

func (limiter *Limiter) Add() {
	limiter.count <- true
	limiter.wg.Add(1)
}

func (limiter *Limiter) Done() {
	<-limiter.count
	limiter.wg.Done()
}

func (limiter *Limiter) Wait() {
	limiter.wg.Wait()
}
