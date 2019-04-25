// Copyright 2019 Ottemo. All rights reserved.

/*
Package test represents set of test writing helpers and global application tests. It provides a set of functions you can
use for starting Ottemo application in a test mode, prepare randomized data, fill DB with sample data, which is use-full
during GO tests writing.

Package also contains a set of benchmarks and tests which related to whole application rather when particular package.
In order to run them use:
  go test [-tags ...] github.com/ottemo/commerce/tests
  go test -bench . [-tags ...] github.com/ottemo/commerce/tests

(refer to http://golang.org/pkg/testing/ for details)
*/
package test
