package env

import (
	"errors"
	"github.com/ottemo/foundation/utils"
)

// ConfigGetValue returns config value or nil if not present
func ConfigGetValue(Path string) interface{} {
	if config := GetConfig(); config != nil {
		return config.GetValue(Path)
	}

	return nil
}

// IniValue returns value from ini file or "" if not present
func IniValue(Path string) string {
	if iniConfig := GetIniConfig(); iniConfig != nil {
		return iniConfig.GetValue(Path, "")
	}
	return ""
}

// Log logs general purpose message
func Log(storage string, prefix string, message string) {
	if logger := GetLogger(); logger != nil {
		logger.Log(storage, prefix, message)
	}
}

// LogMap logs map data
func LogMap(storage string, data map[string]interface{}) {
	if logger := GetLogger(); logger != nil {
		logger.LogMap(storage, data)
	}
}

// LogError logs an error message
func LogError(err error) {
	err = ErrorDispatch(err)
	if logger := GetLogger(); logger != nil {
		logger.LogError(err)
	}
}

// LogMessage is a Log function short form for info messages in default storage
func LogMessage(message string) {
	if logger := GetLogger(); logger != nil {
		logger.LogMessage(message)
	}
}

// ErrorLevel returns error level of given error
func ErrorLevel(err error) int {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorLevel(err)
	}
	return 0
}

// ErrorCode returns error code of given error
func ErrorCode(err error) string {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorCode(err)
	}
	return ""
}

// ErrorMessage returns message of given error
func ErrorMessage(err error) string {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.GetErrorMessage(err)
	}
	return err.Error()
}

// ErrorRegisterListener registers listener for error bus
func ErrorRegisterListener(listener FuncErrorListener) {
	if errorBus := GetErrorBus(); errorBus != nil {
		errorBus.RegisterListener(listener)
	}
}

// ErrorDispatch handles error, returns new one you should pass next
func ErrorDispatch(err error) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.Dispatch(err)
	}
	return err
}

// ErrorModify works similar to ErrorDispatch but allows to set error level, code and module name
func ErrorModify(err error, module string, level int, code string) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.Modify(err, module, level, code)
	}
	return err
}

// Error creates new error (without dispatch)
func Error(module string, level int, code string, message string) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.Prepare(module, level, code, message)
	}
	return errors.New(message)
}

// ErrorNew creates new error and dispatches it
func ErrorNew(module string, level int, code string, message string) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.New(module, level, code, message)
	}
	return errors.New(message)
}

// ErrorRaw creates new error by parsing given string (seek for module name, level and code inside) and dispatches it
func ErrorRaw(message string) error {
	if errorBus := GetErrorBus(); errorBus != nil {
		return errorBus.Raw(message)
	}
	return errors.New(message)
}

// EventRegisterListener registers listener for event bus
func EventRegisterListener(event string, listener FuncEventListener) {
	if eventBus := GetEventBus(); eventBus != nil {
		eventBus.RegisterListener(event, listener)
	}
}

// Event emits new event for registered listeners
func Event(event string, args map[string]interface{}) {
	if eventBus := GetEventBus(); eventBus != nil {
		eventBus.New(event, args)
	}
}

// TypeParse shortcut for utils.DataTypeParse
func TypeParse(typeName string) utils.DataType {
	return utils.DataTypeParse(typeName)
}

// TypeWPrecisionAndScale shortcut for utils.DataTypeWPrecisionAndScale
func TypeWPrecisionAndScale(dataType string, precision int, scale int) string {
	return utils.DataTypeWPrecisionAndScale(dataType, precision, scale)
}

// TypeWPrecision shortcut for utils.DataTypeWPrecision
func TypeWPrecision(dataType string, precision int) string {
	return utils.DataTypeWPrecision(dataType, precision)
}

// TypeArrayOf shortcut for utils.DataTypeArrayOf
func TypeArrayOf(dataType string) string {
	return utils.DataTypeArrayOf(dataType)
}

// TypeIsArray shortcut for utils.DataTypeIsArray
func TypeIsArray(dataType string) bool {
	return utils.DataTypeIsArray(dataType)
}

// TypeIsString shortcut for utils.DataTypeIsString
func TypeIsString(dataType string) bool {
	return utils.DataTypeIsString(dataType)
}

// TypeIsFloat shortcut for utils.DataTypeIsFloat
func TypeIsFloat(dataType string) bool {
	return utils.DataTypeIsFloat(dataType)
}
