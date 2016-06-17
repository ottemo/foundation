package utils

/*
refer to http://www.timeanddate.com/time/zones/
*/

import (
	"fmt"
	"strings"
	"time"
)

// TimeZones list of known Time Zone Abbreviations – Worldwide List
var TimeZones = map[string]time.Duration{
	"A":     1*time.Hour + 00*time.Minute,
	"ACDT":  10*time.Hour + 30*time.Minute,
	"ACST":  9*time.Hour + 30*time.Minute,
	"ACT":   -5*time.Hour + 00*time.Minute,
	"ACWST": 8*time.Hour + 45*time.Minute,
	"ADT":   3*time.Hour + 00*time.Minute,
	"AEDT":  11*time.Hour + 00*time.Minute,
	"AEST":  10*time.Hour + 00*time.Minute,
	"AET":   11*time.Hour + 00*time.Minute,
	"AFT":   4*time.Hour + 30*time.Minute,
	"AKDT":  -8*time.Hour + 00*time.Minute,
	"AKST":  -9*time.Hour + 00*time.Minute,
	"ALMT":  6*time.Hour + 00*time.Minute,
	"AMST":  -3*time.Hour + 00*time.Minute,
	"AMT":   -4*time.Hour + 00*time.Minute,
	"ANAST": 12*time.Hour + 00*time.Minute,
	"ANAT":  12*time.Hour + 00*time.Minute,
	"AoE":   -12*time.Hour + 00*time.Minute,
	"AQTT":  5*time.Hour + 00*time.Minute,
	"ART":   -3*time.Hour + 00*time.Minute,
	"AST":   3*time.Hour + 00*time.Minute,
	"AT":    -4*time.Hour + 00*time.Minute,
	"AWDT":  9*time.Hour + 00*time.Minute,
	"AWST":  8*time.Hour + 00*time.Minute,
	"AZOST": 0*time.Hour + 00*time.Minute,
	"AZOT":  -1*time.Hour + 00*time.Minute,
	"AZST":  5*time.Hour + 00*time.Minute,
	"AZT":   4*time.Hour + 00*time.Minute,
	"B":     2*time.Hour + 00*time.Minute,
	"BNT":   8*time.Hour + 00*time.Minute,
	"BOT":   -4*time.Hour + 00*time.Minute,
	"BRST":  -2*time.Hour + 00*time.Minute,
	"BRT":   -3*time.Hour + 00*time.Minute,
	"BST":   6*time.Hour + 00*time.Minute,
	"BTT":   6*time.Hour + 00*time.Minute,
	"C":     3*time.Hour + 00*time.Minute,
	"CAST":  8*time.Hour + 00*time.Minute,
	"CAT":   2*time.Hour + 00*time.Minute,
	"CCT":   6*time.Hour + 30*time.Minute,
	"CDT":   -5*time.Hour + 00*time.Minute,
	"CEST":  2*time.Hour + 00*time.Minute,
	"CET":   1*time.Hour + 00*time.Minute,
	"CHAST": 12*time.Hour + 45*time.Minute,
	"CHOT":  8*time.Hour + 00*time.Minute,
	"ChST":  10*time.Hour + 00*time.Minute,
	"CHUT":  10*time.Hour + 00*time.Minute,
	"CKT":   -10*time.Hour + 00*time.Minute,
	"CLST":  -3*time.Hour + 00*time.Minute,
	"CLT":   -4*time.Hour + 00*time.Minute,
	"COT":   -5*time.Hour + 00*time.Minute,
	"CST":   -6*time.Hour + 00*time.Minute,
	"CT":    -6*time.Hour + 00*time.Minute,
	"CVT":   -1*time.Hour + 00*time.Minute,
	"CXT":   7*time.Hour + 00*time.Minute,
	"D":     4*time.Hour + 00*time.Minute,
	"DAVT":  7*time.Hour + 00*time.Minute,
	"DDUT":  10*time.Hour + 00*time.Minute,
	"E":     5*time.Hour + 00*time.Minute,
	"EASST": -5*time.Hour + 00*time.Minute,
	"EAST":  -6*time.Hour + 00*time.Minute,
	"EAT":   3*time.Hour + 00*time.Minute,
	"ECT":   -5*time.Hour + 00*time.Minute,
	"EDT":   -4*time.Hour + 00*time.Minute,
	"EEST":  3*time.Hour + 00*time.Minute,
	"EET":   2*time.Hour + 00*time.Minute,
	"EGST":  0*time.Hour + 00*time.Minute,
	"EGT":   -1*time.Hour + 00*time.Minute,
	"EST":   -5*time.Hour + 00*time.Minute,
	"ET":    -5*time.Hour + 00*time.Minute,
	"F":     6*time.Hour + 00*time.Minute,
	"FET":   3*time.Hour + 00*time.Minute,
	"FJT":   12*time.Hour + 00*time.Minute,
	"FKST":  -3*time.Hour + 00*time.Minute,
	"FKT":   -4*time.Hour + 00*time.Minute,
	"FNT":   -2*time.Hour + 00*time.Minute,
	"G":     7*time.Hour + 00*time.Minute,
	"GALT":  -6*time.Hour + 00*time.Minute,
	"GAMT":  -9*time.Hour + 00*time.Minute,
	"GET":   4*time.Hour + 00*time.Minute,
	"GFT":   -3*time.Hour + 00*time.Minute,
	"GILT":  12*time.Hour + 00*time.Minute,
	"GMT":   0*time.Hour + 00*time.Minute,
	"GST":   4*time.Hour + 00*time.Minute,
	"GYT":   -4*time.Hour + 00*time.Minute,
	"H":     8*time.Hour + 00*time.Minute,
	"HADT":  -9*time.Hour + 00*time.Minute,
	"HAST":  -10*time.Hour + 00*time.Minute,
	"HKT":   8*time.Hour + 00*time.Minute,
	"HOVT":  7*time.Hour + 00*time.Minute,
	"I":     9*time.Hour + 00*time.Minute,
	"ICT":   7*time.Hour + 00*time.Minute,
	"IDT":   3*time.Hour + 00*time.Minute,
	"IOT":   6*time.Hour + 00*time.Minute,
	"IRDT":  4*time.Hour + 30*time.Minute,
	"IRKST": 9*time.Hour + 00*time.Minute,
	"IRKT":  8*time.Hour + 00*time.Minute,
	"IRST":  3*time.Hour + 30*time.Minute,
	"IST":   5*time.Hour + 30*time.Minute,
	"JST":   9*time.Hour + 00*time.Minute,
	"K":     10*time.Hour + 00*time.Minute,
	"KGT":   6*time.Hour + 00*time.Minute,
	"KOST":  11*time.Hour + 00*time.Minute,
	"KRAST": 8*time.Hour + 00*time.Minute,
	"KRAT":  7*time.Hour + 00*time.Minute,
	"KST":   9*time.Hour + 00*time.Minute,
	"KUYT":  4*time.Hour + 00*time.Minute,
	"L":     11*time.Hour + 00*time.Minute,
	"LHDT":  11*time.Hour + 00*time.Minute,
	"LHST":  10*time.Hour + 30*time.Minute,
	"LINT":  14*time.Hour + 00*time.Minute,
	"M":     12*time.Hour + 00*time.Minute,
	"MAGST": 12*time.Hour + 00*time.Minute,
	"MAGT":  10*time.Hour + 00*time.Minute,
	"MART":  -9*time.Hour + 30*time.Minute,
	"MAWT":  5*time.Hour + 00*time.Minute,
	"MDT":   -6*time.Hour + 00*time.Minute,
	"MHT":   12*time.Hour + 00*time.Minute,
	"MMT":   6*time.Hour + 30*time.Minute,
	"MSD":   4*time.Hour + 00*time.Minute,
	"MSK":   3*time.Hour + 00*time.Minute,
	"MST":   -7*time.Hour + 00*time.Minute,
	"MT":    -7*time.Hour + 00*time.Minute,
	"MUT":   4*time.Hour + 00*time.Minute,
	"MVT":   5*time.Hour + 00*time.Minute,
	"MYT":   8*time.Hour + 00*time.Minute,
	"N":     -1*time.Hour + 00*time.Minute,
	"NCT":   11*time.Hour + 00*time.Minute,
	"NDT":   -2*time.Hour + 30*time.Minute,
	"NFT":   11*time.Hour + 30*time.Minute,
	"NOVST": 7*time.Hour + 00*time.Minute,
	"NOVT":  6*time.Hour + 00*time.Minute,
	"NPT":   5*time.Hour + 45*time.Minute,
	"NRT":   12*time.Hour + 00*time.Minute,
	"NST":   -3*time.Hour + 30*time.Minute,
	"NUT":   -11*time.Hour + 00*time.Minute,
	"NZST":  12*time.Hour + 00*time.Minute,
	"O":     -2*time.Hour + 00*time.Minute,
	"OMSST": 7*time.Hour + 00*time.Minute,
	"OMST":  6*time.Hour + 00*time.Minute,
	"ORAT":  5*time.Hour + 00*time.Minute,
	"P":     -3*time.Hour + 00*time.Minute,
	"PDT":   -7*time.Hour + 00*time.Minute,
	"PET":   -5*time.Hour + 00*time.Minute,
	"PETST": 12*time.Hour + 00*time.Minute,
	"PETT":  12*time.Hour + 00*time.Minute,
	"PGT":   10*time.Hour + 00*time.Minute,
	"PHT":   8*time.Hour + 00*time.Minute,
	"PKT":   5*time.Hour + 00*time.Minute,
	"PMDT":  -2*time.Hour + 00*time.Minute,
	"PMST":  -3*time.Hour + 00*time.Minute,
	"PONT":  11*time.Hour + 00*time.Minute,
	"PST":   -8*time.Hour + 00*time.Minute,
	"PT":    -8*time.Hour + 00*time.Minute,
	"PWT":   9*time.Hour + 00*time.Minute,
	"PYST":  -3*time.Hour + 00*time.Minute,
	"PYT":   -4*time.Hour + 00*time.Minute,
	"Q":     -4*time.Hour + 00*time.Minute,
	"QYZT":  6*time.Hour + 00*time.Minute,
	"R":     -5*time.Hour + 00*time.Minute,
	"RET":   4*time.Hour + 00*time.Minute,
	"ROTT":  -3*time.Hour + 00*time.Minute,
	"S":     -6*time.Hour + 00*time.Minute,
	"SAKT":  10*time.Hour + 00*time.Minute,
	"SAMT":  4*time.Hour + 00*time.Minute,
	"SAST":  2*time.Hour + 00*time.Minute,
	"SBT":   11*time.Hour + 00*time.Minute,
	"SCT":   4*time.Hour + 00*time.Minute,
	"SGT":   8*time.Hour + 00*time.Minute,
	"SRET":  11*time.Hour + 00*time.Minute,
	"SRT":   -3*time.Hour + 00*time.Minute,
	"SST":   -11*time.Hour + 00*time.Minute,
	"SYOT":  3*time.Hour + 00*time.Minute,
	"T":     -7*time.Hour + 00*time.Minute,
	"TAHT":  -10*time.Hour + 00*time.Minute,
	"TFT":   5*time.Hour + 00*time.Minute,
	"TJT":   5*time.Hour + 00*time.Minute,
	"TLT":   9*time.Hour + 00*time.Minute,
	"TMT":   5*time.Hour + 00*time.Minute,
	"TVT":   12*time.Hour + 00*time.Minute,
	"U":     -8*time.Hour + 00*time.Minute,
	"ULAT":  8*time.Hour + 00*time.Minute,
	"UTC":   0*time.Hour + 00*time.Minute,
	"UYST":  -2*time.Hour + 00*time.Minute,
	"UYT":   -3*time.Hour + 00*time.Minute,
	"UZT":   5*time.Hour + 00*time.Minute,
	"V":     -9*time.Hour + 00*time.Minute,
	"VET":   -4*time.Hour + 30*time.Minute,
	"VLAST": 11*time.Hour + 00*time.Minute,
	"VLAT":  10*time.Hour + 00*time.Minute,
	"VOST":  6*time.Hour + 00*time.Minute,
	"VUT":   11*time.Hour + 00*time.Minute,
	"W":     -10*time.Hour + 00*time.Minute,
	"WAKT":  12*time.Hour + 00*time.Minute,
	"WARST": -3*time.Hour + 00*time.Minute,
	"WAST":  2*time.Hour + 00*time.Minute,
	"WAT":   1*time.Hour + 00*time.Minute,
	"WEST":  1*time.Hour + 00*time.Minute,
	"WET":   0*time.Hour + 00*time.Minute,
	"WFT":   12*time.Hour + 00*time.Minute,
	"WGST":  -2*time.Hour + 00*time.Minute,
	"WGT":   -3*time.Hour + 00*time.Minute,
	"WIB":   7*time.Hour + 00*time.Minute,
	"WIT":   9*time.Hour + 00*time.Minute,
	"WITA":  8*time.Hour + 00*time.Minute,
	"WST":   1*time.Hour + 00*time.Minute,
	"WT":    0*time.Hour + 00*time.Minute,
	"X":     -11*time.Hour + 00*time.Minute,
	"Y":     -12*time.Hour + 00*time.Minute,
	"YAKST": 10*time.Hour + 00*time.Minute,
	"YAKT":  9*time.Hour + 00*time.Minute,
	"YAPT":  10*time.Hour + 00*time.Minute,
	"YEKST": 6*time.Hour + 00*time.Minute,
	"YEKT":  5*time.Hour + 00*time.Minute,
	"Z":     0*time.Hour + 00*time.Minute,
}

