package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	NTCommon "common"
	"common/helper"
	"common/task/log"
	"common/work"
	"go.uber.org/zap"

	YServer "common/task/server"

	graphite "github.com/cyberdelia/go-metrics-graphite"
	"github.com/rcrowley/go-metrics"
	RpcServer "github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"gopkg.in/yaml.v2"
)

var (
	config     = &Config{}
	R          = sync.Mutex{}
	logger     *zap.Logger
	TaskServer *YServer.Server
	TaskClient *YServer.Client
)

type General struct {
	helper.LogParams `yaml:"Log"`
	work.CmdQueue    `yaml:"CmdQueue"`
	Leader           `yaml:"Leader"`
	Caller           `yaml:"Caller"`
}

type Config struct {
	General `yaml:"General"`
	Rpc     `yaml:"Rpc"`
}

type RegisterServer struct {
	Addr         string   `yaml:"addr"`
	EtcdAddr     []string `yaml:"etcdAddrs"`
	GraphiteAddr string   `yaml:"graphiteAddr"`
}

type Rpc struct {
	RegisterServer `yaml:"RegisterServer"`
	TLS            `yaml:"TLS"`
}

type TLS struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type Caller struct {
	Routing []string `yaml:"routing"`
}

type Leader struct {
	Names []string `yaml:"names"`
}

func startMetrics() {
	metrics.RegisterRuntimeMemStats(metrics.DefaultRegistry)
	go metrics.CaptureRuntimeMemStats(metrics.DefaultRegistry, time.Second)

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	go graphite.Graphite(metrics.DefaultRegistry, 1e9, "rpcx.services.host.127_0_0_1", addr)
}

func addRegistryEtcdPlugin(s *RpcServer.Server) {

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + config.Addr,
		EtcdServers:    config.EtcdAddr,
		BasePath:       NTCommon.RPCEtcdRegisteredPath,
		Metrics:        metrics.NewRegistry(), //监测数据
		UpdateInterval: time.Minute,           //自动更新
	}
	err := r.Start()
	if err != nil {
		logger.Fatal("Register Etcd is failed.", zap.Error(err))
	}
	s.Plugins.Add(r)
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

	return true
}

//可重入,不阻塞客户端的后续调用(Client已使用go调用)
//该函数返回即回传res到client,如果处理中使用goroutine,注意要阻塞此函数返回
func rpcCallback(ctx context.Context, req *NTCommon.ESBRequest, res *NTCommon.ESBResponse) error {
	if req.Enqueue {
		r, err := sendTask(req)
		if err != nil {
			logger.Error("Call is error.", zap.Error(err))
			res.SetError(req, err)
		} else {
			//注意,保持res指针地址,不要覆盖
			*res = *r //res.AssignForRes(r, r.Result)
		}
	} else {
		if resResult, err := work.MainEnterDirect(ctx, req); err != nil {
			res.SetError(req, err)
		} else {
			//注意,保持res指针地址,不要覆盖
			//res.AssignForRes(resResult, resResult.Result)
			*res = *resResult
		}
	}

	return nil
}

func sendTask(req *NTCommon.ESBRequest) (*NTCommon.ESBResponse, error) {
	var err error
	//目前僅支持隨機算法
	var nw string
	nc:=len(config.Names)
	if nc<=0{
		return nil,errors.New("Task worker name is null")
	}else if nc==1{
		nw=config.Names[0]
	}else {
		nw = config.Names[rand.Intn(nc)]
	}
	//if taskId, err := TaskClient.SetTaskCtl(TaskClient.RetryCount, 0).Send(nw, "MainEnter", req); err == nil {
	if taskId, err := TaskClient.Send(NTCommon.ESBTaskGroupName+nw, NTCommon.ESBTaskFuncName, req); err == nil {
		//waiting result
		if result, err := TaskClient.GetResult(taskId, 2*time.Hour, 30*time.Millisecond); err == nil {
			if result.IsSuccess() {
				if resJson,err := result.GetString(0); err != nil {
					return nil, err
				} else {
					if resJson==""{
						if errStr,err:=result.GetString(1);err!=nil{
							return nil,err
						}else {
							return nil, errors.New(errStr)
						}
					}
					res := &NTCommon.ESBResponse{}
					if err = res.Unmarshal([]byte(resJson)); err != nil {
						return nil, err
					} else {
						return res, nil
					}
				}
			} else {
				return nil, errors.New("Result get is faild.")
			}
		}
	}

	return nil, err
}

