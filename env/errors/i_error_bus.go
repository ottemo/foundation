package errors

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ottemo/foundation/env"
	"runtime"
	"strconv"
	"strings"
)

// converts error message to OttemoError instance
func parseErrorMessage(message string) *OttemoError {
	resultError := new(OttemoError)

	reResult := MSG_REGEXP.FindStringSubmatch(message)

	if level, err := strconv.ParseInt(reResult[1], 10, 64); err == nil {
		resultError.Level = int(level)
	}
	resultError.Code = reResult[2]
	resultError.Message = reResult[3]

	if resultError.Code == "" {
		hasher := md5.New()
		hasher.Write([]byte(resultError.Message))

		resultError.Code = hex.EncodeToString(hasher.Sum(nil))
	}

	// primitive stack trace
	if INCLUDE_STACK {
		cutStopFlag := false
		skip := 0
		_, file, line, ok := runtime.Caller(skip)
		for ok {
			if cutStopFlag || !strings.Contains(file, "github.com/ottemo/foundation/env/") {
				cutStopFlag = true
				resultError.Stack += file + ":" + strconv.Itoa(line) + "\n"
			}

			skip += 1
			_, file, line, ok = runtime.Caller(skip)
		}
	}

	return resultError
}

// returns error level
func (it *DefaultErrorBus) GetErrorLevel(err error) int {
	if ottemoErr, ok := err.(*OttemoError); ok {
		return ottemoErr.Level
	}
	return parseErrorMessage(err.Error()).Level
}

// returns errors code
func (it *DefaultErrorBus) GetErrorCode(err error) string {
	if ottemoErr, ok := err.(*OttemoError); ok {
		return ottemoErr.Code
	}
	return parseErrorMessage(err.Error()).Code
}

// returns error message
func (it *DefaultErrorBus) GetErrorMessage(err error) string {
	if ottemoErr, ok := err.(*OttemoError); ok {
		return ottemoErr.Message
	}
	return err.Error()
}

// registers error listener
func (it *DefaultErrorBus) RegisterListener(listener env.F_ErrorListener) {
	it.listeners = append(it.listeners, listener)
}

// creates and processes OttemoError
func (it *DefaultErrorBus) New(message string) error {
	return it.Dispatch(parseErrorMessage(message))
}

// converts regular error to OttemoError and passes it through registered listeners
func (it *DefaultErrorBus) Dispatch(err error) error {
	if err == nil {
		return err
	}

	if ottemoErr, ok := err.(*OttemoError); ok {
		if ottemoErr.handled {
			return ottemoErr
		}

		ottemoErr.handled = true

		for _, listener := range it.listeners {
			if listener(ottemoErr) {
				break
			}
		}

		env.LogError(ottemoErr)

		return ottemoErr
	}

	return it.New(err.Error())
}
