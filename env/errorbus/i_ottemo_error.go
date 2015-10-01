package errorbus

import (
	"fmt"
)

// returns error message only
func (it *OttemoError) Error() string {
	if it.Level < hideLevel {
		return hideMessage
	}
	return it.Message
}

// ErrorMessage returns original error message
func (it *OttemoError) ErrorMessage() string {
	return it.Message
}

// ErrorFull returns error detail information about error
func (it *OttemoError) ErrorFull() string {
	message := it.Message
	if it.CallStack != "" {
		message += "\n" + it.CallStack
	}

	module := it.Module
	if module != "" {
		module += ":"
	}
	return fmt.Sprintf("%s%d:%s - %s", module, it.Level, it.Code, message)
}

// ErrorLevel returns error level - if specified or 0
func (it *OttemoError) ErrorLevel() int {
	return it.Level
}

// ErrorCode returns error code (hexadecimal value) if specified, otherwise MD5 over error message
func (it *OttemoError) ErrorCode() string {
	return it.Code
}

// ErrorCallStack returns error functions call stack for error
//   Note: ConstCollectStack constant should be set to true, otherwise, stack information will be blank
func (it *OttemoError) ErrorCallStack() string {
	return it.CallStack
}

// IsHandled returns handled flag
func (it *OttemoError) IsHandled() bool {
	return it.handled
}

// MarkHandled makes error as already processed (prevents from future processing)
func (it *OttemoError) MarkHandled() bool {
	it.handled = true
	return it.handled
}

// IsLogged returns logged flag
func (it *OttemoError) IsLogged() bool {
	return it.logged
}

// MarkLogged makes error as already logged (prevents from future logging)
func (it *OttemoError) MarkLogged() bool {
	it.logged = true
	return it.logged
}
