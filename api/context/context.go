package context

import (
	"fmt"
	"runtime"
	"sync"
)

// package constants
const (
	debugOutput = true
)

// package variables
var (
	contexts      = make(map[string]map[string]interface{})
	contextsMutex sync.RWMutex

	proxies      = make(map[uintptr]uint)
	proxyBase    uint
	proxyStartPC uintptr
)

// init makes package initialization
func init() {

	// determination of available proxies
	var i uint
	discoverProxy := func() {
		pc, _, _, _ := runtime.Caller(1)

		if _, present := proxies[pc]; !present {
			proxies[pc] = i
		} else {
			pc, _, _, _ := runtime.Caller(2)
			proxies[pc] = i - 1
			i = 9999
		}
	}
	for i < 9999 {
		proxy(i, discoverProxy)
		i++
	}
	proxyBase = uint(len(proxies))

	// determination of proxy start point
	discoverStartPC := func() {
		// pointer to MakeContext func address,
		// just after it call proxies comming
		proxyStartPC, _, _, _ = runtime.Caller(2)
	}
	MakeContext(discoverStartPC)

	// debug output
	if debugOutput {
		tmp := fmt.Sprintln("init proxies:")
		for proxy, index := range proxies {
			file, line := runtime.FuncForPC(proxy).FileLine(proxy)
			tmp += fmt.Sprintf("%x=%x - %s:%d\n", index, proxy, file, line)
		}
		tmp += fmt.Sprintf("proxyStartPC = %x\n\n", proxyStartPC)
		fmt.Println(tmp)
	}
}

// getCallStack returns current call stack excluding runtime and self pc's (program counters)
func getCallStack(skip int) []uintptr {
	pcSize := 100
	pc := make([]uintptr, pcSize)

	n := runtime.Callers(2+skip, pc)
	for n >= pcSize {
		pcSize = pcSize * 10
		pc = make([]uintptr, pcSize)
		runtime.Callers(2+skip, pc)
	}

	pcSize = n - 1
	return pc[0 : n-1]
}

// MakeContext encodes new context index in current call-stack and then executes given target
//   - so, all the internal routine within target func, and sub-calls will have own context
//   - context is a map[string]interface{} where you can store any information
//   - context map will be eliminated automatically after target func finish
func MakeContext(target func()) {

	pc := getCallStack(1)

	// converting to string
	pcSource := ""
	for _, caller := range pc {
		pcSource += fmt.Sprintf("%x-", caller)
	}

	// looking for already in use contexts for stack
	var proxyIndex uint
	var contextKey string

	for true {
		contextKey = fmt.Sprintf("%s%d", pcSource, proxyIndex)

		contextsMutex.Lock()
		if _, present := contexts[contextKey]; present {
			contextsMutex.Unlock()
			proxyIndex++
			continue
		}

		contexts[contextKey] = make(map[string]interface{})
		contextsMutex.Unlock()

		break
	}

	if debugOutput {
		debugLogValue := fmt.Sprintln("MakeContext:")
		for _, x := range pc {
			file, line := runtime.FuncForPC(x).FileLine(x)
			debugLogValue += fmt.Sprintf("%x - %s:%d\n", x, file, line)
		}
		debugLogValue += fmt.Sprintln("making context: ", contextKey)
		fmt.Println(debugLogValue)
	}

	defer func() {
		contextsMutex.Lock()
		delete(contexts, contextKey)
		contextsMutex.Unlock()
	}()

	proxy(proxyIndex, target)
}

// GetContext - returns context assigned to current call-stack, or nil if no context
func GetContext() map[string]interface{} {
	pc := getCallStack(1)
	pcSize := len(pc)

	// looking for context beginning
	var baseIndex int
	for i := 0; i < pcSize; i++ {
		if pc[i] == proxyStartPC {
			baseIndex = i
			break
		}
	}

	// calculating context key
	var proxyPosition uint
	var proxyIndex uint

	var pcSource string

	for i := baseIndex - 1; i >= 0; i-- {
		if digitValue, present := proxies[pc[i]]; present {
			proxyIndex += digitValue
			proxyPosition++
		} else {
			break
		}
	}

	for i := baseIndex + 1; i < pcSize; i++ {
		pcSource += fmt.Sprintf("%x-", pc[i])
	}

	contextKey := fmt.Sprintf("%s%d", pcSource, proxyIndex)

	// if no context for this key, function will return nil
	contextsMutex.Lock()
	context := contexts[contextKey]
	contextsMutex.Unlock()

	if debugOutput {
		debugLogValue := fmt.Sprintln("context lookup: ", contextKey)
		for _, x := range pc {
			file, line := runtime.FuncForPC(x).FileLine(x)
			debugLogValue += fmt.Sprintf("%x - %s:%d\n", x, file, line)
		}
		debugLogValue += fmt.Sprint(context)
		fmt.Println(debugLogValue)
	}

	return context
}