//func sendTask(ctx context.Context, routKey string, req *NTCommon.ESBRequest, eta *time.Time) (*NTCommon.ESBResponse, error) {
//
//	if ServerMac == nil {
//		return nil, errors.New("Server is null")
//	}
//
//	span, ctx := opentracing.StartSpanFromContext(ctx, "PushTask")
//	defer span.Finish()
//	id := strconv.FormatInt(req.ID, 10)
//	span.SetBaggageItem("batch.id", id)
//	span.LogFields(opentracing_log.String("batch.id", id))
//
//	//序列化
//	bytes, err := req.Marshal()
//
//	if err != nil {
//		//logger.Error("Request serialize is failed.", zap.Error(err))
//		return nil, err
//	}
//
//	task := tasks.Signature{
//		Name:       "MainEnter",
//		UUID:       id,
//		RoutingKey: config.Tag,
//		ETA:        eta,
//		RetryCount: 0, //不重复
//		Args: []tasks.Arg{
//			{
//				Type:  "string",
//				Value: string(bytes), //注意,Signature再发送给队列时会再次序列化
//			},
//			{
//				Type:  "[]string",
//				Value: config.Caller.Routing, //strings.Join(config.Caller.Routing,"&"), //注意,Signature再发送给队列时会再次序列化
//			},
//		},
//	}
//
//	//ctx_e:=context.WithValue(ctx,"extends",config.Extends)
//
//	//发布到消息队列
//	result, err := ServerMac.SendTaskWithContext(ctx, &task)
//	if err != nil {
//		//logger.Error("Send task queue is error", zap.Error(err))
//		return nil, err
//	}
//
//
//	//定时轮训检测是否有结果
//	if value, err := result.GetWithTimeout(time.Duration(time.Hour * time.Duration(config.WaitTimeout)) , time.Duration(time.Millisecond * 30)); err != nil {
//		return nil, err
//	} else {
//		//默认日志Debug等级会输出结果,这里不必再次输出
//		resJson := tasks.HumanReadableResults(value)
//		res:=&NTCommon.ESBResponse{}
//		if err=res.Unmarshal([]byte(resJson));err!=nil{
//			//logger.Error("Marshal is error", zap.Error(err))
//			return nil, err
//		} else {
//			return res, nil
//		}
//	}
//}

func initTaskServer() {
	broker := work.NewTaskBroker("127.0.0.1", "6379", "", 0, 10)
	backend := work.NewTaskBackend("127.0.0.1", "6379", "", 0, 10)

	TaskServer = work.InitTaskServer(broker, backend, -1, -1)
	_c := TaskServer.GetClient()
	TaskClient = &_c

	for _,nw:=range config.Leader.Names {
		TaskServer.Add(NTCommon.ESBTaskGroupName, nw, work.MainEnter)
	}
}

//func initMachinery() *machinery.Server {
//
//	cleanup, err := helper.SetupTracer("sender")
//	if err != nil {
//		logger.Fatal("Unable to instantiate a tracer:", zap.Error(err))
//	}
//	defer cleanup()
//
//	confMac := &MConfig.Config{
//		Broker:          config.Broker,
//		ResultBackend:   config.Backend,
//		ResultsExpireIn: config.ResultsExpireIn * 3600000, //3600秒 1小时,这里没有使用,因为Publish不负责保存结果
//		//Redis: &MConfig.RedisConfig{
//		//	MaxIdle:                10,
//		//	IdleTimeout:            240,
//		//	ReadTimeout:            15,
//		//	WriteTimeout:           15,
//		//	ConnectTimeout:         15,
//		//	DelayedTasksPollPeriod: 20,
//		//},
//	}
//
//	if s, err := machinery.NewServer(confMac); err != nil {
//		logger.Panic("Connect server is failed.", zap.Error(err))
//	} else {
//		//注册
//		if err = s.RegisterTask("MainEnter", work.MainEnter); err != nil {
//			logger.Panic("Register tasks is failed.", zap.Error(err))
//		}
//		return s
//	}
//	return nil
//}

func initRpcServer() (*RpcServer.Server, error) {
	var server *RpcServer.Server
	if strings.TrimSpace(config.CertFile) != "" && strings.TrimSpace(config.KeyFile) != "" {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		server = RpcServer.NewServer(RpcServer.WithTLSConfig(tlsConfig))
	} else {
		server = RpcServer.NewServer()
	}

	opt := RpcServer.AllowAllCORSOptions()
	server.SetCORS(opt)

	if strings.TrimSpace(config.GraphiteAddr) != "" {
		p := serverplugin.NewMetricsPlugin(metrics.DefaultRegistry)
		server.Plugins.Add(p)
		startMetrics()
	}

	p2 := serverplugin.OpenTracingPlugin{}
	server.Plugins.Add(p2)

	server.RegisterOnShutdown(func(s *RpcServer.Server) {
		_ = server.UnregisterAll()
	})

	if len(config.EtcdAddr) != 0 {
		addRegistryEtcdPlugin(server)
	}

	if err := server.RegisterFunctionName(NTCommon.CommonService, NTCommon.ESBRequestFunction, rpcCallback, ""); err != nil {
		return nil, err
	}

	go func() {
		if err := server.Serve("tcp", config.Addr); err != nil {
			logger.Fatal("Error", zap.Error(err))
		}
	}()

	return server, nil
}

func main() {

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt, syscall.SIGTERM)

	rand.Seed(time.Now().UnixNano())
	readConfigFile()

	logger = helper.NewAdapterLogger(config.LogPath+"/ntleader.log", config.LogSize, config.LogMaxAge, config.LogLevel).Logger
	log.TaskLog=logger
	defer logger.Sync()

	initTaskServer()

	if _, err := initRpcServer(); err != nil {
		logger.Fatal("Error", zap.Error(err))
	}

	<-chanSignal
	//err := server.Shutdown(context.Background())
	//if err != nil {
	//	logger.Panic("Shutdown Error.", zap.Error(err))
	//}
}
