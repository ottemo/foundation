package errors

import (
	"fmt"
)

// returns error message only
func (it *OttemoError) Error() string {
	return it.Message
}

// ErrorFull returns error detail information about error
func (it *OttemoError) ErrorFull() string {
	message := it.Message
	if it.Stack != "" {
		message += "\n" + it.Stack
	}
	return fmt.Sprintf("%d:%s - %s", it.Level, it.Code, message)
}

// ErrorLevel returns error level - if specified or 0
func (it *OttemoError) ErrorLevel() int {
	return it.Level
}

// ErrorCode returns error code (hexadecimal value) if specified, otherwise MD5 over error message
func (it *OttemoError) ErrorCode() string {
	return it.Code
}

// ErrorStack returns error functions call stack for error
//   Note: ConstCollectStack constant should be set to true, otherwise, stack information will be blank
func (it *OttemoError) ErrorStack() string {
	return it.Stack
}
