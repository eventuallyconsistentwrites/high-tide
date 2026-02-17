package countmin

import (
	"fmt"
	"sync"
)

type MapCounter struct {
	mapCounter map[string]int
	mu         sync.RWMutex
}

func NewMapCounter() *MapCounter {
	return &MapCounter{
		mapCounter: make(map[string]int),
	}
}
func (m *MapCounter) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("%#v", m.mapCounter)
}
func (m *MapCounter) Update(value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mapCounter[value]++
}
func (m *MapCounter) PointQuery(value string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.mapCounter[value]
}
func (m *MapCounter) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.mapCounter)
}
