package kvmutex

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestLocker struct {
	Locked bool
}

func (l *TestLocker) Lock() {
	l.Locked = true
}

func (l *TestLocker) Unlock() {
	l.Locked = false
}

type TestLockerFactory struct {
	locker sync.Locker
}

func (f *TestLockerFactory) Get() sync.Locker {
	return f.locker
}

func TestKVMutex_Lock_Locks(t *testing.T) {
	lock := &TestLocker{}
	mutex := NewKVMutex().WithFactory(&TestLockerFactory{locker: lock})

	mutex.Lock("somekey")

	assert.True(t, lock.Locked)
}

func TestKVMutex_Unlock_Unlocks(t *testing.T) {
	lock := &TestLocker{}
	mutex := NewKVMutex().WithFactory(&TestLockerFactory{locker: lock})

	mutex.Lock("somekey")
	assert.True(t, lock.Locked)

	mutex.Unlock("somekey")
	assert.False(t, lock.Locked)
}
