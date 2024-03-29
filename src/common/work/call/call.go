package call

import (
	"context"

	. "common"
	"common/helper"
	pb "common/work/call/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	CertFile string
)

func init() {
	//由配置文件动态读取
	//测试
	p := &Call{}
	for _, v := range []string{
		"Ex/Example/Call/Test",
	} {
		RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type Call struct {
	Params interface{}
}

func (self *Call) Init() {
	*self = Call{}
}

func (self *Call) Parse(params interface{}) error {
	self.Params = params
	return nil
}

func (self *Call) Do(ctx context.Context) (interface{}, error) {
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
	conn, err := grpc.Dial(":10036", op)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewCallClient(conn)

	req := &pb.Request{
		Method: "",
		Params: self.Params.(string),
	}

	res, err := client.Call(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
