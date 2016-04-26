package utils

import (
	"errors"
	"reflect"
	"sync"
)

var (
	locks      = make(map[uintptr]*syncMutex)
	locksMutex sync.Mutex
)

// syncMutex an extended variation of sync.Mutex
type syncMutex struct {
	index uintptr
	refs  int
	lock  bool
	mutex sync.Mutex
}

// Lock holds the lock on a mutex
func (it *syncMutex) Lock() {
	it.mutex.Lock()
	it.lock = true
}

// Unlock releases the mutex lock
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

// GetIndex returns the pointer of an holden subject (the index in global mutexes map)
func (it *syncMutex) GetIndex() uintptr {
	return it.index
}

// IsLocked checks if the mutex if currently locked without actual lock
func (it *syncMutex) IsLocked() bool {
	return it.lock
}

// Refs returns the amount of references to a mutex
func (it *syncMutex) Refs() int {
	return it.refs
}

// GetPointer returns the memory address value for a given subject or an error for scalar values
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

	if value.CanAddr() {
		value = value.Addr()
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

	// debug.PrintStack()
	return 0, errors.New("can't get pointer to " + value.Type().String())
}

// SyncMutex creates a mutex element for a given subject or error for a scalar values (un-addressable values)
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

// SyncLock holds mutex on a given subject
func SyncLock(subject interface{}) error {
	mutex, err := SyncMutex(subject)
	if err != nil {
		return err
	}
	mutex.Lock()
	return nil
}

// SyncUnlock releases mutex on a given subject
func SyncUnlock(subject interface{}) error {
	mutex, err := SyncMutex(subject)
	if err != nil {
		return err
	}
	mutex.Unlock()
	return nil
}

// pathItem represents the stack of (index, value, mutex) collected during path-walk on a tree like type value
type pathItem struct {
	parent *pathItem
	mutex  *syncMutex
	locked bool
	key    reflect.Value
	value  reflect.Value
}

// Lock is a SyncMutex.Lock() implementation for a pathItem
// (it uses own it.locked flag to safely handle Lock/Unlock only for a current instance)
func (it *pathItem) Lock() {
	if it.mutex != nil && !it.locked {
		it.locked = true
		it.mutex.Lock()
	}
}

// Unlock is a SyncMutex.Lock() implementation for a pathItem
// (it uses own it.locked flag to safely handle Lock/Unlock only for a current instance)
func (it *pathItem) Unlock() {
	if it.mutex != nil && it.locked {
		it.mutex.Unlock()
		it.locked = false
	}
}

// LockStack holds the lock for a whole path item stack including self
func (it *pathItem) LockStack() {
	for x := it; x != nil; x = x.parent {
		x.Lock()
	}
}

// UnlockStack releases the locks for a whole path item stack including self
func (it *pathItem) UnlockStack() {
	for x := it; x != nil; x = x.parent {
		x.Unlock()
	}
}

// Update refreshes the parent items information. The function going through the stack and updates the values
// (such functionality is required for the the unlocked slice parents, so while the new slice item created ahed of slice capacity,
// the copy of slice is made -with such case the parent items information will address unused memory and should be updated )
func (it *pathItem) Update(newSubject reflect.Value) {
	stack := make([]*pathItem, 0, 100)
	wasLocked := make([]*syncMutex, 0, 100)

	// collecting pathItem stack and locking items
	for x := it.parent; x != nil; x = x.parent {
		x.mutex.Lock()
		stack = append(stack, x)
		wasLocked = append(wasLocked, x.mutex)
	}

	// updating pathItem references
	for i := len(stack) - 1; i > 0; i-- {
		oldValue := stack[i-1].value

		switch stack[i].value.Kind() {
		case reflect.Map:
			stack[i-1].value = stack[i].value.MapIndex(stack[i].key)

		case reflect.Slice, reflect.Array:
			idx := int(stack[i].key.Int())
			stack[i-1].value = stack[i].value.Index(idx)
		}

		if oldValue != stack[i-1].value {
			mutex, err := SyncMutex(stack[i-1].value)
			if err != nil {
				panic(err)
			}
			mutex.Lock()

			wasLocked[i-1].Unlock()
			wasLocked[i-1] = mutex

			stack[i-1].locked = true
			stack[i-1].mutex = mutex
			stack[i-1].parent = stack[i]
		}
	}

	// updating the item with new key value
	switch it.value.Kind() {
	case reflect.Map:
		it.value.SetMapIndex(it.key, newSubject)

	case reflect.Slice, reflect.Array:
		idx := int(it.key.Int())
		it.value.Index(idx).Set(newSubject)
	}

	// un-locking locked items
	for _, x := range wasLocked {
		x.Unlock()
	}
}