// ParseTimeZone returns zone name and it's offset from GMT
func ParseTimeZone(timeZone string) (string, time.Duration) {

	zoneName := strings.TrimSpace(timeZone)
	var zoneOffset time.Duration

	if strings.Contains(timeZone, "+") || strings.Contains(timeZone, "-") || strings.Contains(timeZone, " ") {

		offsetIndex := strings.IndexAny(timeZone, "+- ")
		zoneName = timeZone[:offsetIndex]
		offsetString := timeZone[offsetIndex+1:]

		offsetString = strings.Replace(offsetString, ":", "", -1)
		if len(offsetString) >= 6 {
			zoneOffset += time.Second * time.Duration(InterfaceToInt(offsetString[4:6]))
		}
		if len(offsetString) >= 4 {
			zoneOffset += time.Minute * time.Duration(InterfaceToInt(offsetString[2:4]))
		}
		if len(offsetString) >= 2 {
			zoneOffset += time.Hour * time.Duration(InterfaceToInt(offsetString[0:2]))
		}
		if len(offsetString) == 1 {
			zoneOffset += time.Hour * time.Duration(InterfaceToInt(offsetString[0:1]))
		}

		if timeZone[offsetIndex] == '-' {
			zoneOffset = -zoneOffset
		}
	} else {
		zoneOffset = time.Hour * time.Duration(InterfaceToInt(timeZone))
		zoneName = strings.Replace(zoneName, InterfaceToString(InterfaceToInt(timeZone)), "", -1)
	}
	// set zone name to UTC when zone name not specified
	if zoneName == "" {
		zoneName = "UTC"
	}

	return zoneName, zoneOffset
}

