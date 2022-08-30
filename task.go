package goworkerpool

import "github.com/pkg/errors"

var ERROR_PARAM_ERROR = errors.New("Pool len illegal")

type TaskInput struct {
	Id        string
	Guid      string
	InputData interface{}
}

type TaskResult struct {
	Id         string //回调ID
	Guid       string
	OutputData interface{}
	Err        error
	Boolret    bool
}

//
type WorkFunc func(task TaskInput) *TaskResult
