package env

import (
	"time"

	"github.com/ottemo/foundation/utils"
)

// Package global constants
const (
	ConstConfigTypeID       = utils.ConstDataTypeID
	ConstConfigTypeBoolean  = utils.ConstDataTypeBoolean
	ConstConfigTypeVarchar  = utils.ConstDataTypeVarchar
	ConstConfigTypeText     = utils.ConstDataTypeText
	ConstConfigTypeInteger  = utils.ConstDataTypeInteger
	ConstConfigTypeDecimal  = utils.ConstDataTypeDecimal
	ConstConfigTypeMoney    = utils.ConstDataTypeMoney
	ConstConfigTypeFloat    = utils.ConstDataTypeFloat
	ConstConfigTypeDatetime = utils.ConstDataTypeDatetime
	ConstConfigTypeJSON     = utils.ConstDataTypeJSON
	ConstConfigTypeHTML     = utils.ConstDataTypeHTML
	ConstConfigTypeGroup    = "group"
	ConstConfigTypeSecret   = "secret"

	ConstLogPrefixError   = "ERROR"
	ConstLogPrefixWarning = "WARNING"
	ConstLogPrefixDebug   = "DEBUG"
	ConstLogPrefixInfo    = "INFO"

	ConstErrorLevelAPI        = 10
	ConstErrorLevelModel      = 9
	ConstErrorLevelActor      = 8
	ConstErrorLevelHelper     = 7
	ConstErrorLevelService    = 4
	ConstErrorLevelServiceAct = 3
	ConstErrorLevelCritical   = 2
	ConstErrorLevelStartStop  = 1
	ConstErrorLevelExternal   = 0

	ConstErrorModule = "env"
	ConstErrorLevel  = ConstErrorLevelService
)
var(
	ConstCountries = map[string]string{
		"AD": "Andorra",
		"AE": "United Arab Emirates",
		"AF": "Afghanistan",
		"AG": "Antigua and Barbuda",
		"AI": "Anguilla",
		"AL": "Albania",
		"AM": "Armenia",
		"AN": "Netherlands Antilles",
		"AO": "Angola",
		"AR": "Argentina",
		"AT": "Austria",
		"AU": "Australia",
		"AW": "Aruba",
		"AX": "Aland Island (Finland)",
		"AZ": "Azerbaijan",
		"BA": "Bosnia-Herzegovina",
		"BB": "Barbados",
		"BD": "Bangladesh",
		"BE": "Belgium",
		"BF": "Burkina Faso",
		"BG": "Bulgaria",
		"BH": "Bahrain",
		"BI": "Burundi",
		"BJ": "Benin",
		"BM": "Bermuda",
		"BN": "Brunei Darussalam",
		"BO": "Bolivia",
		"BR": "Brazil",
		"BS": "Bahamas",
		"BT": "Bhutan",
		"BW": "Botswana",
		"BY": "Belarus",
		"BZ": "Belize",
		"CA": "Canada",
		"CC": "Cocos Island (Australia)",
		"CD": "Congo, Democratic Republic of the",
		"CF": "Central African Republic",
		"CG": "Congo, Republic of the",
		"CH": "Switzerland",
		"CI": "Ivory Coast (Cote d Ivoire)",
		"CK": "Cook Islands (New Zealand)",
		"CL": "Chile",
		"CM": "Cameroon",
		"CN": "China",
		"CO": "Colombia",
		"CR": "Costa Rica",
		"CU": "Cuba",
		"CV": "Cape Verde",
		"CX": "Christmas Island (Australia)",
		"CY": "Cyprus",
		"CZ": "Czech Republic",
		"DE": "Germany",
		"DJ": "Djibouti",
		"DK": "Denmark",
		"DM": "Dominica",
		"DO": "Dominican Republic",
		"DZ": "Algeria",
		"EC": "Ecuador",
		"EE": "Estonia",
		"EG": "Egypt",
		"ER": "Eritrea",
		"ES": "Spain",
		"ET": "Ethiopia",
		"FI": "Finland",
		"FJ": "Fiji",
		"FK": "Falkland Islands",
		"FM": "Micronesia, Federated States of",
		"FO": "Faroe Islands",
		"FR": "France",
		"GA": "Gabon",
		"GB": "Great Britain and Northern Ireland",
		"GD": "Grenada",
		"GE": "Georgia, Republic of",
		"GF": "French Guiana",
		"GH": "Ghana",
		"GI": "Gibraltar",
		"GL": "Greenland",
		"GM": "Gambia",
		"GN": "Guinea",
		"GP": "Guadeloupe",
		"GQ": "Equatorial Guinea",
		"GR": "Greece",
		"GS": "South Georgia (Falkland Islands)",
		"GT": "Guatemala",
		"GW": "Guinea-Bissau",
		"GY": "Guyana",
		"HK": "Hong Kong",
		"HN": "Honduras",
		"HR": "Croatia",
		"HT": "Haiti",
		"HU": "Hungary",
		"ID": "Indonesia",
		"IE": "Ireland",
		"IL": "Israel",
		"IN": "India",
		"IQ": "Iraq",
		"IR": "Iran",
		"IS": "Iceland",
		"IT": "Italy",
		"JM": "Jamaica",
		"JO": "Jordan",
		"JP": "Japan",
		"KE": "Kenya",
		"KG": "Kyrgyzstan",
		"KH": "Cambodia",
		"KI": "Kiribati",
		"KM": "Comoros",
		"KN": "Saint Kitts (Saint Christopher and Nevis)",
		"KP": "North Korea (Korea, Democratic People&#39;s Republic of)",
		"KR": "South Korea (Korea, Republic of)",
		"KW": "Kuwait",
		"KY": "Cayman Islands",
		"KZ": "Kazakhstan",
		"LA": "Laos",
		"LB": "Lebanon",
		"LC": "Saint Lucia",
		"LI": "Liechtenstein",
		"LK": "Sri Lanka",
		"LR": "Liberia",
		"LS": "Lesotho",
		"LT": "Lithuania",
		"LU": "Luxembourg",
		"LV": "Latvia",
		"LY": "Libya",
		"MA": "Morocco",
		"MC": "Monaco (France)",
		"MD": "Moldova",
		"MG": "Madagascar",
		"MK": "Macedonia, Republic of",
		"ML": "Mali",
		"MM": "Burma",
		"MN": "Mongolia",
		"MO": "Macao",
		"MQ": "Martinique",
		"MR": "Mauritania",
		"MS": "Montserrat",
		"MT": "Malta",
		"MU": "Mauritius",
		"MV": "Maldives",
		"MW": "Malawi",
		"MX": "Mexico",
		"MY": "Malaysia",
		"MZ": "Mozambique",
		"NA": "Namibia",
		"NC": "New Caledonia",
		"NE": "Niger",
		"NG": "Nigeria",
		"NI": "Nicaragua",
		"NL": "Netherlands",
		"NO": "Norway",
		"NP": "Nepal",
		"NR": "Nauru",
		"NZ": "New Zealand",
		"OM": "Oman",
		"PA": "Panama",
		"PE": "Peru",
		"PF": "French Polynesia",
		"PG": "Papua New Guinea",
		"PH": "Philippines",
		"PK": "Pakistan",
		"PL": "Poland",
		"PM": "Saint Pierre and Miquelon",
		"PN": "Pitcairn Island",
		"PT": "Portugal",
		"PY": "Paraguay",
		"QA": "Qatar",
		"RE": "Reunion",
		"RO": "Romania",
		"RS": "Serbia",
		"RU": "Russia",
		"RW": "Rwanda",
		"SA": "Saudi Arabia",
		"SB": "Solomon Islands",
		"SC": "Seychelles",
		"SD": "Sudan",
		"SE": "Sweden",
		"SG": "Singapore",
		"SH": "Saint Helena",
		"SI": "Slovenia",
		"SK": "Slovak Republic",
		"SL": "Sierra Leone",
		"SM": "San Marino",
		"SN": "Senegal",
		"SO": "Somalia",
		"SR": "Suriname",
		"ST": "Sao Tome and Principe",
		"SV": "El Salvador",
		"SY": "Syrian Arab Republic",
		"SZ": "Swaziland",
		"TC": "Turks and Caicos Islands",
		"TD": "Chad",
		"TG": "Togo",
		"TH": "Thailand",
		"TJ": "Tajikistan",
		"TK": "Tokelau (Union Group) (Western Samoa)",
		"TL": "East Timor (Timor-Leste, Democratic Republic of)",
		"TM": "Turkmenistan",
		"TN": "Tunisia",
		"TO": "Tonga",
		"TR": "Turkey",
		"TT": "Trinidad and Tobago",
		"TV": "Tuvalu",
		"TW": "Taiwan",
		"TZ": "Tanzania",
		"UA": "Ukraine",
		"UG": "Uganda",
		"UY": "Uruguay",
		"UZ": "Uzbekistan",
		"VA": "Vatican City",
		"VC": "Saint Vincent and the Grenadines",
		"VE": "Venezuela",
		"VG": "British Virgin Islands",
		"VN": "Vietnam",
		"VU": "Vanuatu",
		"WF": "Wallis and Futuna Islands",
		"WS": "Western Samoa",
		"YE": "Yemen",
		"YT": "Mayotte (France)",
		"ZA": "South Africa",
		"ZM": "Zambia",
		"ZW": "Zimbabwe",
		"US": "United States",
	}
)