// SyncSet - synchronized write access to tree like variables
//   - the value could be a func(oldValue {type}) {type} which would be synchronized called
//   - the -1 index for a slice means it's extensions for a new element
func SyncSet(subject interface{}, value interface{}, path ...interface{}) error {

	pathItem, err := getPathItem(subject, path, true, 2, nil)
	if err != nil {
		pathItem.UnlockStack()
		return err
	}

	rSubject := pathItem.value
	if !rSubject.IsValid() {
		pathItem.UnlockStack()
		return errors.New("invalid subject")
	}

	// kind := rSubject.Kind()
	// if kind != reflect.Ptr && pathItem.parent != nil {
	if pathItem.mutex == nil && pathItem.parent != nil {
		pathItem.Unlock()
		pathItem = pathItem.parent
		rSubject = pathItem.value
	}

	rSubject = reflect.Indirect(rSubject)
	rKey := pathItem.key

	// new value validation
	rValue := reflect.ValueOf(value)
	rValueType := rValue.Type()

	// allowing to have setter function instead of just value
	funcValue := func(oldValue reflect.Value, valueType reflect.Type) reflect.Value {
		if !oldValue.IsValid() {
			oldValue = reflect.New(valueType).Elem()
		}
		if rValue.Kind() == reflect.Func {
			if rValueType.NumOut() == 1 && rValueType.NumIn() == 1 {
				// oldValueType := oldValue.Type()
				// !rValueType.In(0).AssignableTo(oldValueType) &&
				//!rValueType.Out(0).AssignableTo(oldValueType) {
				return rValue.Call([]reflect.Value{oldValue})[0]
			}
		}
		return rValue
	}

	switch rSubject.Kind() {
	case reflect.Map:
		oldValue := rSubject.MapIndex(rKey)
		rSubject.SetMapIndex(rKey, funcValue(oldValue, oldValue.Type()))

	case reflect.Slice, reflect.Array:
		idx := int(rKey.Int())
		oldValue := rSubject.Index(idx)
		oldValue.Set(funcValue(oldValue, oldValue.Type()))

	default:
		rSubject.Set(funcValue(rSubject, rSubject.Type()))
	}

	pathItem.UnlockStack()
	return nil
}

// SyncGet - synchronized read access tree like variables
// 	- the -1 index for a slice means to init new blank slice value
//
// sample:
// 	var x map[string]map[int][]int = ...
//
// 	// un-synchronized access to value - causes simultaneous read/write access panics for parallel threade
// 	y := x["a"][5][1]
//
//	// synchronized access to variable
//	y := SyncGet(x, false, "a", 5, 1) // returns x["a"][5][1], or (nil, error) if item does not exist
//      y := SyncGet(x, true, "a", 5, -1) // returns new blank value of list for a possibly new element x["a"][5]
func SyncGet(subject interface{}, initBlank bool, path ...interface{}) (interface{}, error) {
	result, err := getPathItem(subject, path, initBlank, 0, nil)
	if err != nil {
		return nil, err
	}
	return result.value.Interface(), nil
}