// MakeUTCTime returns given time in UTC-00:00 timezone (current time.Locale() ignores)
// and offset from GMT in duration format
func MakeUTCTime(inTime time.Time, timeZone string) (time.Time, time.Duration) {
	UTCTime := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), inTime.Minute(), inTime.Second(), inTime.Nanosecond(), time.UTC)

	zoneName, zoneOffset := ParseTimeZone(timeZone)

	if zoneUTCOffset, present := TimeZones[zoneName]; present {
		zoneOffset += time.Duration(zoneUTCOffset)
	} else {
		return time.Time{}, time.Duration(0)
	}

	UTCTime = UTCTime.Add(-zoneOffset)

	return UTCTime, zoneOffset
}

// MakeUTCOffsetTime returns given time in UTC offset timezone (i.e. PST = UTC-08:00, current time.Locale() ignores)
// and offset from GMT in duration format
func MakeUTCOffsetTime(inTime time.Time, timeZone string) (time.Time, time.Duration) {
	zoneName, zoneOffset := ParseTimeZone(timeZone)

	if timeZoneOffset, present := TimeZones[zoneName]; present {
		zoneOffset += timeZoneOffset
	} else {
		return time.Time{}, time.Duration(0)
	}

	resultTime := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), inTime.Minute(), inTime.Second(), inTime.Nanosecond(), time.UTC)
	resultTime = resultTime.Add(zoneOffset)

	return resultTime, zoneOffset
}

