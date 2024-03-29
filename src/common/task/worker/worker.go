package worker

import (
	"errors"
	"fmt"
	"common/task/controller"
	"common/task/log"
	"common/task/message"
	"common/task/util"
	"go.uber.org/zap"

	"reflect"
)

type WorkerInterface interface {
	Run(ctl *controller.TaskCtl, funcArgs []string, result *message.Result) error
	WorkerName() string
}

type FuncWorker struct {
	F    interface{}
	Name string
}

func (f FuncWorker) Run(ctl *controller.TaskCtl, funcArgs []string, result *message.Result) error {
	return runFunc(f.F, ctl, funcArgs, result)
}
func (f FuncWorker) WorkerName() string {
	return f.Name
}

func runFunc(f interface{}, ctl *controller.TaskCtl, funcArgs []string, result *message.Result) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			result.Status = message.ResultStatus.Failure
			t, ok := e.(error)
			if ok {
				err = t
			} else {
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}
	}()
	funcValue := reflect.ValueOf(f)
	funcType := reflect.TypeOf(f)
	var inStart = 0
	var inValue []reflect.Value
	if funcType.NumIn() > 0 && funcType.In(0) == reflect.TypeOf(&controller.TaskCtl{}) {
		inStart = 1
	}

	inValue, err = util.GetCallInArgs(funcValue, funcArgs, inStart)
	if err != nil {
		return
	}
	if inStart == 1 {
		inValue = append(inValue, reflect.Value{})
		copy(inValue[1:], inValue)
		inValue[0] = reflect.ValueOf(ctl)

	}

	funcOut := funcValue.Call(inValue)

	if ctl.GetError() != nil {
		err = ctl.GetError()
	} else {
		result.Status = message.ResultStatus.Success
		if len(funcOut) > 0 {
			re, err2 := util.GoValuesToYJsonSlice(funcOut)
			if err2 != nil {
				log.TaskLog.Error("Error:",zap.Error(err2))
			} else {
				result.FuncReturn = re
			}
		}
	}

	return
}
