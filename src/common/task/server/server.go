package server

import (
	"context"
	"fmt"
	"common/task/backends"
	"common/task/brokers"
	"common/task/config"
	"common/task/controller"
	"common/task/log"
	"common/task/message"
	"common/task/worker"
	"go.uber.org/zap"

	"reflect"
	"sync"
)

// [workerName]worker
type workerMap map[string]worker.WorkerInterface

type Server struct {
	sync.Map

	workerGroup map[string]workerMap // [groupName]workerMap

	broker  brokers.BrokerInterface
	backend backends.BackendInterface

	workerReadyChan chan struct{}
	msgChan         chan message.Message

	getMessageGoroutineStopChan chan struct{}
	workerGoroutineStopChan     chan struct{}

	safeStopChan chan struct{}

	// config
	StatusExpires int // second, -1:forever
	ResultExpires int // second, -1:forever
}

func NewServer(c config.Config) Server {

	g := make(map[string]workerMap)
	return Server{
		workerGroup:   g,
		broker:        c.Broker,
		backend:       c.Backend,
		safeStopChan:  make(chan struct{}),
		StatusExpires: c.StatusExpires,
		ResultExpires: c.ResultExpires,
	}
}

func (t *Server) GetQueryName(groupName string) string {
	return "ESB.Core.Task.Pushed:" + groupName
}

func (t *Server) MakeWorkerReady() {
	defer func() {
		recover()
	}()
	t.workerReadyChan <- struct{}{}
}

func (t *Server) Run(groupName string, numWorkers int) {

	workerMap, ok := t.workerGroup[groupName]
	if !ok {
		panic("not find group '" + groupName + "'")
	}
	_, ok = t.Load("isRunning")
	if ok {
		panic("Running multiple groups is not supported, this feature is already under development")
	}
	t.Store("isRunning", struct{}{})

	if t.broker != nil {
		if t.broker.GetPoolSize() <= 0 {
			t.broker.SetPoolSize(1)
		}
		t.broker.Activate()
	}
	if t.backend != nil {
		if t.backend.GetPoolSize() <= 0 {
			t.backend.SetPoolSize(numWorkers)
		}
		t.backend.Activate()
	}

	log.TaskLog.Info("Start",zap.String("group", groupName),zap.Int("numWorkers", numWorkers))

	log.TaskLog.Info("group workers:")
	for name := range workerMap {
		log.TaskLog.Info("  - " + name)
	}

	t.workerReadyChan = make(chan struct{}, numWorkers)
	t.msgChan = make(chan message.Message, numWorkers)

	t.getMessageGoroutineStopChan = make(chan struct{}, 1)
	go t.GetNextMessageGoroutine(groupName)

	t.workerGoroutineStopChan = make(chan struct{}, 1)
	go t.WorkerGoroutine(groupName)

	for i := 0; i < numWorkers; i++ {
		t.MakeWorkerReady()
	}

}

func (t *Server) safeStop() {
	log.TaskLog.Info("waiting for incomplete tasks ")

	// stop get message goroutine
	close(t.workerReadyChan)
	t.Store("isStop", struct{}{})
	<-t.getMessageGoroutineStopChan

	// stop worker goroutine
	close(t.msgChan)
	<-t.workerGoroutineStopChan

	close(t.safeStopChan)

}

func (t *Server) Shutdown(ctx context.Context) error {

	go func() {
		t.safeStop()
	}()

	select {
	case <-t.safeStopChan:
	case <-ctx.Done():
		return ctx.Err()
	}

	log.TaskLog.Info("Shutdown!")
	return nil
}

// add worker to group
// w : worker func
func (t *Server) Add(groupName string, workerName string, w interface{}) {
	_, ok := t.workerGroup[groupName]
	if !ok {
		t.workerGroup[groupName] = make(workerMap)
	}
	wType := reflect.TypeOf(w).Kind().String()
	if wType == "func" {
		funcWorker := worker.FuncWorker{
			Name: workerName,
			F:    w,
		}
		t.workerGroup[groupName][workerName] = funcWorker
	} else {
		panic(fmt.Sprintf("worker must be func"))
	}

}

func (t *Server) Next(groupName string) (message.Message, error) {
	return t.broker.Next(t.GetQueryName(groupName))
}

func (t *Server) SetResult(result message.Result) error {
	var exTime int
	if result.IsFinish() {
		exTime = t.ResultExpires
	} else {
		exTime = t.StatusExpires
	}
	if exTime == 0 {
		return nil
	}
	return t.backend.SetResult(result, exTime)
}

// send msg to queue
// t.Send("groupName", "workerName" , 1,"hi",1.2)
//
func (t *Server) Send(groupName string, workerName string, ctl controller.TaskCtl, args ...interface{}) (string, error) {
	var msg = message.NewMessage(ctl)
	msg.WorkerName = workerName
	err := msg.SetArgs(args...)
	if err != nil {
		return "", err
	}

	return msg.Id, t.broker.Send(t.GetQueryName(groupName), msg)

}

func (t *Server) GetResult(id string) (message.Result, error) {
	result := message.NewResult(id)
	return t.backend.GetResult(result.GetBackendKey())
}

func (t *Server) GetClient() Client {
	if t.broker != nil {
		if t.broker.GetPoolSize() <= 0 {
			t.broker.SetPoolSize(10)
		}
		t.broker.Activate()
	}
	if t.backend != nil {
		if t.backend.GetPoolSize() <= 0 {
			t.backend.SetPoolSize(10)
		}
		t.backend.Activate()
	}
	return Client{
		server: t,
		ctl:    controller.NewTaskCtl(),
	}
}