// InterfaceSchedule is an interface to system schedule service
type InterfaceSchedule interface {
	Execute()
	RunTask(params map[string]interface{}) error

	Enable() error
	Disable() error

	Set(param string, value interface{}) error
	Get(param string) interface{}

	GetInfo() map[string]interface{}
}

// InterfaceScheduler is an interface to system scheduler service
type InterfaceScheduler interface {
	ListTasks() []string
	RegisterTask(name string, task FuncCronTask) error

	ScheduleAtTime(scheduleTime time.Time, taskName string, taskParams map[string]interface{}) (InterfaceSchedule, error)
	ScheduleRepeat(cronExpr string, taskName string, taskParams map[string]interface{}) (InterfaceSchedule, error)
	ScheduleOnce(cronExpr string, taskName string, taskParams map[string]interface{}) (InterfaceSchedule, error)

	ListSchedules() []InterfaceSchedule
}

// InterfaceEventBus is an interface to system event processor
type InterfaceEventBus interface {
	RegisterListener(event string, listener FuncEventListener)
	New(event string, eventData map[string]interface{})
}

// InterfaceErrorBus is an interface to system error processor
type InterfaceErrorBus interface {
	GetErrorLevel(error) int
	GetErrorCode(error) string
	GetErrorMessage(error) string

	RegisterListener(FuncErrorListener)

	Dispatch(err error) error
	Modify(err error, module string, level int, code string) error

	Prepare(module string, level int, code string, message string) error
	New(module string, level int, code string, message string) error
	Raw(message string) error
}

