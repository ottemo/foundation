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

	var err error
	var rKey reflect.Value
	var rKeyType reflect.Type

	// if the path is specified - taking element to work with
	pathLen := len(path)
	if len(path) > 1 {
		subject, err = SyncGet(subject, true, path[:pathLen-1]...)
		if err != nil {
			return err
		}

		rKey = reflect.ValueOf( path[pathLen-1] )
		rKeyType = rKey.Type()
	}

	// checking subject
	if subject == nil {
		return errors.New("subject is nil")
	}
	rSubject := reflect.ValueOf(subject)

	// taking mutex to subject
	m, err := SyncMutex(rSubject)
	if err != nil {
		return err
	}

	// if pointer is given taking reflected element to work with
	rSubjectKind := rSubject.Kind()
	if rSubjectKind == reflect.Ptr || rSubjectKind == reflect.Interface {
		rSubject = rSubject.Elem()
		rSubjectKind = rSubject.Kind()
	}
	rSubjectType := rSubject.Type()

	// set value validation
	rValue := reflect.ValueOf(value)
	rValueType := rValue.Type()

	// allowing to have setter function instead of just value
	funcValue := func(oldValue reflect.Value) reflect.Value {
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

	// setting the subject value according it's type
	switch rSubject.Kind() {
	case reflect.Map:
		if rKeyType != rSubjectType.Key() {
			return errors.New("invalid map key type (" + rKeyType.String() + " != " + rSubjectType.Key().String() + ")")
		}

		m.Lock()
		oldValue := rSubject.MapIndex(rKey)
		rSubject.SetMapIndex(rKey, funcValue(oldValue))
		m.Unlock()

	case reflect.Slice, reflect.Array:
		if rKey.IsValid() {
			// slice index was specified
			if rKey.Kind() != reflect.Int {
				return errors.New("invalid index type - should be integer")
			}
			idx := int(rKey.Int())
			if rSubject.Len() <= idx {
				return errors.New("index out of bound")
			}

			m.Lock()
			oldValue := rSubject.Index(idx)
			oldValue.Set(funcValue(oldValue))
			m.Unlock()
		} else {
			// the new element supposed
			return errors.New("not implemented")
		}

	case reflect.Ptr, reflect.Chan, reflect.Func:
		m.Lock()
		rSubject.Set(funcValue(rValue))
		m.Unlock()

	default:
		return errors.New("invalid acceptor, must be pointer and not scalar value")
	}

	return err
}

func SyncGet(subject interface{}, initBlank bool, path ...interface{}) (interface{}, error) {

	if subject == nil {
		return nil, errors.New("nil subject")
	}

	rSubject := reflect.ValueOf(subject)
	var rSubjectType reflect.Type
	var rSubjectKind reflect.Kind

	// function to make a blankValue based on given type
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

	var rKey reflect.Value
	var rKeyType reflect.Type

	for idx, key := range path {

		if !rSubject.IsValid() {
			return nil, errors.New("invalid path")
		}

		rSubjectKind = rSubject.Kind()
		if rSubjectKind == reflect.Ptr || rSubjectKind == reflect.Interface {
			rSubject = rSubject.Elem()
			rSubjectKind = rSubject.Kind()
		}
		rSubjectType = rSubject.Type()

		switch rSubjectKind {
		case reflect.Map:
			rKey = reflect.ValueOf(key)
			rKeyType = rKey.Type()
			if rKeyType != rSubjectType.Key() {
				return nil, errors.New(fmt.Sprintf("invalid type of path item %d", idx))
			}

			// taking mutex for item
			m, err := SyncMutex(rSubject)
			if err != nil {
				return nil, err
			}

			// time critical access to element
			m.Lock()
			rSubjectValue := rSubject.MapIndex(rKey)
			if !rSubjectValue.IsValid() && initBlank {
				if rSubjectValue, err = initBlankValue(rSubjectType.Elem()); err != nil {
					m.Unlock()
					return nil, err
				}
				rSubject.SetMapIndex(rKey, rSubjectValue)
			}
			m.Unlock()

			rSubject = rSubjectValue
			rSubjectType = rSubject.Type()

		case reflect.Slice, reflect.Array:
			// key should be index
			if rKey.Kind() != reflect.Int {
				return errors.New("invalid index type - should be integer")
			}
			idx := int(rKey.Int())
			if rSubject.Len() <= idx {
				return errors.New("index out of bound")
			}

			// taking mutex for item
			m, err := SyncMutex(rSubject)
			if err != nil {
				return nil, err
			}

			// time critical access to element
			m.Lock()
			rSubject := rSubject.Index(idx)
			if !rSubject.IsValid() && initBlank {
				if rSubjectValue, err = initBlankValue(rSubjectType.Elem()); err != nil {
					m.Unlock()
					return nil, err
				}
				rSubject.Set(rSubjectValue)
			}
			m.Unlock()

			rSubjectType = rSubject.Type()

		default:
			return nil, errors.New(fmt.Sprintf("invalid element type on path item %d", idx))
		}
	}

	return rSubject.Interface(), nil
}