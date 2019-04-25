// Copyright 2019 Ottemo. All rights reserved.

/*
Package logger is a default implementation of InterfaceLogger declared in "github.com/ottemo/commerce/env" package.

Default logger is pretty simple implementation, which takes an message and puts it into a "storage" (a file with specific
name in this case). If for some reason message can not be places in file (file access denied, etc.) message will be
printed to stdout. Message time (in RFC3339 format) and specified prefix adds to message before output.
*/
package logger
