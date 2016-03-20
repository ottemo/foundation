package utils

import (
	"sync"
	"reflect"
	"crypto/md5"
	"fmt"
)

var (
	locks = make(map[string]*syncMutex)
	locksMutex sync.Mutex
)

type syncMutex struct {
	index string
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

func (it *syncMutex) GetIndex() string {
	return it.index
}

func (it *syncMutex) IsLocked() bool {
	return it.lock
}

func (it *syncMutex) Refs() int {
	return it.refs
}

func getMutexIndex(combination []interface{}) string {
	index := ""
	for _, item := range combination {
		value := reflect.ValueOf(item)

		var indexElement string = ""
		switch value.Kind() {
		case reflect.String:
			indexElement += value.String()

		case reflect.Chan,
			reflect.Map,
			reflect.Ptr,
			reflect.UnsafePointer,
			reflect.Func,
			reflect.Slice,
			reflect.Array:

			indexElement = string(value.Pointer())
			// indexElement = fmt.Sprintf("%p", item)

		default:
			indexElement = fmt.Sprintf("%v", item)
		}

		index += "/" + indexElement
	}

	if len(index) > 32 {
		index = string(md5.New().Sum([]byte(index)))
	}

	return index
}

func GetMutex(combination ...interface{}) *syncMutex {
	locksMutex.Lock()
	defer locksMutex.Unlock()

	index := getMutexIndex(combination)

	mutex, present := locks[index]
	if !present {
		mutex = new(syncMutex)
		mutex.index = index
		locks[index] = mutex
	}
	mutex.refs++
	return mutex
}
