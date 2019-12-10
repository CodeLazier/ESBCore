package edi

import (
	"context"
	"fmt"
	"sync"

	. "common"
	"common/helper"
	pb "common/work/edi/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"
)

var (
	GlobalGrpcOptions grpc.DialOption
	ConnPool sync.Map
)
type _configs struct {
	A struct {
		B struct {
			C struct {
				D struct {
					Server string `yaml:"server"`
					Cert   string `yaml:"cert"`
				} `yaml:"gRpc"`
			} `yaml:"EDI"`
		} `yaml:"Caller"`
	} `yaml:"General"`
}

var configs struct {
	GRPC_Server   string
	GRPC_CertFile string
}

func LoadCfg(content []byte) error {
	c := &_configs{}
	if err := yaml.Unmarshal(content, c); err != nil {
		return err
	}
	configs.GRPC_CertFile = c.A.B.C.D.Cert
	configs.GRPC_Server = c.A.B.C.D.Server
	return nil
}

func init() {
	//由配置文件动态读取
	p := &EDICall{}
	for _, v := range []string{
		"NT/EDI/",
	} {
		RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}

	if helper.IsExists(configs.GRPC_CertFile) {
		cred, err := credentials.NewClientTLSFromFile("key/server.crt", "")
		if err != nil {
			fmt.Println(err)
		}else {
			GlobalGrpcOptions = grpc.WithTransportCredentials(cred)
		}
	} else {
		GlobalGrpcOptions = grpc.WithInsecure()
	}

	//ConnPool=make(map[string]*grpc.ClientConn)
}

type EDICall struct {
	Params EDIRequest
}

func (self *EDICall) GetConn(addr string)(*grpc.ClientConn,error){
	if c,ok:=ConnPool.Load(addr);ok{
		return c.(*grpc.ClientConn),nil
	}
	if c,err:= grpc.Dial(addr, GlobalGrpcOptions, grpc.WithDisableRetry());err!=nil{
		return nil,err
	}else{
		ConnPool.Store(addr,c)
		return c,nil
	}
}

//此Init非Package Init
func (self *EDICall) Init() {
	*self = EDICall{}
}

func (self *EDICall) Parse(params interface{}) error {
	return self.Params.Unmarshal(params)
}

func (self *EDICall) Do(ctx context.Context) (interface{}, error) {
	//pool?
	//TODO port Configurable
	conn, err := self.GetConn(fmt.Sprintf("%s", configs.GRPC_Server))
	if err != nil {
		return nil, err
	}

	client := pb.NewEDICallClient(conn)

	req := &pb.EDIRequest{
		Method: self.Params.Method,
		Params: self.Params.Params,
	}

	res, err := client.Call(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
