package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	TServer "common/task/server"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	. "common"
	"common/helper"
	"common/task/log"
	"common/work"
)

var (
	logger     *zap.Logger
	config     = &Config{}
	R          = sync.Mutex{}
	TaskServer *TServer.Server
)

type Worker struct {
	Concurrency int    `yaml:"concurrency"`
	Names       string `yaml:"names"`
}

type General struct {
	helper.LogParams `yaml:"Log"`
	work.CmdQueue    `yaml:"CmdQueue"`
	Worker           `yaml:"Worker"`
	Caller           `yaml:"Caller"`
}

type Config struct {
	General `yaml:"General"`
}

type Caller struct {
	TLS string `yaml:"TLS"`
}

func readConfigFile() bool {
	R.Lock()
	defer R.Unlock()
	fmt.Println("Ready to read configuration file...")
	content, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		fmt.Println("Read conf.yaml is failed.", err)
		return false
	}

	if err = yaml.Unmarshal(content, config); err != nil {
		fmt.Println("conf.yaml is format invalid.", err)
		return false
	}

	return work.LoadAllCfg(content) == nil
}

func initTaskServer() {
	broker := work.NewTaskBroker("127.0.0.1", "6379", "", 0, 0)
	backend := work.NewTaskBackend("127.0.0.1", "6379", "", 0, 0)

	TaskServer = work.InitTaskServer(broker, backend, -1, -1)
	TaskServer.Add(ESBTaskGroupName+config.Worker.Names, ESBTaskFuncName, work.MainEnter)
	TaskServer.Run(ESBTaskGroupName+config.Worker.Names, 1)

}

//func initMachinery() *machinery.Worker {
//
//	conf := &MConfig.Config{
//		Broker:          config.Broker,
//		DefaultQueue:    config.Tag, //接受消费的队列,publish可以通过RoutingKey来指定交给那个消费worker进行处理
//		ResultBackend:   config.Backend,
//		ResultsExpireIn: config.ResultsExpireIn,
//		//Redis: &MConfig.RedisConfig{
//		//	MaxIdle:                3,
//		//	IdleTimeout:            240,
//		//	ReadTimeout:            15,
//		//	WriteTimeout:           15,
//		//	ConnectTimeout:         15,
//		//	DelayedTasksPollPeriod: 20,
//		//},
//	}
//
//	cleanup, err := helper.SetupTracer(config.Tag)
//	if err != nil {
//		logger.Fatal("Unable to instantiate a tracer:", zap.Error(err))
//	}
//	defer cleanup()
//
//	if server, err := machinery.NewServer(conf); err != nil {
//		logger.Error("Server is failed.", zap.Error(err))
//	} else {
//
//		if err = server.RegisterTask("MainEnter", work.MainEnter); err != nil {
//			logger.Error("RegisterTasks is failed.", zap.Error(err))
//		} else {
//
//			worker := server.NewWorker(config.Tag, config.Concurrency) //consumerTag 仅起个名字方便调试
//
//			//错误回调
//			errorhandler := func(err error) {
//				logger.Error("Work is error.", zap.Error(err))
//			}
//
//			//执行前回调
//			pretaskhandler := func(signature *tasks.Signature) {
//				//
//			}
//
//			//执行后回调
//			posttaskhandler := func(signature *tasks.Signature) {
//				//
//			}
//
//			//register handle
//			worker.SetPostTaskHandler(posttaskhandler)
//			worker.SetErrorHandler(errorhandler)
//			worker.SetPreTaskHandler(pretaskhandler)
//
//			if err = worker.Launch(); err != nil {
//				logger.Error("Work launch is error.", zap.Error(err))
//			}
//
//			return worker
//		}
//	}
//	return nil
//}

func main() {
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt, syscall.SIGTERM)

	readConfigFile()
	logger = helper.NewAdapterLogger(config.LogPath+"/ntworker.log", config.LogSize, config.LogMaxAge, config.LogLevel).Logger
	log.TaskLog=logger
	defer logger.Sync()

	initTaskServer()

	<-chanSignal
	TaskServer.Shutdown(context.Background())
}
