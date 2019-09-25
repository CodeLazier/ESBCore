package edi

import (
	"context"

	"common/fundef"
	"common/helper"
	pb "common/work/edi/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	CertFile string
)

func init() {
	//由配置文件动态读取
	p := &EDICall{}
	for _, v := range []string{
		"NT/EDI/runGroup",
	} {
		fundef.RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type EDICall struct {
	Params []byte
}

func (self *EDICall) Init() {
	*self = EDICall{}
}

func (self *EDICall) Parse(params []byte) error {
	self.Params = params
	return nil
}

func (self *EDICall) Do(ctx context.Context, method string) (interface{}, error) {
	var op grpc.DialOption
	if helper.IsExists(CertFile) {
		creds, err := credentials.NewClientTLSFromFile("key/server.crt", "")
		if err != nil {
			return nil, err
		}
		op = grpc.WithTransportCredentials(creds)
	} else {
		op = grpc.WithInsecure()
	}

	//pool?
	//TODO port Configurable
	conn, err := grpc.Dial(":10038", op)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewEDICallClient(conn)

	req := &pb.EDIRequest{
		Method: method,
		Params: string(self.Params),
	}

	res, err := client.Call(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
