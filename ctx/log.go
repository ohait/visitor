package ctx

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var Out = json.NewEncoder(os.Stdout)

func caller(i int) string {
	_, file, line, _ := runtime.Caller(i + 1)
	file = strings.TrimLeft(file, prefix)
	return fmt.Sprintf("%s:%d", file, line)
}

type Line struct {
	Message string                     `json:"message"`
	Leval   string                     `json:"level"`
	Time    time.Time                  `json:"time"`
	Tags    map[string]json.RawMessage `json:"tags"`
	Src     string                     `json:"src"`
	Error   error                      `json:"err"`
}

func Log(c C) Line {
	return Line{
		Src:  caller(1),
		Time: time.Now(),
		Tags: getTags(c),
	}
}

func (l Line) Tag(key string, val interface{}) Line {
	j, err := json.Marshal(val)
	if err != nil {
		j, _ = json.Marshal(err.Error())
	}
	l.Tags[key] = j
	return l
}

func (l Line) Debugf(args ...interface{}) {
	l.Leval = "debug"
	l.logf(args...)
}

func (l Line) Warnf(args ...interface{}) {
	l.Leval = "warn"
	l.logf(args...)
}

func (l Line) Errf(args ...interface{}) {
	l.Leval = "warn"
	l.logf(args...)
}

func (l Line) Err(err error) Line {
	l.Error = Wrap(nil, err)
	return l
}

func (l Line) logf(args ...interface{}) {
	msg := ""
	switch len(args) {
	case 0: // empty message? :shrug:
	case 1:
		msg = fmt.Sprintf("%s", args[0])
		if l.Error == nil {
			switch v := args[0].(type) {
			case error:
				l = l.Err(v)
			}
		}
	default:
		msg = fmt.Sprintf(args[0].(string), args[1:]...)
		if l.Error == nil {
			err, found := scanError(args)
			if found {
				l.Error = err.asPtr()
			}
		}
	}
	l.Message = msg

	err := Out.Encode(l)
	if err != nil {
		panic(err) // if we can't log, we panic
	}
}
