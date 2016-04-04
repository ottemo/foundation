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
	funcValue := func(oldValue reflect.Value, valueType reflect.Type) reflect.Value {
		if !oldValue.IsValid() {
			oldValue = reflect.New(valueType).Elem()
		}
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
		rSubject.SetMapIndex(rKey, funcValue(oldValue, rSubjectType.Elem()))
		m.Unlock()

	case reflect.Slice, reflect.Array:
		// checking if path was not specified
		if !rKey.IsValid() {
			// no path key was specified - i.e. assigning value should be slice/array
			if !rSubject.CanAddr() {
				return errors.New("unadressable sbject")
			}

			m.Lock()
			newValue := funcValue(rSubject, rSubjectType)
			if newValue.Type().AssignableTo(rSubjectType) {
				rSubject.Set(newValue)
			}
			m.Unlock()
		} else {
			// path key was specified - i.e. slice item modification or adding new
			if rKey.Kind() != reflect.Int {
				return errors.New("invalid index type - should be integer")
			}

			idx := int(rKey.Int())
			if rSubject.Len() <= idx {
				return errors.New("index out of bound (SyncSet)")
			}

			// (idx = -1) is a condition to create new item
			if idx >= 0 {
				// changing existing item
				m.Lock()
				oldValue := rSubject.Index(idx)
				oldValue.Set(funcValue(oldValue, rSubjectType.Elem()))
				m.Unlock()
			} else {
				// making new item
				if !rSubject.CanAddr() {
					return errors.New("invalid acceptor: " + rSubjectType.String() )
				}

				m.Lock()
				oldValue := reflect.ValueOf(nil)
				newValue := funcValue(oldValue, rSubjectType.Elem())

				length := rSubject.Len()
				if rSubject.Cap() < length {
					rSubject.SetLen(length + 1)
					rSubject.Index(length).Set(newValue)
				} else {
					rSubject.Set(reflect.Append(rSubject, newValue))
				}
				m.Unlock()
			}
		}

	case reflect.Ptr, reflect.Chan, reflect.Func:
		m.Lock()
		rSubject.Set(funcValue(rValue, rSubjectType))
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
			value := reflect.New(valueType).Elem()
			value.Set(reflect.MakeSlice(valueType, 0, 10))
			return value, nil
		case reflect.Chan, reflect.Func:
			break
		default:
			return reflect.New(valueType).Elem(), nil
		}
		return reflect.ValueOf(nil), errors.New("unsuported blank value type " + valueType.String())
	}

	var rKey reflect.Value
	var rKeyType reflect.Type

	// If path is specified then taking path element as subject
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

		// taking path key type
		rKey = reflect.ValueOf(key)
		rKeyType = rKey.Type()

		switch rSubjectKind {
		case reflect.Map:
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

		// (idx = -1) is a condition to create new item
		case reflect.Slice, reflect.Array:

			// key should be index
			if rKey.Kind() != reflect.Int {
				return nil, errors.New("invalid slice/array index type (" + rKey.Kind().String() + "), should be integer")
			}

			idx := int(rKey.Int())
			if rSubject.Len() <= idx {
				return nil, errors.New("index out of bound (SyncGet)")
			}

			// taking mutex for item
			m, err := SyncMutex(rSubject)
			if err != nil {
				return nil, err
			}


			if idx >= 0 {
				// access to existing item (time critical)
				m.Lock()
				rSubject = rSubject.Index(idx).Addr() //?
				if !rSubject.IsValid() && initBlank {
					// item value is nil
					rSubjectValue, err := initBlankValue(rSubjectType.Elem())
					if err != nil {
						m.Unlock()
						return nil, err
					}
					rSubject.Set(rSubjectValue)
				}
				m.Unlock()
			} else {
				// making new item
				if (!initBlank) {
					return nil, errors.New("invalid index -1 as initBlank = false")
				}

				if !rSubject.CanAddr() {
					return nil, errors.New("not addresable subject")
				}

				// checking if capacity allows to increase slice/array length
				m.Lock()
				newItemValue, err := initBlankValue(rSubjectType.Elem())
				if err != nil {
					return nil, err
				}

				length := rSubject.Len()
				if rSubject.Cap() < length {
					rSubject.SetLen(length + 1)
					rSubject.Index(length).Set(newItemValue)
					rSubject = rSubject.Index(length).Addr()

				} else {
					rSubject.Set(reflect.Append(rSubject, newItemValue))
				}
				rSubject = rSubject.Index(length).Addr()
				m.Unlock()
			}
			rSubjectType = rSubject.Type()

		default:
			return nil, errors.New(fmt.Sprintf("invalid element type on path item %d", idx))
		}
	}

	return rSubject.Interface(), nil
}