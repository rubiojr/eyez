package tui

import (
	"errors"
	"sync"
)

type RStack struct {
	lock sync.Mutex
	s    []*Record
}

func NewRStack() *RStack {
	return &RStack{sync.Mutex{}, make([]*Record, 0)}
}

func (s *RStack) Push(e *Record) {
	n := []*Record{e}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.s = append(n, s.s...)
}

func (s *RStack) Pop() (*Record, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]

	return res, nil
}

func (s *RStack) Each(f func(*Record)) {
	for {
		rec, err := s.Pop()
		if err != nil {
			break
		}
		f(rec)
	}
}
