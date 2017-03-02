package errorbus

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ottemo/foundation/env"
	"runtime"
	"strconv"
	"strings"
	"fmt"
)

// backtrace collects current call stack information and modifies given error with
func (it *DefaultErrorBus) backtrace(ottemoErr *OttemoError) {
	if ConstCollectCallStack {
		cutStopFlag := false
		step := 0

		_, file, line, ok := runtime.Caller(step)
		for ok {
			if !cutStopFlag && !strings.Contains(file, "env/helpers.go") && !strings.Contains(file, "env/errorbus/") {
				cutStopFlag = true
			}

			if cutStopFlag {
				ottemoErr.CallStack += file + ":" + strconv.Itoa(line) + "\n"
			}

			step++
			_, file, line, ok = runtime.Caller(step)
		}
	}
}

// process is set of routines to handle error within system
func (it *DefaultErrorBus) process(ottemoErr *OttemoError) *OttemoError {

	if ottemoErr.handled {
		return ottemoErr
	}

	ottemoErr.handled = true

	it.backtrace(ottemoErr)

	for _, listener := range it.listeners {
		if listener(ottemoErr) {
			break
		}
	}

	if debug {
		env.LogError(ottemoErr)
	}

	return ottemoErr
}

// converts error message to OttemoError instance
func (it *DefaultErrorBus) parseErrorMessage(message string) *OttemoError {
	resultError := new(OttemoError)

	reResult := ConstMsgRegexp.FindStringSubmatch(message)

	if level, err := strconv.ParseInt(reResult[2], 10, 64); err == nil {
		resultError.Level = int(level)
	}
	resultError.Module = reResult[1]
	resultError.Code = reResult[3]
	resultError.Message = reResult[4]

	if resultError.Code == "" {
		hasher := md5.New()
		if _, err := hasher.Write([]byte(resultError.Message)); err != nil {
			fmt.Println("37a7ed91-1673-46f6-8160-685ee1cf2828", err.Error())
		}

		resultError.Code = hex.EncodeToString(hasher.Sum(nil))
	}

	rawCode := strings.Replace(resultError.Code, "-", "", -1)
	if len(rawCode) == 32 {
		resultError.Code = rawCode[0:8] + "-" + rawCode[8:12] + "-" + rawCode[12:16] + "-" + rawCode[16:20] + "-" + rawCode[20:]
	}

	return resultError
}

// GetErrorLevel returns error level
func (it *DefaultErrorBus) GetErrorLevel(err error) int {
	if ottemoErr, ok := err.(*OttemoError); ok {
		return ottemoErr.Level
	}
	return it.parseErrorMessage(err.Error()).Level
}

// GetErrorCode returns errors code
func (it *DefaultErrorBus) GetErrorCode(err error) string {
	if ottemoErr, ok := err.(*OttemoError); ok {
		return ottemoErr.Code
	}
	return it.parseErrorMessage(err.Error()).Code
}

// GetErrorMessage returns error message
func (it *DefaultErrorBus) GetErrorMessage(err error) string {
	if ottemoErr, ok := err.(*OttemoError); ok {
		return ottemoErr.Message
	}
	return err.Error()
}

// RegisterListener registers error listener
func (it *DefaultErrorBus) RegisterListener(listener env.FuncErrorListener) {
	it.listeners = append(it.listeners, listener)
}

// Prepare creates OttemoError without processing
func (it *DefaultErrorBus) Prepare(module string, level int, code string, message string) error {
	ottemoErr := new(OttemoError)

	ottemoErr.Module = module
	ottemoErr.Level = level
	ottemoErr.Code = code
	ottemoErr.Message = message

	return ottemoErr
}

// New creates and processes OttemoError
func (it *DefaultErrorBus) New(module string, level int, code string, message string) error {
	ottemoErr := new(OttemoError)

	ottemoErr.Module = module
	ottemoErr.Level = level
	ottemoErr.Code = code
	ottemoErr.Message = message

	return it.process(ottemoErr)
}

// Raw creates and processes OttemoError encoded in given string
func (it *DefaultErrorBus) Raw(message string) error {
	ottemoErr := it.parseErrorMessage(message)
	ottemoErr.Level = 10

	return it.process(ottemoErr)
}

// Modify works similar to Dispatch but allows to specify some additional information
func (it *DefaultErrorBus) Modify(err error, module string, level int, code string) error {
	var ottemoErr *OttemoError

	if typedError, ok := err.(*OttemoError); !ok {
		ottemoErr = typedError
	} else {
		ottemoErr = it.parseErrorMessage(err.Error())
	}

	ottemoErr.Module = module
	ottemoErr.Level = level
	ottemoErr.Code = code

	return it.process(ottemoErr)
}

// Dispatch converts regular error to OttemoError and passes it through registered listeners
func (it *DefaultErrorBus) Dispatch(err error) error {
	if err == nil {
		return err
	}

	if ottemoErr, ok := err.(*OttemoError); ok {
		return it.process(ottemoErr)
	}

	ottemoErr := it.parseErrorMessage(err.Error())
	ottemoErr.Level = 0

	return it.process(ottemoErr)
}