// proxy points - the place which makes unique pc for target function call
//   - their amount specifies proxy base
//   - the more proxy base you have, then less stack growing (so 3 stack calls can provide 0xff^3=16777216 contexts)
func proxy(index uint, target func()) {

	switch index {
	case 0x00:
		target()
	case 0x01:
		target()
	case 0x02:
		target()
	case 0x03:
		target()
	case 0x04:
		target()
	case 0x05:
		target()
	case 0x06:
		target()
	case 0x07:
		target()
	case 0x08:
		target()
	case 0x09:
		target()
	case 0x0a:
		target()
	case 0x0b:
		target()
	case 0x0c:
		target()
	case 0x0d:
		target()
	case 0x0e:
		target()
	case 0x0f:
		target()
	case 0x10:
		target()
	case 0x11:
		target()
	case 0x12:
		target()
	case 0x13:
		target()
	case 0x14:
		target()
	case 0x15:
		target()
	case 0x16:
		target()
	case 0x17:
		target()
	case 0x18:
		target()
	case 0x19:
		target()
	case 0x1a:
		target()
	case 0x1b:
		target()
	case 0x1c:
		target()
	case 0x1d:
		target()
	case 0x1e:
		target()
	case 0x1f:
		target()
	case 0x20:
		target()
	case 0x21:
		target()
	case 0x22:
		target()
	case 0x23:
		target()
	case 0x24:
		target()
	case 0x25:
		target()
	case 0x26:
		target()
	case 0x27:
		target()
	case 0x28:
		target()
	case 0x29:
		target()
	case 0x2a:
		target()
	case 0x2b:
		target()
	case 0x2c:
		target()
	case 0x2d:
		target()
	case 0x2e:
		target()
	case 0x2f:
		target()
	case 0x30:
		target()
	case 0x31:
		target()
	case 0x32:
		target()
	case 0x33:
		target()
	case 0x34:
		target()
	case 0x35:
		target()
	case 0x36:
		target()
	case 0x37:
		target()
	case 0x38:
		target()
	case 0x39:
		target()
	case 0x3a:
		target()
	case 0x3b:
		target()
	case 0x3c:
		target()
	case 0x3d:
		target()
	case 0x3e:
		target()
	case 0x3f:
		target()
	case 0x40:
		target()
	case 0x41:
		target()
	case 0x42:
		target()
	case 0x43:
		target()
	case 0x44:
		target()
	case 0x45:
		target()
	case 0x46:
		target()
	case 0x47:
		target()
	case 0x48:
		target()
	case 0x49:
		target()
	case 0x4a:
		target()
	case 0x4b:
		target()
	case 0x4c:
		target()
	case 0x4d:
		target()
	case 0x4e:
		target()
	case 0x4f:
		target()
	case 0x50:
		target()
	case 0x51:
		target()
	case 0x52:
		target()
	case 0x53:
		target()
	case 0x54:
		target()
	case 0x55:
		target()
	case 0x56:
		target()
	case 0x57:
		target()
	case 0x58:
		target()
	case 0x59:
		target()
	case 0x5a:
		target()
	case 0x5b:
		target()
	case 0x5c:
		target()
	case 0x5d:
		target()
	case 0x5e:
		target()
	case 0x5f:
		target()
	case 0x60:
		target()
	case 0x61:
		target()
	case 0x62:
		target()
	case 0x63:
		target()
	case 0x64:
		target()
	case 0x65:
		target()
	case 0x66:
		target()
	case 0x67:
		target()
	case 0x68:
		target()
	case 0x69:
		target()
	case 0x6a:
		target()
	case 0x6b:
		target()
	case 0x6c:
		target()
	case 0x6d:
		target()
	case 0x6e:
		target()
	case 0x6f:
		target()
	case 0x70:
		target()
	case 0x71:
		target()
	case 0x72:
		target()
	case 0x73:
		target()
	case 0x74:
		target()
	case 0x75:
		target()
	case 0x76:
		target()
	case 0x77:
		target()
	case 0x78:
		target()
	case 0x79:
		target()
	case 0x7a:
		target()
	case 0x7b:
		target()
	case 0x7c:
		target()
	case 0x7d:
		target()
	case 0x7e:
		target()
	case 0x7f:
		target()
	case 0x80:
		target()
	case 0x81:
		target()
	case 0x82:
		target()
	case 0x83:
		target()
	case 0x84:
		target()
	case 0x85:
		target()
	case 0x86:
		target()
	case 0x87:
		target()
	case 0x88:
		target()
	case 0x89:
		target()
	case 0x8a:
		target()
	case 0x8b:
		target()
	case 0x8c:
		target()
	case 0x8d:
		target()
	case 0x8e:
		target()
	case 0x8f:
		target()
	case 0x90:
		target()
	case 0x91:
		target()
	case 0x92:
		target()
	case 0x93:
		target()
	case 0x94:
		target()
	case 0x95:
		target()
	case 0x96:
		target()
	case 0x97:
		target()
	case 0x98:
		target()
	case 0x99:
		target()
	case 0x9a:
		target()
	case 0x9b:
		target()
	case 0x9c:
		target()
	case 0x9d:
		target()
	case 0x9e:
		target()
	case 0x9f:
		target()
	case 0xa0:
		target()
	case 0xa1:
		target()
	case 0xa2:
		target()
	case 0xa3:
		target()
	case 0xa4:
		target()
	case 0xa5:
		target()
	case 0xa6:
		target()
	case 0xa7:
		target()
	case 0xa8:
		target()
	case 0xa9:
		target()
	case 0xaa:
		target()
	case 0xab:
		target()
	case 0xac:
		target()
	case 0xad:
		target()
	case 0xae:
		target()
	case 0xaf:
		target()
	case 0xb0:
		target()
	case 0xb1:
		target()
	case 0xb2:
		target()
	case 0xb3:
		target()
	case 0xb4:
		target()
	case 0xb5:
		target()
	case 0xb6:
		target()
	case 0xb7:
		target()
	case 0xb8:
		target()
	case 0xb9:
		target()
	case 0xba:
		target()
	case 0xbb:
		target()
	case 0xbc:
		target()
	case 0xbd:
		target()
	case 0xbe:
		target()
	case 0xbf:
		target()
	case 0xc0:
		target()
	case 0xc1:
		target()
	case 0xc2:
		target()
	case 0xc3:
		target()
	case 0xc4:
		target()
	case 0xc5:
		target()
	case 0xc6:
		target()
	case 0xc7:
		target()
	case 0xc8:
		target()
	case 0xc9:
		target()
	case 0xca:
		target()
	case 0xcb:
		target()
	case 0xcc:
		target()
	case 0xcd:
		target()
	case 0xce:
		target()
	case 0xcf:
		target()
	case 0xd0:
		target()
	case 0xd1:
		target()
	case 0xd2:
		target()
	case 0xd3:
		target()
	case 0xd4:
		target()
	case 0xd5:
		target()
	case 0xd6:
		target()
	case 0xd7:
		target()
	case 0xd8:
		target()
	case 0xd9:
		target()
	case 0xda:
		target()
	case 0xdb:
		target()
	case 0xdc:
		target()
	case 0xdd:
		target()
	case 0xde:
		target()
	case 0xdf:
		target()
	case 0xe0:
		target()
	case 0xe1:
		target()
	case 0xe2:
		target()
	case 0xe3:
		target()
	case 0xe4:
		target()
	case 0xe5:
		target()
	case 0xe6:
		target()
	case 0xe7:
		target()
	case 0xe8:
		target()
	case 0xe9:
		target()
	case 0xea:
		target()
	case 0xeb:
		target()
	case 0xec:
		target()
	case 0xed:
		target()
	case 0xee:
		target()
	case 0xef:
		target()
	case 0xf0:
		target()
	case 0xf1:
		target()
	case 0xf2:
		target()
	case 0xf3:
		target()
	case 0xf4:
		target()
	case 0xf5:
		target()
	case 0xf6:
		target()
	case 0xf7:
		target()
	case 0xf8:
		target()
	case 0xf9:
		target()
	case 0xfa:
		target()
	case 0xfb:
		target()
	case 0xfc:
		target()
	case 0xfd:
		target()
	case 0xfe:
		target()
	case 0xff:
		target()
	}

	if index > 0xff {
		proxy(index-0xff, target)
	}
}
