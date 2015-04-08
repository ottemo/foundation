// Copyright 2014 The Ottemo Authors. All rights reserved.

/*
Package ini is a default implementation of InterfaceIniConfig declared in "github.com/ottemo/foundation/env" package.

Ini config is a config.ini file where startup information located. Ini file could be separated on a sections (refer to
https://github.com/vaughan0/go-ini) for ini file lookup. Ini config takes a value from ini file "current" section and
if not found in there - looking for a global section.

Special section [ConstTestSectionName] ("test") used for "test mode" application startup. To start application in that
mode [ConstCmdArgTestFlag] "--test" should be used.

    Example 1:
    ----------
        if iniConfig := env.GetIniConfig(); iniConfig != nil {
            if iniValue := iniConfig.GetValue("db.sqlite3.uri", uri); iniValue != "" {
                uri = iniValue
            }
        }

    Example 2:
    ----------
        uri := env.IniValue("db.sqlite3.uri")
*/
package ini
