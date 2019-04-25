package context

import (
	"sync"
	"reflect"
	"runtime"
)

// proxies holds set of pass-through functions
var proxies = [...] func(target func()) {
	func(target func()) { target() }, // 00
	func(target func()) { target() }, // 01
	func(target func()) { target() }, // 02
	func(target func()) { target() }, // 03
	func(target func()) { target() }, // 04
	func(target func()) { target() }, // 05
	func(target func()) { target() }, // 06
	func(target func()) { target() }, // 07
	func(target func()) { target() }, // 08
	func(target func()) { target() }, // 09
	func(target func()) { target() }, // 0a
	func(target func()) { target() }, // 0b
	func(target func()) { target() }, // 0c
	func(target func()) { target() }, // 0d
	func(target func()) { target() }, // 0e
	func(target func()) { target() }, // 0f
	func(target func()) { target() }, // 10
	func(target func()) { target() }, // 11
	func(target func()) { target() }, // 12
	func(target func()) { target() }, // 13
	func(target func()) { target() }, // 14
	func(target func()) { target() }, // 15
	func(target func()) { target() }, // 16
	func(target func()) { target() }, // 17
	func(target func()) { target() }, // 18
	func(target func()) { target() }, // 19
	func(target func()) { target() }, // 1a
	func(target func()) { target() }, // 1b
	func(target func()) { target() }, // 1c
	func(target func()) { target() }, // 1d
	func(target func()) { target() }, // 1e
	func(target func()) { target() }, // 1f
	func(target func()) { target() }, // 20
	func(target func()) { target() }, // 21
	func(target func()) { target() }, // 22
	func(target func()) { target() }, // 23
	func(target func()) { target() }, // 24
	func(target func()) { target() }, // 25
	func(target func()) { target() }, // 26
	func(target func()) { target() }, // 27
	func(target func()) { target() }, // 28
	func(target func()) { target() }, // 29
	func(target func()) { target() }, // 2a
	func(target func()) { target() }, // 2b
	func(target func()) { target() }, // 2c
	func(target func()) { target() }, // 2d
	func(target func()) { target() }, // 2e
	func(target func()) { target() }, // 2f
	func(target func()) { target() }, // 30
	func(target func()) { target() }, // 31
	func(target func()) { target() }, // 32
	func(target func()) { target() }, // 33
	func(target func()) { target() }, // 34
	func(target func()) { target() }, // 35
	func(target func()) { target() }, // 36
	func(target func()) { target() }, // 37
	func(target func()) { target() }, // 38
	func(target func()) { target() }, // 39
	func(target func()) { target() }, // 3a
	func(target func()) { target() }, // 3b
	func(target func()) { target() }, // 3c
	func(target func()) { target() }, // 3d
	func(target func()) { target() }, // 3e
	func(target func()) { target() }, // 3f
	func(target func()) { target() }, // 40
	func(target func()) { target() }, // 41
	func(target func()) { target() }, // 42
	func(target func()) { target() }, // 43
	func(target func()) { target() }, // 44
	func(target func()) { target() }, // 45
	func(target func()) { target() }, // 46
	func(target func()) { target() }, // 47
	func(target func()) { target() }, // 48
	func(target func()) { target() }, // 49
	func(target func()) { target() }, // 4a
	func(target func()) { target() }, // 4b
	func(target func()) { target() }, // 4c
	func(target func()) { target() }, // 4d
	func(target func()) { target() }, // 4e
	func(target func()) { target() }, // 4f
	func(target func()) { target() }, // 50
	func(target func()) { target() }, // 51
	func(target func()) { target() }, // 52
	func(target func()) { target() }, // 53
	func(target func()) { target() }, // 54
	func(target func()) { target() }, // 55
	func(target func()) { target() }, // 56
	func(target func()) { target() }, // 57
	func(target func()) { target() }, // 58
	func(target func()) { target() }, // 59
	func(target func()) { target() }, // 5a
	func(target func()) { target() }, // 5b
	func(target func()) { target() }, // 5c
	func(target func()) { target() }, // 5d
	func(target func()) { target() }, // 5e
	func(target func()) { target() }, // 5f
	func(target func()) { target() }, // 60
	func(target func()) { target() }, // 61
	func(target func()) { target() }, // 62
	func(target func()) { target() }, // 63
	func(target func()) { target() }, // 64
	func(target func()) { target() }, // 65
	func(target func()) { target() }, // 66
	func(target func()) { target() }, // 67
	func(target func()) { target() }, // 68
	func(target func()) { target() }, // 69
	func(target func()) { target() }, // 6a
	func(target func()) { target() }, // 6b
	func(target func()) { target() }, // 6c
	func(target func()) { target() }, // 6d
	func(target func()) { target() }, // 6e
	func(target func()) { target() }, // 6f
	func(target func()) { target() }, // 70
	func(target func()) { target() }, // 71
	func(target func()) { target() }, // 72
	func(target func()) { target() }, // 73
	func(target func()) { target() }, // 74
	func(target func()) { target() }, // 75
	func(target func()) { target() }, // 76
	func(target func()) { target() }, // 77
	func(target func()) { target() }, // 78
	func(target func()) { target() }, // 79
	func(target func()) { target() }, // 7a
	func(target func()) { target() }, // 7b
	func(target func()) { target() }, // 7c
	func(target func()) { target() }, // 7d
	func(target func()) { target() }, // 7e
	func(target func()) { target() }, // 7f
	func(target func()) { target() }, // 80
	func(target func()) { target() }, // 81
	func(target func()) { target() }, // 82
	func(target func()) { target() }, // 83
	func(target func()) { target() }, // 84
	func(target func()) { target() }, // 85
	func(target func()) { target() }, // 86
	func(target func()) { target() }, // 87
	func(target func()) { target() }, // 88
	func(target func()) { target() }, // 89
	func(target func()) { target() }, // 8a
	func(target func()) { target() }, // 8b
	func(target func()) { target() }, // 8c
	func(target func()) { target() }, // 8d
	func(target func()) { target() }, // 8e
	func(target func()) { target() }, // 8f
	func(target func()) { target() }, // 90
	func(target func()) { target() }, // 91
	func(target func()) { target() }, // 92
	func(target func()) { target() }, // 93
	func(target func()) { target() }, // 94
	func(target func()) { target() }, // 95
	func(target func()) { target() }, // 96
	func(target func()) { target() }, // 97
	func(target func()) { target() }, // 98
	func(target func()) { target() }, // 99
	func(target func()) { target() }, // 9a
	func(target func()) { target() }, // 9b
	func(target func()) { target() }, // 9c
	func(target func()) { target() }, // 9d
	func(target func()) { target() }, // 9e
	func(target func()) { target() }, // 9f
	func(target func()) { target() }, // a0
	func(target func()) { target() }, // a1
	func(target func()) { target() }, // a2
	func(target func()) { target() }, // a3
	func(target func()) { target() }, // a4
	func(target func()) { target() }, // a5
	func(target func()) { target() }, // a6
	func(target func()) { target() }, // a7
	func(target func()) { target() }, // a8
	func(target func()) { target() }, // a9
	func(target func()) { target() }, // aa
	func(target func()) { target() }, // ab
	func(target func()) { target() }, // ac
	func(target func()) { target() }, // ad
	func(target func()) { target() }, // ae
	func(target func()) { target() }, // af
	func(target func()) { target() }, // b0
	func(target func()) { target() }, // b1
	func(target func()) { target() }, // b2
	func(target func()) { target() }, // b3
	func(target func()) { target() }, // b4
	func(target func()) { target() }, // b5
	func(target func()) { target() }, // b6
	func(target func()) { target() }, // b7
	func(target func()) { target() }, // b8
	func(target func()) { target() }, // b9
	func(target func()) { target() }, // ba
	func(target func()) { target() }, // bb
	func(target func()) { target() }, // bc
	func(target func()) { target() }, // bd
	func(target func()) { target() }, // be
	func(target func()) { target() }, // bf
	func(target func()) { target() }, // c0
	func(target func()) { target() }, // c1
	func(target func()) { target() }, // c2
	func(target func()) { target() }, // c3
	func(target func()) { target() }, // c4
	func(target func()) { target() }, // c5
	func(target func()) { target() }, // c6
	func(target func()) { target() }, // c7
	func(target func()) { target() }, // c8
	func(target func()) { target() }, // c9
	func(target func()) { target() }, // ca
	func(target func()) { target() }, // cb
	func(target func()) { target() }, // cc
	func(target func()) { target() }, // cd
	func(target func()) { target() }, // ce
	func(target func()) { target() }, // cf
	func(target func()) { target() }, // d0
	func(target func()) { target() }, // d1
	func(target func()) { target() }, // d2
	func(target func()) { target() }, // d3
	func(target func()) { target() }, // d4
	func(target func()) { target() }, // d5
	func(target func()) { target() }, // d6
	func(target func()) { target() }, // d7
	func(target func()) { target() }, // d8
	func(target func()) { target() }, // d9
	func(target func()) { target() }, // da
	func(target func()) { target() }, // db
	func(target func()) { target() }, // dc
	func(target func()) { target() }, // dd
	func(target func()) { target() }, // de
	func(target func()) { target() }, // df
	func(target func()) { target() }, // e0
	func(target func()) { target() }, // e1
	func(target func()) { target() }, // e2
	func(target func()) { target() }, // e3
	func(target func()) { target() }, // e4
	func(target func()) { target() }, // e5
	func(target func()) { target() }, // e6
	func(target func()) { target() }, // e7
	func(target func()) { target() }, // e8
	func(target func()) { target() }, // e9
	func(target func()) { target() }, // ea
	func(target func()) { target() }, // eb
	func(target func()) { target() }, // ec
	func(target func()) { target() }, // ed
	func(target func()) { target() }, // ee
	func(target func()) { target() }, // ef
	func(target func()) { target() }, // f0
	func(target func()) { target() }, // f1
	func(target func()) { target() }, // f2
	func(target func()) { target() }, // f3
	func(target func()) { target() }, // f4
	func(target func()) { target() }, // f5
	func(target func()) { target() }, // f6
	func(target func()) { target() }, // f7
	func(target func()) { target() }, // f8
	func(target func()) { target() }, // f9
	func(target func()) { target() }, // fa
	func(target func()) { target() }, // fb
	func(target func()) { target() }, // fc
	func(target func()) { target() }, // fd
	func(target func()) { target() }, // fe
	func(target func()) { target() }, // ff
}

