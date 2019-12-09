package server

import (
	"sync"

	"common/task/log"
	"common/task/message"
	"common/task/worker"
	"common/task/yerrors"
	"go.uber.org/zap"
)

// get next message if worker is ready
func (t *Server) GetNextMessageGoroutine(groupName string) {
	var msg message.Message
	var err error
	for range t.workerReadyChan {
		_, ok := t.Load("isStop")
		if ok {
			break
		}
		msg, err = t.Next(groupName)

		if err != nil {
			go t.MakeWorkerReady()
			TaskErr, ok := err.(yerrors.TaskError)
			if ok && TaskErr.Type() != yerrors.ErrTypeEmptyQuery {
				log.TaskLog.Error("get msg error, ",zap.Error(err))
			}
			continue
		}
		log.TaskLog.Info("",zap.Any("new msg",msg))
		t.msgChan <- msg
	}

	t.getMessageGoroutineStopChan <- struct{}{}

}

// start worker to run
func (t *Server) WorkerGoroutine(groupName string) {
	workerMap, _ := t.workerGroup[groupName]
	waitWorkerWG := sync.WaitGroup{}

	for msg := range t.msgChan {
		go func(msg message.Message) {
			defer func() {
				e := recover()
				if e != nil {
					log.TaskLog.Error("run worker panic",zap.String("name",msg.WorkerName),zap.Any("e",e))
				}
			}()

			defer func() { go t.MakeWorkerReady() }()

			w, ok := workerMap[msg.WorkerName]
			if !ok {
				log.TaskLog.Error("not found worker",zap.String("name",msg.WorkerName))
				return
			}

			waitWorkerWG.Add(1)
			defer waitWorkerWG.Done()

			result, err := t.GetResult(msg.Id)
			if err != nil {
				if yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
					result = message.NewResult(msg.Id)
				} else {
					log.TaskLog.Error("get result error ", zap.Error(err))
					result = message.NewResult(msg.Id)
				}
			}

			t.workerGoroutine_RunWorker(w, &msg, &result)

		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
}

func (t *Server) workerGoroutine_RunWorker(w worker.WorkerInterface, msg *message.Message, result *message.Result) {

	ctl := msg.TaskCtl

RUN:

	result.SetStatusRunning()
	t.workerGoroutine_SaveResult(*result)

	err := w.Run(&ctl, msg.FuncArgs, result)
	if err == nil {
		result.Status = message.ResultStatus.Success
		t.workerGoroutine_SaveResult(*result)

		return
	}
	log.TaskLog.Error("run worker[%s] error %s",zap.String("name",msg.WorkerName),zap.Error(err))

	if ctl.CanRetry() {
		result.Status = message.ResultStatus.WaitingRetry
		ctl.RetryCount -= 1
		msg.TaskCtl = ctl
		log.TaskLog.Info("retry task ",zap.Any("value",msg))
		ctl.SetError(nil)

		goto RUN
	} else {
		result.Status = message.ResultStatus.Failure
		t.workerGoroutine_SaveResult(*result)

	}

}

func (t *Server) workerGoroutine_SaveResult(result message.Result) {
	log.TaskLog.Debug("save result",zap.Any("Value", result))
	err := t.SetResult(result)
	if err != nil {
		log.TaskLog.Error("save result error ",zap.Error(err))
	}
}
