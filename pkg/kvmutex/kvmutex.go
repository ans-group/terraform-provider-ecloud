package kvmutex

import (
	"sync"
)

type LockerFactory interface {
	Get() sync.Locker
}

type KVMutexLockerFactory struct{}

func (f *KVMutexLockerFactory) Get() sync.Locker {
	return &sync.Mutex{}
}

type KVMutex struct {
	mutexesLock sync.Mutex
	mutexes     map[string]sync.Locker
	factory     LockerFactory
}

// NewKVMutex returns an initialised KVMutex struct
func NewKVMutex() *KVMutex {
	return &KVMutex{
		mutexes: make(map[string]sync.Locker),
		factory: &KVMutexLockerFactory{},
	}
}

func (m *KVMutex) WithFactory(f LockerFactory) *KVMutex {
	m.factory = f
	return m
}

func (m *KVMutex) Lock(key string) func() {
	mutex := m.get(key)
	mutex.Lock()
	return func() {
		mutex.Unlock()
	}
}

func (m *KVMutex) Unlock(key string) {
	m.get(key).Unlock()
}

func (m *KVMutex) get(key string) sync.Locker {
	m.mutexesLock.Lock()
	defer m.mutexesLock.Unlock()
	_, ok := m.mutexes[key]
	if !ok {
		m.mutexes[key] = m.factory.Get()
	}
	return m.mutexes[key]
}
