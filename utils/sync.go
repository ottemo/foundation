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

	var value reflect.Value

	if rValue, ok := subject.(reflect.Value); ok {
		value = rValue
	} else {
		value = reflect.ValueOf(subject)
	}

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

func SyncMutex(subject interface{}) *syncMutex {
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

	// time critical segment (because of lock)
	initBlankValue := func() error {
		switch rSubjectType.Kind() {
		case reflect.Map:
			rSubject = reflect.MakeMap(rSubjectType)
		}
		return errors.New("subject is nil")
	}

	length := len(path)
	last := length-1
	for idx, key := range path {

		if rSubject.IsNil() && rSubjectType != nil {
		}

		switch rSubject.Kind() {
		case reflect.Map:
			rKey = reflect.ValueOf(key)
			rKeyType = rKey.Type()
			if rKeyType != rSubjectType.Key() {
				return errors.New("invalid path key")
			}

			if idx != last {
				m := SyncMutex(rSubject)
				if m == nil {
					return errors.New(fmt.Sprintf("invalid mutex on '%v' type '%v'", rKey, rSubjectType))
				}

				// time critical segment (because of lock)
				m.Lock()
				rSubjectValue := rSubject.MapIndex(rKey)
				if rSubjectValue.IsNil() {
					if err := initBlankValue(); err != nil {
						m.Unlock()
						return err
					}
				}
				m.Unlock()

				rSubject = rSubjectValue
				rSubjectType = rSubject.Type()
			}
		}



	}

	m := SyncMutex(rSubject.Pointer())
	if m == nil {
		return errors.New("invalid mutex creation for subject")
	}

	// allowing to pass value as setter function "func(oldValue) => newValue"
	isLocked := false
	if rValue.Kind() == reflect.Func {
		if rValueType.NumOut()==1 && rValueType.NumIn()==1 &&
			!rValueType.In(0).AssignableTo(rSubjectType) &&
			!rValueType.Out(0).AssignableTo(rSubjectType) {

			// the result is dependable on input, so we need
			// to keep read lock while the value would not
			// be updated
			m.Lock()
			isLocked = true
			rValue = rValue.Call([]reflect.Value{rValue})[0]
		}
	}

	switch rSubject.Kind() {
	case reflect.Map:
		if rKeyType != rSubjectType.Key() {
			return errors.New("invalid map key type")
		}

		if rKeyType != rValueType {
			return errors.New("invalid value type")
		}

		if !isLocked { m.Lock() }
		rSubject.SetMapIndex(rKey, rValue)
		m.Unlock()

	case reflect.Slice, reflect.Array:
		if rKey.Kind() != reflect.Int {
			return errors.New("invalid key, should be integer")
		}
		idx := int(rKey.Int())
		if rSubject.Len() <= idx {
			return errors.New("out of bound")
		}

		if !isLocked { m.Lock() }
		rSubject.Index(idx).Set(rValue)
		m.Unlock()

	case reflect.Ptr, reflect.Chan, reflect.Func:
		if !isLocked { m.Lock() }
		rSubject.Set(rValue)
		m.Unlock()

	default:
		return errors.New("invalid subject - must be pointer and not scalar value")
	}

	return nil
}

func SyncGet(subject interface{}, path ...interface{}) (interface{}, error) {



	return nil, nil
}