package fundef

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Work interface {
	//序列化参数信息,注意针对同一函数,序列化前保留上次调用参数值,序列化未全部覆盖的话遗留的参数旧值会带入
	//所以这里使用Init来恢复默认状态(重新初始化)
	Init()
	Parse([]byte) error
	Do(context.Context, string) (interface{}, error)
}

var (
	RegisterWorkMap       = make(map[string]Work)
	NoImplFunError  error = errors.New("Function Not implemented!")
)

type RequestParams struct {
	ID        int64  //snowflake id
	Topic     string //所属Topic
	TimeStamp time.Time
	Token string //jwt token
	Extends interface{} //预留扩展
	Params string // 请求参数,用json表示
}

type ResponseResult struct {
	ID        int64 //request的ID
	Topic     string
	TimeStamp time.Time
	Err       error //逻辑错误,非技术性調用錯誤

	Result interface{} //结果

}

func UnmarshalForRequestParams(data []byte) (*RequestParams, error) {
	req := RequestParams{}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	} else {
		return &req, nil
	}
}

func UnmarshalForResponseResult(data string) (*ResponseResult, error) {
	res := ResponseResult{}
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		return nil, err
	} else {
		return &res, nil
	}
}

func (req *RequestParams) Marshal() ([]byte, error) {
	if b, err := json.Marshal(req); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (res *ResponseResult) Marshal() ([]byte, error) {
	if b, err := json.Marshal(res); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (res *ResponseResult) SetError(id int64, topic string, err error) {
	res.TimeStamp = time.Now()
	res.ID = id
	res.Topic = topic
	res.Err = err
}

func (res *ResponseResult) SetResult(id int64, topic string, result interface{}) *ResponseResult {
	res.TimeStamp = time.Now()
	res.ID = id
	res.Topic = topic
	res.Result = result
	return res
}

type Func_Impl func(context.Context, string) (interface{}, error)

//RPC返回值必须是指针类型
type RequestFunc func(context.Context, *RequestParams, *ResponseResult) error