// proxiesBase holds the length of proxies[] - it is context ID digit encoding base (0xff per function call in stack)
var proxiesBase uint = uint(len(proxies))

// proxiesDict is a decoding table for proxy functions (proxy function call in stack trace have appropriate numeric value)
var proxiesDict map[uintptr]uint

// proxiesStart the pointer to a function which starts the proxy calls
var proxiesStart uintptr

// contexts are the key-values maps which are accessible in any nested function after context creation call
var contexts = make(map[uint]map[string]interface{})

// contextsMutex synchronizes the access to contexts variable
var contextsMutex sync.RWMutex

// init performs the package self-initialization routine
func init() {
	proxiesStart = reflect.ValueOf(RunInContext).Pointer()

	proxiesDict = make(map[uintptr]uint)
	for idx, val := range proxies {
		proxiesDict[reflect.ValueOf(val).Pointer()] = uint(idx)
	}
	proxiesDict[reflect.ValueOf(proxyLoop).Pointer()] = proxiesBase
}

// proxyLoop is the service recursive function for context ID value encoding
func proxyLoop(target func (), i uint) {
	i = i - proxiesBase
	if i >= proxiesBase {
		proxyLoop(target, i)
	} else {
		proxies[i](target)
	}
}

// getCallStack returns array of current call stack function entries pointers
func getCallStack(skip int) []uintptr {
	pcSize := 100
	pc := make([]uintptr, pcSize)

	n := runtime.Callers(skip, pc)
	for n >= pcSize {
		pcSize = pcSize * 10
		pc = make([]uintptr, pcSize)
		runtime.Callers(skip, pc)
	}

	result := make([]uintptr, 0, n)
	frames := runtime.CallersFrames(pc)
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		result = append(result, frame.Entry)
		// prints the call stack - for debugging purposes
		// fmt.Printf("%d:%d - %s:%d - %s\n", frame.PC, frame.Entry, frame.File, frame.Line, frame.Function)
	}

	return result
}

