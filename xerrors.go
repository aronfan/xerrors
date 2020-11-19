package xerrors

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	defaultSysInternalError = int64(-12345678)
)

// Xerror - error extension
type Xerror struct {
	err     error
	retcode int64
	stack   string
	message string
}

// SetSysInternalError - set the default system internal error code
func SetSysInternalError(rc int32) {
	defaultSysInternalError = int64(rc)
}

// Wrap - wrap the stderror
func Wrap(err error) *Xerror {
	switch err.(type) {
	case *Xerror:
		e := &Xerror{err: err, retcode: defaultSysInternalError}
		e.withCaller(2)
		return e
	case nil:
		panic(-1)
	default:
		e := &Xerror{err: err, retcode: defaultSysInternalError}
		e.withCaller(2)
		return e
	}
}

// Int - get the retcode
func Int(err error) int32 {
	curr := err
	for {
		switch t := curr.(type) {
		case *Xerror:
			switch t.err.(type) {
			case *Xerror:
				curr = t.err
			default:
				return int32(t.retcode)
			}
		case nil:
			return 0
		default:
			return int32(defaultSysInternalError)
		}
	}
}

// Error - the error string
func (e *Xerror) Error() string {
	stack := ""
	curr := e
	for {
		stack += curr.stack
		switch t := curr.err.(type) {
		case *Xerror:
			curr = t
			stack += " -> "
		default:
			goto LOOPEXIT
		}
	}
LOOPEXIT:
	return fmt.Sprintf("stderr:[%s] retcode:[%d] stack:[%s] message:[%s]", curr.err.Error(), curr.retcode, stack, curr.message)
}

// WithInt - set the retcode
func (e *Xerror) WithInt(rc int32) *Xerror {
	e.retcode = int64(rc)
	return e
}

// WithMessage - set the message
func (e *Xerror) WithMessage(message string) *Xerror {
	e.message = message
	return e
}

func (e *Xerror) withCaller(skip int) *Xerror {
	_, file, line, _ := runtime.Caller(skip)
	index := strings.LastIndex(file, "/")
	if index != -1 {
		file = file[index+1:]
	}
	e.stack = fmt.Sprintf("%s:%d", file, line)
	return e
}