// initBlankValue initializes a new blankValue for a given type
func initBlankValue(valueType reflect.Type) (reflect.Value, error) {
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

// getPathItem - the shared function used by SyncSet, SyncGet to synchronized access to tree like type items
// 	- the "initBlank" argument controls new element creation if index does not exists
//	- the "lockLevel" controls the amount of path items should be locked at the end
//		-1 - hold lock on whole path items
// 		 0 - hold lock to nothing
// 		 1 - hold lock to result item only
// 		 2 - hold lock to result item and it's parent item
//		 3 - ...
//
//		if "lockLevel" != 0, the locks should be unlocked by caller with pathItem.Unlock() call
// 		or routines like:
//		  level := 1
//		  for x := pathItem; lockLevel >0 && level <= lockLevel && x.parent != nil; x = x.parent {
//		      x.mutex.Unlock()
//		      level++
//		  }
//	- the "parent" element is used for a recursion, should be nil for an initial call
func getPathItem(subject interface{}, path []interface{}, initBlank bool, lockLevel int, parent *pathItem) (*pathItem, error) {
	var err error

	// do nothing for nil objects
	if subject == nil {
		return nil, errors.New("nil subject")
	}

	// checking for reflect.Value in subject
	rSubject, ok := subject.(reflect.Value)
	if !ok {
		rSubject = reflect.ValueOf(subject)
	}

	// handling pointers
	rSubject = reflect.Indirect(rSubject)

	// checking subject for a zero value
	if !rSubject.IsValid() {
		return nil, errors.New("invalid subject")
	}

	// taking subject type and kind
	rSubjectKind := rSubject.Kind()
	rSubjectType := rSubject.Type()

	// initializing result item
	pathItem := &pathItem{
		parent: parent,
		value:  rSubject,
		mutex:  nil,
	}

	if len(path) == 0 {
		return pathItem, nil
	}

	// taking mutex for subject if possible
	mutex, err := SyncMutex(rSubject)
	if err != nil {
		if len(path) != 0 {
			return nil, err
		}
		mutex = nil
	}
	pathItem.mutex = mutex

	// locking access to subject
	pathItem.Lock()

	// we have optional unlock (see "lockLevel" argument)
	// so this function should decide where to unlock
	unlock := func() {
		if lockLevel == -1 {
			return
		}

		if len(path)+1 > lockLevel {
			pathItem.Unlock()
		}
	}

	// checking for end of path, if so we are done
	if len(path) == 0 {
		unlock()
		return pathItem, nil
	}

	// otherwise the first item of a path as current key item
	rKey := reflect.ValueOf(path[0])
	if !rKey.IsValid() {
		return nil, errors.New("invalid path")
	}
	rKeyType := rKey.Type()
	pathItem.key = rKey

	newPath := path[1:]

	// getting path item from subject based on it's type
	switch rSubjectKind {

	case reflect.Map:
		// comparing given path key type to subject key type
		if rKeyType != rSubjectType.Key() {
			unlock()
			return nil, errors.New("invalid path item type " +
				rKeyType.String() + " != " + rSubjectType.Key().String())
		}

		// taking subject key item
		rSubjectItem := rSubject.MapIndex(rKey)

		// checking if item is not defined and we should make new value
		if !rSubjectItem.IsValid() && initBlank {
			if rSubjectItem, err = initBlankValue(rSubjectType.Elem()); err != nil {
				unlock()
				return nil, err
			}
			rSubject.SetMapIndex(rKey, rSubjectItem)
		}

		unlock()
		return getPathItem(rSubjectItem, newPath, initBlank, lockLevel, pathItem)

	case reflect.Slice, reflect.Array:
		// for the slices the key should be integer index
		if rKey.Kind() != reflect.Int {
			unlock()
			return nil, errors.New("invalid path item type: " +
				rKey.Kind().String() + " != Int")
		}
		idx := int(rKey.Int())

		// checking the length of slice
		if rSubject.Len() <= idx {
			unlock()
			return nil, errors.New("index " + rKey.String() + " is out of bound")
		}

		// (idx = -1) is used to create new item, otherwise it is existing item
		if idx >= 0 {
			rSubjectItem := rSubject.Index(idx)

			// checking if existing item is nil but we should initialize it
			if !rSubjectItem.IsValid() && initBlank {
				rSubjectValue, err := initBlankValue(rSubjectType.Elem())
				if err != nil {
					unlock()
					return nil, err
				}
				rSubjectItem.Set(rSubjectValue)
			}

			rSubject = rSubjectItem
		} else {
			// checking if new item creation was specified, and we can create it
			if !initBlank {
				unlock()
				return nil, errors.New("invalid index -1 as initBlank = false")
			}

			// checking subject value to be a reference
			if !rSubject.CanAddr() {
				unlock()
				return nil, errors.New("not addresable subject")
			}

			// initializing new blank item
			newItem, err := initBlankValue(rSubjectType.Elem())
			if err != nil {
				unlock()
				return nil, err
			}

			// checking if slice capacity allows to increase length
			length := rSubject.Len()
			if rSubject.Cap() < length {
				rSubject.SetLen(length + 1)
				rSubject.Index(length).Set(newItem)
			} else {
				// new slice creation required (worst scenario)
				newSubject := reflect.New(rSubjectType).Elem()
				newSubject.Set(reflect.Append(rSubject, newItem))
				rSubject.Set(newSubject)

				if parent != nil {
					if lockLevel != -1 {
						parent.Update(newSubject)
					}
					pathItem.parent = parent
				}
			}

			// updating path item info
			pathItem.value = rSubject
			pathItem.key = reflect.ValueOf(length)
			rSubject = rSubject.Index(length)
		}

		unlock()
		return getPathItem(rSubject, newPath, initBlank, lockLevel, pathItem)

	default:
		unlock()
		return nil, errors.New("invalid subject, path can not be applied to: " + rSubjectType.String())
	}

	return pathItem, nil
}