// GetContextId returns current identifier (positive number) or 0 if it does not exists
func GetContextId() uint {
	var contextsId uint

	pointers := getCallStack(2)
	for idx, ptr := range pointers {
		if ptr == proxiesStart {
			for idx--; idx > 0; idx-- {
				if value, present := proxiesDict[pointers[idx]]; present {
					contextsId += value
				} else {
					return contextsId + 1
				}
			}
		}
	}
	return contextsId
}

// GetContextById returns context map by a given identifier or nil if it does not exists
func GetContextById(id uint) map[string]interface{}  {
	var result map[string]interface{} = nil

	contextsMutex.Lock()
	if value, present := contexts[id-1]; present {
		result = value
	}
	contextsMutex.Unlock()

	return result
}

// RunInContext executes given function within a new or given a "context"
// i.e. context map will be accessible with the GetContext() call within given function or its sub-calls
func RunInContext(target func(), context map[string]interface{}) map[string]interface{} {
	if context == nil {
		context = make(map[string]interface{})
	}

	contextsMutex.Lock()
	var contextsIdx uint
	for ; true; contextsIdx++ {
		if _, present := contexts[contextsIdx]; !present {
			break
		}
	}

	contexts[contextsIdx] = context
	contextsMutex.Unlock()

	defer func() {
		contextsMutex.Lock()
		delete(contexts, contextsIdx)
		contextsMutex.Unlock()
	}()

	if contextsIdx >= proxiesBase {
		proxyLoop(target, contextsIdx)
	} else {
		proxies[contextsIdx](target)
	}
	return context
}

// GetContextValue returns the value in current context or nil if context or value does not exists
func GetContextValue(key string) interface{} {
	contextId := GetContextId() - 1

	contextsMutex.Lock()
	defer func () {
		contextsMutex.Unlock()
	} ()

	if context, present := contexts[contextId]; present {
		if value, present := context[key]; present {
			return value
		}
	}
	return nil
}

// GetContext returns the context map associated to a current stack or nil of there are no context available
func SetContextValue(key string, value interface{}) interface{} {
	contextId := GetContextId() - 1
	var result interface{}

	contextsMutex.Lock()
	defer func () {
		contextsMutex.Unlock()
	} ()

	if context, present := contexts[contextId]; present {
		if oldValue, present := context[key]; present {
			result = oldValue
		}
		context[key] = value
	}
	return result
}

// GetContext returns the context map associated to a current stack or nil of there are no context available
func GetContext() map[string]interface{} {
	return GetContextById(GetContextId())
}

// MakeContext executes given function within a new context (alias for RunInContext)
func MakeContext(target func()) {
	RunInContext(target, nil)
}
