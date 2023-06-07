package devrunner

import (
	"bytes"
	"sync"
)

type ThreadSafeBuffer struct {
	b *bytes.Buffer
	m sync.RWMutex
}

func newThreadSafeBuffer() *ThreadSafeBuffer {
	return &ThreadSafeBuffer{
		b: bytes.NewBuffer(nil),
	}
}

func (b *ThreadSafeBuffer) Read(p []byte) (n int, err error) {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.b.Read(p)
}
func (b *ThreadSafeBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *ThreadSafeBuffer) String() string {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.b.String()
}
