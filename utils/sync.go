package utils

import (
	"sync"
	"reflect"
	"errors"
)

var (
	locks = make(map[uintptr]*syncMutex)
	locksMutex sync.Mutex
)

type syncMutex struct {
	index uintptr
	refs  int
	lock  bool
	mutex sync.Mutex
}

func (it *syncMutex) Lock() {
	it.mutex.Lock()
	it.lock = true
}

func (it *syncMutex) Unlock() {
	locksMutex.Lock()
	defer locksMutex.Unlock()

	it.lock = false
	it.refs--
	if it.refs == 0 {
		delete(locks, it.index)
	}

	it.mutex.Unlock()
}

func (it *syncMutex) GetIndex() uintptr {
	return it.index
}

func (it *syncMutex) IsLocked() bool {
	return it.lock
}

func (it *syncMutex) Refs() int {
	return it.refs
}

func GetPointer(subject interface{}) uintptr {
	if subject == nil {
		return 0
	}

	value := reflect.ValueOf(subject)
	switch value.Kind() {
	case reflect.Chan,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer,
		reflect.Func,
		reflect.Slice,
		reflect.Array:

		return value.Pointer()
	}

	return 0
}

func GetMutex(subject interface{}) *syncMutex {
	locksMutex.Lock()
	defer locksMutex.Unlock()

	if index := GetPointer(subject); index != 0 {
		mutex, present := locks[index]
		if !present {
			mutex = new(syncMutex)
			mutex.index = index
			locks[index] = mutex
		}
		mutex.refs++
		return mutex
	}

	return nil
}

func SyncSet(subject interface{}, value interface{}, path ...interface{}) error {
	rSubject := reflect.ValueOf(subject)
	for _, x := range path {
		if rSubject.IsNil() {
			return errors.New("path not found")
		}

		switch rSubject.Kind() {
		case reflect.Map:
			xType := reflect.TypeOf(x)
			rSubjectType := rSubject.Type()
			if xType != rSubjectType {
				return errors.New("wrong key type")
			}

			m := GetMutex(rSubject.Pointer())
			if m == nil {
				return errors.New("invalid mutex")
			}

			m.Lock()
			rSubject = rSubject.MapIndex(reflect.ValueOf(x))
			m.Unlock()
		}
	}

	if rSubject.IsNil() {
		return errors.New("invalid object")
	}

	valueType := reflect.TypeOf(value)
	rSubjectType := rSubject.Type()
	if rSubjectType != valueType {
		return errors.New("invalid value type")
	}

	m := GetMutex(rSubject.Pointer())
	if m == nil {
		return errors.New("invalid mutex")
	}

	m.Lock()
	rSubject = rSubject.MapIndex(reflect.ValueOf(x))
	m.Unlock()

	rSubject.Set(value)
	return nil
}

func SyncGet(subject interface{}, path ...interface{}) (interface{}, error) {
	return nil, nil
}