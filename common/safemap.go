package common

import (
	"sync"
)

type SafeMap struct {
	sync.RWMutex
	Map map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[interface{}]interface{})
	return sm
}

func (sm *SafeMap) ReadMap(key interface{}) interface{} {
	sm.RLock()
	defer sm.RUnlock()
	value := sm.Map[key]
	return value
}

func (sm *SafeMap) WriteMap(key interface{}, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.Map[key] = value
}
