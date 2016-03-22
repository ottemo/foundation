package utils

import (
	"sync"
	"reflect"
	"errors"
	"fmt"
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
	if subject == nil {
		return errors.New("subject is nil")
	}

	rSubject := reflect.ValueOf(subject)
	rSubjectType := rSubject.Type()

	var rKey reflect.Value
	var rKeyType reflect.Type

	rValue := reflect.ValueOf(value)
	rValueType := rValue.Type()

	length := len(path)
	last := length-1
	for idx, key := range path {
		switch rSubject.Kind() {
		case reflect.Map:
			rKey = reflect.ValueOf(key)
			rKeyType = rKey.Type()
			if rKey.IsNil() || rKeyType != rSubjectType.Key() {
				return errors.New(fmt.Sprintf("invalid path key %d - %v" % rKey))
			}

			if idx != last {
				m := GetMutex(rSubject.Pointer())
				if m == nil {
					return errors.New(fmt.Sprintf("invalid mutex on %v" % rSubject))
				}

				m.Lock()
				rSubject = rSubject.MapIndex(rKey)
				m.Unlock()
			}
		}

		if rSubject.IsNil() {
			return errors.New("subject is nil")
		}
	}


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