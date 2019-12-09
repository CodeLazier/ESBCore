package server

import (
	"common/task/log"
	"common/task/message"
	"common/task/worker"
	"common/task/yerrors"
	"sync"
)

// get next message if worker is ready
func (t *Server) GetNextMessageGoroutine(groupName string) {
	log.TaskLog.WithField("goroutine", "GetNextMessage").Debug("start")
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
				log.TaskLog.WithField("goroutine", "GetNextMessage").Error("get msg error, ", err)
			}
			continue
		}
		log.TaskLog.WithField("goroutine", "GetNextMessage").Infof("new msg %+v", msg)
		t.msgChan <- msg
	}

	t.getMessageGoroutineStopChan <- struct{}{}
	log.TaskLog.WithField("goroutine", "GetNextMessage").Debug("stop")

}

// start worker to run
func (t *Server) WorkerGoroutine(groupName string) {
	log.TaskLog.WithField("goroutine", "worker").Debug("start")

	workerMap, _ := t.workerGroup[groupName]
	waitWorkerWG := sync.WaitGroup{}

	for msg := range t.msgChan {
		go func(msg message.Message) {
			defer func() {
				e := recover()
				if e != nil {
					log.TaskLog.WithField("goroutine", "worker").Errorf("run worker[%s] panic %v", msg.WorkerName, e)
				}
			}()

			defer func() { go t.MakeWorkerReady() }()

			w, ok := workerMap[msg.WorkerName]
			if !ok {
				log.TaskLog.WithField("goroutine", "worker").Error("not found worker", msg.WorkerName)
				return
			}

			waitWorkerWG.Add(1)
			defer waitWorkerWG.Done()

			result, err := t.GetResult(msg.Id)
			if err != nil {
				if yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
					result = message.NewResult(msg.Id)
				} else {
					log.TaskLog.WithField("goroutine", "worker").
						Error("get result error ", err)
					result = message.NewResult(msg.Id)
				}
			}

			t.workerGoroutine_RunWorker(w, &msg, &result)

		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
	log.TaskLog.WithField("goroutine", "worker").Debug("stop")

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
	log.TaskLog.WithField("goroutine", "worker").Errorf("run worker[%s] error %s", msg.WorkerName, err)

	if ctl.CanRetry() {
		result.Status = message.ResultStatus.WaitingRetry
		ctl.RetryCount -= 1
		msg.TaskCtl = ctl
		log.TaskLog.WithField("goroutine", "worker").Infof("retry task %s", msg)
		ctl.SetError(nil)

		goto RUN
	} else {
		result.Status = message.ResultStatus.Failure
		t.workerGoroutine_SaveResult(*result)

	}

}

func (t *Server) workerGoroutine_SaveResult(result message.Result) {
	log.TaskLog.WithField("goroutine", "worker").
		Debugf("save result %+v", result)
	err := t.SetResult(result)
	if err != nil {
		log.TaskLog.WithField("goroutine", "worker").Errorf("save result error ", err)
	}
}
