package utils

import (
	"sync"
	"reflect"
	"errors"
	"fmt"
	"runtime/debug"
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

func GetPointer(subject interface{}) (uintptr, error) {
	if subject == nil {
		return 0, errors.New("can't get pointer to nil")
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

		return value.Pointer(), nil
	}

	debug.PrintStack()
	return 0, errors.New("can't get pointer to " + value.Type().String())
}

func SyncMutex(subject interface{}) (*syncMutex, error) {
	locksMutex.Lock()
	defer locksMutex.Unlock()

	index, err := GetPointer(subject)
	if err != nil {
		return nil, err
	}
	if index == 0 {
		return nil, errors.New("mutex to zero pointer")
	}
	mutex, present := locks[index]
	if !present {
		mutex = new(syncMutex)
		mutex.index = index
		locks[index] = mutex
	}
	mutex.refs++
	return mutex, nil
}

func SyncSet(subject interface{}, value interface{}, path ...interface{}) error {
	if subject == nil {
		return errors.New("subject is nil")
	}

	var err error

	rSubject := reflect.ValueOf(subject)
	rSubjectType := rSubject.Type()

	var rKey reflect.Value
	var rKeyType reflect.Type

	rValue := reflect.ValueOf(value)
	rValueType := rValue.Type()

	getValue := func(oldValue reflect.Value) reflect.Value {
		if rValue.Kind() == reflect.Func {
			if rValueType.NumOut()==1 && rValueType.NumIn()==1 {
				// oldValueType := oldValue.Type()
				// !rValueType.In(0).AssignableTo(oldValueType) &&
				//!rValueType.Out(0).AssignableTo(oldValueType) {
				return rValue.Call([]reflect.Value{oldValue})[0]
			}
		}
		return rValue
	}

	initBlankValue := func(valueType reflect.Type) (reflect.Value, error) {
		switch valueType.Kind() {
		case reflect.Map:
			return reflect.MakeMap(valueType), nil
		case reflect.Slice, reflect.Array:
			return reflect.MakeSlice(valueType, 0, 10), nil
		case reflect.Chan, reflect.Func:
			break
		default:
			return reflect.New(valueType).Elem(), nil
		}
		return reflect.ValueOf(nil), errors.New("unsuported blank value type " + valueType.String())
	}

	// If path specified - going to path element in subject
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
				return errors.New(fmt.Sprintf("invalid type of key %d", idx))
			}

			if idx != last {
				m, err := SyncMutex(rSubject)
				if err != nil {
					return err
				}

				// time critical segment (because of lock)
				m.Lock()
				rSubjectValue := rSubject.MapIndex(rKey)
				if !rSubjectValue.IsValid() {
					if rSubjectValue, err = initBlankValue(rSubjectType.Elem()); err != nil {
						m.Unlock()
						return err
					}
					rSubject.SetMapIndex(rKey, rSubjectValue)
				}
				m.Unlock()

				rSubject = rSubjectValue
				rSubjectType = rSubject.Type()
			}
		}



	}

	// setting value to subject
	m, err := SyncMutex(rSubject)
	if err != nil {
		return err
	}

	// setting the value according to subject type
	switch rSubject.Kind() {
	case reflect.Map:
		if rKeyType != rSubjectType.Key() {
			return errors.New("invalid map key type (" + rKeyType.String() + " != " + rSubjectType.Key().String() + ")")
		}

		m.Lock()
		oldValue := rSubject.MapIndex(rKey)
		if !oldValue.IsValid() {
			if oldValue, err = initBlankValue(rSubjectType.Elem()); err != nil {
				m.Unlock()
				return err
			}
		}
		rSubject.SetMapIndex(rKey, getValue(oldValue))
		m.Unlock()

	case reflect.Slice, reflect.Array:
		if rKey.IsValid() {
			if rKey.Kind() != reflect.Int {
				return errors.New("invalid index type - should be integer")
			}
			idx := int(rKey.Int())
			if rSubject.Len() <= idx {
				return errors.New("index out of bound")
			}

			m.Lock()
			oldValue := rSubject.Index(idx)
			oldValue.Set(getValue(oldValue))
			m.Unlock()
		} else {
			m.Lock()
			reflect.Append(rSubject, getValue(rSubject))
			m.Unlock()
		}

	case reflect.Ptr, reflect.Chan, reflect.Func:
		m.Lock()
		rSubject.Set(getValue(rValue))
		m.Unlock()

	default:
		return errors.New("invalid acceptor, must be pointer and not scalar value")
	}

	return err
}

func SyncGet(subject interface{}, path ...interface{}) (interface{}, error) {



	return nil, nil
}