package yerrors

import (
	"fmt"
)

const (
	ErrTypeEmptyQuery      = 1
	ErrTypeUnsupportedType = 2
	ErrTypeOutOfRange      = 3
	ErrTypeNilResult       = 4
	ErrTypeTimeOut         = 5
)

func IsEqual(err error, errType int) bool {
	yerr, ok := err.(TaskError)
	if !ok {
		return ok
	}
	if yerr.Type() == errType {
		return true
	}
	return false
}

type TaskError interface {
	Error() string
	Type() int
}

type ErrEmptyQuery struct {
}

func (e ErrEmptyQuery) Error() string {
	return "Task: empty query"
}

func (e ErrEmptyQuery) Type() int {
	return ErrTypeEmptyQuery
}

type ErrUnsupportedType struct {
	T string
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("Task: UnsupportedType: %s", e.T)
}

func (e ErrUnsupportedType) Type() int {
	return ErrTypeUnsupportedType
}

type ErrOutOfRange struct {
}

func (e ErrOutOfRange) Error() string {
	return "Task: index out of range"
}

func (e ErrOutOfRange) Type() int {
	return ErrTypeOutOfRange
}

type ErrNilResult struct {
}

func (e ErrNilResult) Error() string {
	return "Task: nil result"
}

func (e ErrNilResult) Type() int {
	return ErrTypeNilResult
}

type ErrTimeOut struct {
}

func (e ErrTimeOut) Error() string {
	return "Task: timeout"
}

func (e ErrTimeOut) Type() int {
	return ErrTypeTimeOut
}
