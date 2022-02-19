package ctx

import (
	"runtime"
	"strings"
)

var prefix = initPrefix()

func initPrefix() string {
	_, file, _, _ := runtime.Caller(0)
	return strings.TrimSuffix(file, `visitor/ctx/init.go`)
}

var stackBottom = initStackBottom()

func initStackBottom() string {
	_, file, _, _ := runtime.Caller(2)
	return file
}
