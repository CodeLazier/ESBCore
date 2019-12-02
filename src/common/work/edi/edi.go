package edi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"common/fundef"
	"common/helper"
	pb "common/work/edi/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"
)

type _configs struct {
	A  struct{
		B struct{
			C struct {
				D  struct {
					Server string `yaml:"server"`
					Cert string `yaml:"cert"`
				} `yaml:"gRpc"`
			} `yaml:"EDI"`
		} `yaml:"Caller"`
	} `yaml:"General"`
}


var configs  struct{
	GRPC_Server string
	GRPC_CertFile string
}


func LoadCfg(content []byte) error{
	c:=&_configs{}
	if err:=yaml.Unmarshal(content,c);err!=nil{
		return err
	}
	configs.GRPC_CertFile=c.A.B.C.D.Cert
	configs.GRPC_Server=c.A.B.C.D.Server
	return nil
}

func init() {
	//由配置文件动态读取
	p := &EDICall{}
	for _, v := range []string{
		"NT/EDI/",
	} {
		fundef.RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type EDICall struct {
	Params []byte
}

//此Init非Package Init
func (self *EDICall) Init() {
	*self = EDICall{}
}

func (self *EDICall) Parse(params []byte) error {
	self.Params = params
	return nil
}

func (self *EDICall) Do(ctx context.Context, method string) (interface{}, error) {
	var op grpc.DialOption
	if helper.IsExists(configs.GRPC_CertFile) {
		cred, err := credentials.NewClientTLSFromFile("key/server.crt", "")
		if err != nil {
			return nil, err
		}
		op = grpc.WithTransportCredentials(cred)
	} else {
		op = grpc.WithInsecure()
	}

	//pool?
	//TODO port Configurable
	conn, err := grpc.Dial(fmt.Sprintf("%s",configs.GRPC_Server), op,grpc.WithDisableRetry())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewEDICallClient(conn)

	_method:=method
	if _method=="" {
		_m :=struct{
			Method string `json:"method"`
		}{""}
		if err:=json.Unmarshal(self.Params,&_m);err!=nil || _m.Method==""{
			return nil,errors.New("Method field is not found")
		}
		_method=_m.Method;
	}

	req := &pb.EDIRequest{
		Method: _method,
		Params: string(self.Params),
	}

	res, err := client.Call(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
