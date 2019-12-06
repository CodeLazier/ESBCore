package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	NTCommon "common"
	"common/helper"
	"github.com/beinan/fastid"
	"github.com/pkg/errors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/emicklei/go-restful"
	RpcClient "github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type ServiceSelect struct {
	Algorithm int `yaml:"algorithm"`
}

type ServiceDiscovery struct {
	Category int      `yaml:"category"`
	Addrs    []string `yaml:"addrs"`
}

type Rpc struct {
	ServiceDiscovery `yaml:"ServiceDiscovery"`
	ServiceSelect    `yaml:"Select"`
	RPCTLS           `yaml:"TLS"`
}

type General struct {
	helper.LogParams `yaml:"Log"`
}

type TLS struct {
	CAFile         string `yaml:"caFile"`
	ClientCertFile string `yaml:"clientCert"`
	ClientKeyFile  string `yaml:"clientKey"`
}

type RPCTLS struct {
	Enable bool `yaml:"enable"`
}

type MQTTServer struct {
	Enable   bool   `yaml:"enable"`
	URI      string `yaml:"URI"`
	ClientID string `yaml:"clientId"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	//其实是cleanSession,并不是Mqtt中的retain
	Retain bool `yaml:"retain"`
	TLS    `yaml:"TLS"`
}

type Subscribe struct {
	AutoRefresh int      `yaml:"autoRefresh"`
	Topics      []string `yaml:"Topics"`
}

type WebService struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
	TLS    TLS  `yaml:"TLS"`
}

type Config struct {
	MQTTServer `yaml:"MQTTServer"`
	Subscribe  `yaml:"Subscribe"`
	General    `yaml:"General"`
	Rpc        `yaml:"Rpc"`
	WebService `yaml:"WebService"`
}

var (
	config            = &Config{}
	R                 = sync.RWMutex{}
	readConfigDone    = make(chan bool)
	topics            = make(map[string]byte)
	logger            *zap.Logger
	mqttClient        MQTT.Client
	mqttClientAdapter MQTT.Client
	xClient           RpcClient.XClient
)

func NewTLSConfig(cafile, ccfile, ckeyfile string) *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(cafile)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(ccfile, ckeyfile)
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}

	//fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

//如果ntleader离线,要有重连机制
func initRpcClient() RpcClient.XClient {
	//etcdAddr := "localhost:2379"
	//basePath := "/rpcx_test"

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	var sd RpcClient.ServiceDiscovery

	//NewBidirectionalXClient

	if config.Category == 0 {
		//etcd
		sd = RpcClient.NewEtcdV3Discovery(NTCommon.RPCEtcdRegisteredPath, NTCommon.CommonService, config.Addrs, nil)
	} else if config.Category == 1 { //multServer
		var kvp []*RpcClient.KVPair
		for _, addr := range config.Addrs {
			kvp = append(kvp, &RpcClient.KVPair{Key: addr})
		}
		sd = RpcClient.NewMultipleServersDiscovery(kvp)
	} else if config.Category == 2 { //singleServer
		sd = RpcClient.NewPeer2PeerDiscovery("tcp@"+config.Addrs[0], "")
	}

	ops := RpcClient.DefaultOption
	if config.RPCTLS.Enable {
		ops.TLSConfig = conf
	}
	ops.Retries = 1      //否则不会触发
	ops.Heartbeat = true //心跳
	ops.HeartbeatInterval = time.Second * 5

	return RpcClient.NewXClient(NTCommon.CommonService, RpcClient.Failover, RpcClient.RoundRobin, sd, ops)

	//pool:=RpcClient.NewXClientPool(100,"",RpcClient.Failover, RpcClient.RoundRobin, d, options)

	//ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	//for {
	//	reply := &example.Reply{}
	//	err := xclient.Call(context.Background(), "Mul", args, reply)
	//	if err != nil {
	//		log.Printf("failed to call: %v\n", err)
	//		time.Sleep(5 * time.Second)
	//		continue
	//	}

	//log.Printf("%d * %d = %d", args.A, args.B, reply.C)

	//time.Sleep(5 * time.Second)
	//}
}

func callRpcServer(req *NTCommon.ESBRequest, timeout time.Duration) (*NTCommon.ESBResponse, error) {
	var reply NTCommon.ESBResponse
	if xClient != nil {
		ctx := context.Background()
		if timeout > 0 {
			ctx, _ = context.WithTimeout(ctx, timeout)
		}
		err := xClient.Call(ctx, NTCommon.ESBRequestFunction, req, &reply)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Client is null")
	}

	return &reply, nil
}

func readConfigFile() bool {
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

func updateTopics() {
	//已经RLock
	for _, t := range config.Subscribe.Topics {
		var qos int = 2
		var topic string
		var err error
		s := strings.Split(t, ",")
		l := len(s)
		if l >= 2 {
			topic = s[0]
			if qos, err = strconv.Atoi(s[1]); err != nil {
				logger.Warn("Qos is failed", zap.String("Value", s[1]))
			} else {
				if l > 2 {
					logger.Warn("Extension is currently not supported", zap.Skip())
				}
			}
		} else if l == 1 {
			topic = t
		} else {
			logger.Warn("Subscribe is invalid format", zap.String("Value", t))
		}

		if topic != "" {
			if _, ok := topics[topic]; !ok {
				topics[topic] = byte(qos)
				logger.Info("Registered", zap.String("Topic", topic), zap.Int("Qos", qos))
				if mqttClient != nil && mqttClient.IsConnected() {
					if token := mqttClient.Subscribe(topic, byte(qos), onMsgReceived); token.Wait() && token.Error() != nil {
						logger.Error("Subscribe is error", zap.Error(token.Error()))
					}
				}

			}
		}
	}
}

//该事件会阻塞重入,所以尽快完成事件处理防止阻塞后续消息的到来
func onMsgReceived(client MQTT.Client, message MQTT.Message) {
	go func() {
		logger.Debug("Received message on topic", zap.Strings("Content:", []string{message.Topic(),
			strconv.Itoa(int(message.Qos())),
			strconv.Itoa(int(message.MessageID())),
			string(message.Payload()),
		}))

		var token = &struct {
			Token           string `json:"token"`
			CallID          string `json:"callid"`
			SubscribeResult bool   `json:"result"`
		}{SubscribeResult: true}

		payload := message.Payload()
		err := json.Unmarshal(payload, token)
		if err != nil || token.Token == "" {
			logger.Error("Token is invalid or null", zap.Skip())
			return
		}

		authInfo, _ := helper.ParseJWT(token.Token)
		if err != nil {
			logger.Error("Token is unknow or expired", zap.Skip())
			return
		}

		_ = authInfo

		id := fastid.CommonConfig.GenInt64ID()
		res, err := callRpcServer(&NTCommon.ESBRequest{
			ID:        id,
			Topic:     message.Topic(),
			TimeStamp: time.Now(),
			Payload:   string(payload),
		}, 0)

		if err != nil {
			logger.Error("Call RPC Server is Failed", zap.Int64("ID", id),
				zap.String("Topic", message.Topic()),
				zap.Error(err))
		} else {
			if res.IsError() != nil {
				logger.Error("Response is error", zap.Error(res.IsError()))
			} else {
				logger.Debug("Call result is ", zap.Int64("RequestID", id), zap.Any("Response", res))
			}
			if token.SubscribeResult {
				//发布mqtt消息给订阅者或约定的topic
				//try loop?
				looptry := 0
				for !mqttClient.IsConnected() {
					if looptry > 30 {
						break
					}
					looptry++
					select {
					case <-time.After(time.Minute):
					}
				}
				if mqttClient.IsConnected() {
					callID := token.CallID
					if callID == "" {
						callID = strconv.FormatInt(id, 10)
					}
					ptopic := fmt.Sprintf("%s/Result/%s", authInfo["name"], callID)
					if resByte, err := res.Marshal(); err != nil {
						logger.Error("Response marshal is error", zap.Error(err))
					} else if t := mqttClient.Publish(ptopic, byte(2), false, resByte); t.WaitTimeout(time.Minute) && t.Error() != nil {
						logger.Warn("Publish mqtt topic", zap.Error(t.Error()))
					}
				}
			}
		}
	}()
}

func startServerAPI() {
	wsContainer := restful.NewContainer()
	u := AccessEnter{}
	u.RegisterTo(wsContainer)

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-ESB-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: true,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	addr := fmt.Sprintf(":%d", config.WebService.Port)
	server := &http.Server{Addr: addr, Handler: wsContainer}
	//defer server.Close()
	if config.WebService.TLS.ClientCertFile != "" && config.WebService.TLS.ClientKeyFile != "" {
		go server.ListenAndServeTLS(config.WebService.TLS.ClientCertFile, config.WebService.TLS.ClientKeyFile)
	} else {
		go server.ListenAndServe()
	}
}

type AccessEnter struct{}

func (u AccessEnter) RegisterTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/ESB").Consumes(restful.MIME_JSON, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("NT/EDI/v1").To(u.EDIRequest_v1).Writes(NTCommon.ESBResponse{}))

	//ws.Route(ws.POST("Mqtt").To())
	container.Add(ws)
}

func initMqttClientAdapter() {
	connOpts := MQTT.NewClientOptions().AddBroker(config.URI).SetClientID(config.ClientID + "_adapter")

	mqttClientAdapter := MQTT.NewClient(connOpts)
	if token := mqttClientAdapter.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("MQTT client connect error", zap.Error(token.Error()))
	}
}

func (u AccessEnter) requestMqttAdapterFunc(request *restful.Request, response *restful.Response) {
	topic := ""
	body, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		return
	}

	if mqttClientAdapter != nil && mqttClientAdapter.IsConnected() {
		t := mqttClientAdapter.Publish(topic, byte(2), false, body)
		if t.Wait() && t.Error() != nil {
			logger.Error("MQTT client is error", zap.Error(t.Error()))
			return
		}
	}
	//等待?

	response.WriteHeaderAndJson(200, nil, "application/json")
}

func (u AccessEnter) EDIRequest_v1(request *restful.Request, response *restful.Response) {
	defer request.Request.Body.Close()

	body, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		logger.Error("500", zap.Error(err))
		response.WriteHeaderAndJson(500, err.Error(), "application/json")

	} else {

		req := &NTCommon.ESBRequest{Enqueue: true}

		if err = req.Unmarshal(body); err != nil {
			logger.Error("400", zap.Error(err))
			response.WriteHeaderAndJson(400, err.Error(), "application/json")

		} else {
			//init
			req.TimeStamp, req.Topic, req.ID = time.Now(), "NT/EDI/", fastid.CommonConfig.GenInt64ID()

			res, err := callRpcServer(req, 0)
			if err != nil {
				logger.Error("Call RPC Server is Failed", zap.Int64("ID", req.ID),
					zap.String("Topic", req.Topic),
					zap.Error(err))
			} else {
				if res.IsError() != nil {
					logger.Error("Response is error", zap.Error(res.IsError()))
				} else {
					logger.Debug("Call result is ", zap.Int64("RequestID", req.ID), zap.Any("Response", res))
				}
			}
			//res.Err=errors.WithMessage(res.Err,res.Err.Error())
			//bytes,err:=res.Marshal()
			//if err!=nil{
			//	return
			//}
			//writer.Write(bytes)
			response.WriteHeaderAndJson(200, res, "application/json")
			//response.ResponseWriter.Write(bytes)
			//io.WriteString(response.ResponseWriter, "this would be a normal response")
		}
	}
}

func main() {
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt, syscall.SIGTERM)

	//没有竞态
	if !readConfigFile() {
		return
	}

	logger = helper.NewAdapterLogger(config.LogPath+"/ntagent.log", config.LogSize, config.LogMaxAge, config.LogLevel).Logger
	defer logger.Sync()

	if config.WebService.Enable {
		startServerAPI()
	}

	//m,_:=helper.ParseJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTY1NzcyNDI1LCJuYW1lIjoidGVzdCJ9.2LwopohhuBDUo1i-jiK4YLnapUVqQi5XDVK57eBH2QQ")
	//fmt.Println(m)

	//cli,err:=NTConfig.NewEtcdClient()
	//if err!=nil{
	//	logger.Fatal("Connect etcd is failed",zap.Error(err))
	//}

	//s:=[]string{
	//	"NT/Test/#@2",
	//	"NT/EDI/#@2",
	//	"NT/Sample/#",
	//	"Other/API/+@1",
	//	"Ex/Example/Call/#",
	//}
	//
	//NTConfig.SetRegCall(cli,s)
	//
	//regFun,err:=NTConfig.GetRegCall(cli)
	//if err!=nil{
	//	logger.Error("Read regfun from etcd is failed.",zap.Error(err))
	//}
	//
	//logger.Info(strings.Join(regFun,","),zap.Skip())

	//cli.Close()
	if config.MQTTServer.Enable {
		R.Lock()
		updateTopics()
		R.Unlock()

		//auto refresh config
		go func() {
			for {
				select {
				case <-readConfigDone:
					return

				case <-time.After(time.Duration(time.Second * time.Duration(config.Subscribe.AutoRefresh))):
					R.Lock()
					readConfigFile()
					updateTopics()
					R.Unlock()
				}
			}
		}()

		R.RLock()
		connOpts := MQTT.NewClientOptions().AddBroker(config.URI).SetClientID(config.ClientID)
		if config.Retain {
			//默认是清除上次的Session
			connOpts.SetCleanSession(!config.Retain)
		}
		if helper.IsExists(config.CAFile) && helper.IsExists(config.ClientCertFile) && helper.IsExists(config.ClientKeyFile) {
			connOpts.SetTLSConfig(NewTLSConfig(config.CAFile, config.ClientCertFile, config.ClientKeyFile))
		}

		logger.Info("Connecting MQTT Server", zap.String("URI", config.URI),
			zap.String("ClientId", config.ClientID),
		)

		//连接回调
		connOpts.OnConnect = func(c MQTT.Client) {
			if token := c.SubscribeMultiple(topics, onMsgReceived); token.Wait() && token.Error() != nil {
				logger.Panic("Subscribe is failed", zap.Error(token.Error()))
			}
		}

		//连接
		mqttClient = MQTT.NewClient(connOpts)
		//try
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			logger.Error("Connecting is failed.", zap.Error(token.Error()))
		} else {
			logger.Info("Connected is successfully.", zap.Skip())
		}
		R.RUnlock()
	}

	xClient = initRpcClient()

	<-chanSignal

	if config.MQTTServer.Enable {
		readConfigDone <- true
		if mqttClient.IsConnected() {
			mqttClient.Disconnect(300)
			logger.Warn("Client closed")
		}
	}

	<-time.After(time.Millisecond * 500)
}
