package ini

const (
	INI_GLOBAL_SECTION   = ""
	ASK_FOR_VALUE_PREFIX = "?"
)

type DefaultIniConfig struct {
	iniFilePath string

	iniFileValues  map[string]map[string]string
	currentSection string

	keysToStore map[string]bool
	storeAll    bool
}
