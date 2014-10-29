package ini

const (
	ASK_FOR_VALUE_PREFIX = "?"
)

type DefaultIniConfig struct {
	iniFilePath string

	iniFileValues  map[string]map[string]string
	currentSection string

	keysToStore []string
	storeAll    bool
}
