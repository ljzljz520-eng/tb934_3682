package clock

import "sync"

type Clock interface {
	Now() string
}

type Fixed struct {
	mu    sync.RWMutex
	value string
}

func NewFixed(value string) *Fixed {
	return &Fixed{value: value}
}

func (f *Fixed) Now() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.value
}

func (f *Fixed) Set(value string) {
	f.mu.Lock()
	f.value = value
	f.mu.Unlock()
}

type Sequence struct {
	mu     sync.Mutex
	values []string
	index  int
}

func NewSequence(values ...string) *Sequence {
	return &Sequence{values: append([]string(nil), values...)}
}

func (s *Sequence) Now() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.values) == 0 {
		return ""
	}
	value := s.values[s.index%len(s.values)]
	s.index++
	return value
}
