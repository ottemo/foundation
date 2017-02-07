package impex

import (
	"fmt"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"bytes"
	"strings"
)

// init makes package self-initialization routine
func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	if err := RegisterImportCommand("IMPORT", new(ImportCmdImport)); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := RegisterImportCommand("INSERT", new(ImportCmdInsert)); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := RegisterImportCommand("UPDATE", new(ImportCmdUpdate)); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := RegisterImportCommand("DELETE", new(ImportCmdDelete)); err != nil {
		_ = env.ErrorDispatch(err)
	}

	if err := RegisterImportCommand("STORE", new(ImportCmdStore)); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := RegisterImportCommand("ALIAS", new(ImportCmdAlias)); err != nil {
		_ = env.ErrorDispatch(err)
	}

	if err := RegisterImportCommand("MEDIA", new(ImportCmdMedia)); err != nil {
		_ = env.ErrorDispatch(err)
	}

	if err := RegisterImportCommand("ATTRIBUTE_ADD", new(ImportCmdAttributeAdd)); err != nil {
		_ = env.ErrorDispatch(err)
	}

	// initializing column conversion functions
	ConversionFuncs["log"] = func(args ...interface{}) string {
		env.Log(ConstLogFileName, env.ConstLogPrefixDebug, fmt.Sprint(args))
		return ""
	}

	ConversionFuncs["logf"] = func(format string, args ...interface{}) string {
		env.Log(ConstLogFileName, env.ConstLogPrefixDebug, fmt.Sprintf(format, args))
		return ""
	}

	ConversionFuncs["set"] = func(context map[string]interface{}, key string, value interface{}) string {
		keyLevels := strings.Split(key, ".")
		key = keyLevels[len(keyLevels)-1]
		for _, key := range keyLevels[:len(keyLevels)-1] {
			if value, present := context[key]; present {
				if dictValue, ok := value.(map[string]interface{}); ok {
					context = dictValue
					continue
				}
			}
			newValue := make(map[string]interface{})
			context[key] = newValue
			context = newValue
		}
		context[key] = value

		return ""
	}

	ConversionFuncs["alias"] = func(context map[string]interface{}, args ...interface{}) string {
		var value string
		var aliases map[string]interface{}

		// looking for aliaves and value
		if len(args) > 0 {
			value = utils.InterfaceToString(args[0])
		}

		// checking we have value in given context and using it if not specified
		if contextValue, present := context["value"]; present && len(args) == 0 {
			value = utils.InterfaceToString(contextValue)
		}

		if len(args) > 1 {
			if aliasesDict, ok := args[1].(map[string]interface{}); ok {
				aliases = aliasesDict
			}
		}

		// checking we have aliases dictionary in given context and using it if not specified
		if aliasesValue, present := context["alias"]; present && aliases == nil {
			if aliasesDict, ok := aliasesValue.(map[string]interface{}); ok {
				aliases = aliasesDict
			}
		}

		// checking it for a sense to continue
		if aliases != nil && value != "" {

			// looking for full-text alias
			if aliasValue, present := aliases[value]; present {
				return utils.InterfaceToString(aliasValue)
			}

			// searching in-text aliases prefixed with "@"
			if strings.Contains(value, "@") {
				var result bytes.Buffer
				startIdx := -1
				for idx, chr := range value {
					switch {
					case chr == '@':
						if startIdx > -1 {
							foundAlias := value[startIdx:idx]
							if aliasValue, present := aliases[foundAlias[1:]]; present {
								foundAlias = utils.InterfaceToString(aliasValue)
							}
							result.WriteString(foundAlias)
						}
						startIdx = idx

					case !(chr >= 'A' && chr <= 'Z') && !(chr >= 'a' && chr <= 'z') && chr != '_':
						if startIdx > -1 {
							foundAlias := value[startIdx:idx]
							if aliasValue, present := aliases[foundAlias[1:]]; present {
								foundAlias = utils.InterfaceToString(aliasValue)
							}
							result.WriteString(foundAlias)

							startIdx = -1
						}
						result.WriteRune(chr)

					default:
						if startIdx == -1 {
							result.WriteRune(chr)
						}
					}
				}

				if startIdx > -1 {
					foundAlias := value[startIdx:]
					if aliasValue, present := aliases[foundAlias[1:]]; present {
						foundAlias = utils.InterfaceToString(aliasValue)
					}
					result.WriteString(foundAlias)
				}

				return result.String()
			}
		}

		return ""
	}

}