// InterfaceLogger is an interface to system logging service
type InterfaceLogger interface {
	Log(storage string, prefix string, message string)

	LogError(err error)
	LogEvent(f LogFields, eventName string)
}

type LogFields map[string]interface{}

// InterfaceIniConfig is an interface to startup configuration predefined values service
type InterfaceIniConfig interface {
	SetWorkingSection(sectionName string) error
	SetValue(valueName string, value string) error

	GetSectionValue(sectionName string, valueName string, defaultValue string) string
	GetValue(valueName string, defaultValue string) string

	ListSections() []string
	ListItems() []string
	ListSectionItems(sectionName string) []string
}

// InterfaceConfig is an interface to configuration values managing service
type InterfaceConfig interface {
	RegisterItem(Item StructConfigItem, Validator FuncConfigValueValidator) error
	UnregisterItem(Path string) error

	ListPathes() []string
	GetValue(Path string) interface{}
	SetValue(Path string, Value interface{}) error

	GetGroupItems() []StructConfigItem
	GetItemsInfo(Path string) []StructConfigItem

	Load() error
	Reload() error
}

// InterfaceOttemoError is an interface to errors generated by error bus service
type InterfaceOttemoError interface {
	ErrorFull() string
	ErrorLevel() int
	ErrorCode() string
	ErrorMessage() string
	ErrorCallStack() string

	IsHandled() bool
	MarkHandled() bool

	IsLogged() bool
	MarkLogged() bool

	error
}

// FuncConfigValueValidator is a configuration value validator callback function prototype
type FuncConfigValueValidator func(interface{}) (interface{}, error)

// FuncEventListener is an event listener callback function prototype
//   - return value is continue flag, so listener should return false to stop event propagation
type FuncEventListener func(string, map[string]interface{}) bool

// FuncErrorListener is an error listener callback function prototype
type FuncErrorListener func(error) bool

// FuncCronTask is a callback function prototype executes by scheduler
type FuncCronTask func(params map[string]interface{}) error

// StructConfigItem is a structure to hold information about particular configuration value
type StructConfigItem struct {
	Path  string
	Value interface{}

	Type string

	Editor  string
	Options interface{}

	Label       string
	Description string

	Image string
}
