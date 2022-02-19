package ctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Err   error
	Stack []string
	Tags  Tags
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"err":   e.Error(),
		"stack": e.Stack,
		"tags":  e.Tags,
	})
}

func (err Error) asPtr() *Error {
	return &err
}

func (err Error) Error() string {
	return err.Err.Error()
}

func scanError(args []interface{}) (Error, bool) {
	var err Error
	for _, e := range args {
		switch e := e.(type) {
		case error:
			if errors.As(e, &err) {
				return err, true
			}
		}
	}
	return err, false
}

func Wrap(c C, parent error) error {
	if parent == nil {
		return nil
	}
	var err Error
	if errors.As(parent, &err) {
		return &err
	}
	return &Error{
		Err:   parent,
		Stack: stack(0),
		Tags:  getTags(c),
	}
}

func Wrapf(c C, f string, args ...interface{}) Error {
	var err Error
	parent := fmt.Errorf(f, args...)
	err, found := scanError(args)
	if found {
		err.Err = parent
		return err
	}
	return Error{
		Err:   parent,
		Stack: stack(1),
		Tags:  getTags(c),
	}
}

func callerFull(above int) (pkg, fname, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(above+2, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	if frame.PC == 0 {
		return
	}
	ok = true
	i := strings.LastIndex(frame.Function, ".")
	if i >= 0 {
		pkg = frame.Function[0:i]
		fname = frame.Function[i+1:]
	}
	file = frame.File
	line = frame.Line
	return
}

func stack(i int) []string {
	out := []string{}
	rpc := make([]uintptr, 1)
	for j := i + 1; j < 30; j++ {
		_ = runtime.Callers(j+1, rpc[:])
		frame, _ := runtime.CallersFrames(rpc).Next()
		//_, file, line, ok := runtime.Caller(j)
		//if !ok {
		//	break
		//}
		if frame.File == stackBottom || frame.Function == `runtime.goexit` {
			break
		}
		file := strings.TrimLeft(frame.File, prefix)
		src := fmt.Sprintf("%s:%d   %s", file, frame.Line, frame.Func.Name())
		out = append(out, src)
	}
	return out
}
