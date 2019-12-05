package common

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

const (
	CommonService         = "NT.ESB.CommonService"
	ESBRequestFunction    = "RequestFunc"
	JWTSecretKey          = "30495887dfgkjhbhj&*(*()@#$*&xmh/.d,/,."
	RPCEtcdRegisteredPath = "/NT_ESB_Core"
)

type Work interface {
	//序列化参数信息,注意针对同一函数,序列化前保留上次调用参数值,序列化未全部覆盖的话遗留的参数旧值会带入
	//所以这里使用Init来恢复默认状态(重新初始化)
	Init()
	Parse(interface{}) error
	Do(context.Context) (interface{}, error)
}

var (
	RegisterWorkMap       = make(map[string]Work)
	NoImplFunError  error = errors.New("Function Not implemented!")
)

type ESBRequest struct {
	ID        int64       `json:"-"` //snowflake id
	Topic     string      `json:"-"` //所属Topic
	TimeStamp time.Time   `json:"-"`
	Token     string      `json:"token,omitempty"`   //jwt token
	Extends   interface{} `json:"extends,omitempty"` //预留扩展
	Payload    interface{}      `json:"payload,omitempty"`  // 请求参数,用json表示
	Tag       string      `json:"tag,omitempty"`     //自定key
	Enqueue   bool        `json:"enqueue,omitempty"` //是否入队列
}

type ESBResponse struct {
	ID        int64       `json:"id"` //request的ID
	Topic     string      `json:"topic"`
	Tag       string      `json:"tag"`
	TimeStamp time.Time   `json:"timeStamp"`
	err       error       `json:"-"`                //逻辑错误,非技术性調用錯誤
	ErrNo     int         `json:"errNo"`            //编号
	ErrMsg    string      `json:"errMsg,omitempty"` //or序列化不进去?
	Result    interface{} `json:"result"`           //结果

}

//清空当前内容
func (req *ESBRequest) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	return nil
}

func (req *ESBRequest) Marshal() ([]byte, error) {
	if b, err := json.Marshal(req); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (res *ESBResponse) Marshal() ([]byte, error) {
	if res.err != nil && res.ErrMsg == "" {
		res.ErrMsg = res.err.Error()
	}
	if b, err := json.Marshal(res); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

//会清空内容
func (res *ESBResponse) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	return nil
}

func (res *ESBResponse) IsError() error {
	return res.err
}

func (res *ESBResponse) SetError(req *ESBRequest, err error) {
	res.TimeStamp = time.Now()
	res.ID = req.ID
	res.Topic = req.Topic
	res.err = err
	res.ErrNo = -1
	res.Tag = req.Tag
	res.Result = nil
	if err != nil {
		res.ErrMsg = err.Error()
	}
}

func (res *ESBResponse) AssignForRes(r *ESBResponse, result interface{}) *ESBResponse {
	res.TimeStamp = time.Now()
	res.ID = r.ID
	res.Topic = r.Topic
	res.Tag = r.Tag
	res.err = nil
	res.ErrMsg = ""
	res.ErrNo = 0
	res.Result = result
	return res
}

func (res *ESBResponse) AssignForReq(r *ESBRequest, result interface{}) *ESBResponse {
	res.TimeStamp = time.Now()
	res.ID = r.ID
	res.Topic = r.Topic
	res.Tag = r.Tag
	res.err = nil
	res.ErrMsg = ""
	res.ErrNo = 0
	res.Result = result
	return res
}

type ESBFunc_Impl func(context.Context, string) (interface{}, error)

//RPC返回值必须是指针类型
type ESBRequestFunc func(context.Context, *ESBRequest, *ESBResponse) error