// MakeTZTime returns given time in specified timezone (current time.Locale() ignores)
// and offset from GMT in duration format
func MakeTZTime(inTime time.Time, timeZone string) (time.Time, time.Duration) {
	// TODO: move to using time.LoadLocation() instead of tz to support
	//     areas using Daylight Savings Time : https://www.pivotaltracker.com/story/show/108345622
	zoneName, zoneOffset := ParseTimeZone(timeZone)

	if timeZoneOffset, present := TimeZones[zoneName]; present {
		if zoneOffset.Hours() != 0 {
			zoneName = fmt.Sprintf("%s%+0.2d", zoneName, InterfaceToInt(zoneOffset.Hours()))
		}
		zoneOffset += timeZoneOffset
	} else {
		return time.Time{}, time.Duration(0)
	}
	locale := time.FixedZone(zoneName, InterfaceToInt(zoneOffset.Seconds()))

	resultTime := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), inTime.Minute(), inTime.Second(), inTime.Nanosecond(), locale)
	resultTime = resultTime.Add(zoneOffset)
	resultTime = resultTime.In(locale)

	return resultTime, zoneOffset
}

// TimeToUTCTime returns a given time in UTC timezone (converts time.Locale() to offset for UTC time)
func TimeToUTCTime(inTime time.Time) time.Time {
	result := time.Date(inTime.Year(), inTime.Month(), inTime.Day(), inTime.Hour(), inTime.Minute(), inTime.Second(), inTime.Nanosecond(), time.UTC)

	zoneName, offset := inTime.Zone()
	if zoneOffset, present := TimeZones[zoneName]; present {
		result = result.Add(-zoneOffset)
	}

	result = result.Add(time.Second * time.Duration(-offset))

	return result
}

// SetTimeZoneName set name to zone for given time object
func SetTimeZoneName(inTime time.Time, zoneName string) time.Time {

	_, offset := inTime.Zone()
	loc := time.FixedZone(zoneName, offset)

	return inTime.In(loc)
}
