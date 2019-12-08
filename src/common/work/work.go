package work

import (
	"context"
	"errors"

	. "common"
	//========================
	//register
	_ "common/work/calc"
	_ "common/work/call"
	"common/work/edi"
	_ "common/work/sample"
	_ "common/work/test"
)

func LoadAllCfg(content []byte) error {
	err := edi.LoadCfg(content)

	return err
}

type CmdQueue struct {
	Broker          string `yaml:"broker"`
	Backend         string `yaml:"backend"`
	ResultsExpireIn int    `yaml:"resultsExpireIn"`
	WaitTimeout     int    `yaml:"waitTimeout"`
}

//return resJson,errStr
func MainEnter(req *ESBRequest) (string,string) {
	//impl
	var result interface{}
	var err error
	res := &ESBResponse{}
	if work, ok := RegisterWorkMap[req.Topic]; ok {
		//init,防止状态带入
		work.Init()
		//parse
		if err = work.Parse(req.Payload); err == nil {
			if result, err = work.Do(context.Background()); err == nil {
				//ok
			}
		}

	} else {
		err = NoImplFunError
	}

	if err!=nil{
		return "",err.Error()
	}

	//marshal
	resJson, err := res.AssignForReq(req, result).Marshal()
	if err != nil {
		return "",err.Error()
	}

	return string(resJson),""
}

func MainEnterDirect(ctx context.Context, req *ESBRequest) (*ESBResponse, error) {
	resJson, errStr:= MainEnter(req)
	if errStr!=""{
		return nil,errors.New(errStr)
	}

	res := &ESBResponse{}
	if err := res.Unmarshal([]byte(resJson)); err != nil {
		return nil, err
	}

	return res,nil
}

//TODO Json Marshal/Unmarshal对性能损害较大(reflace),因改用自解析或其他第三方解析库(msgPack?)
//func MainEnter2(ctx context.Context, reqStr string, routing []string) (string, error) {
//	//unmarshal
//	req := &ESBRequest{}
//	if err := req.Unmarshal([]byte(reqStr)); err != nil {
//		return "", err
//	}
//	return reqStr, nil
//	return MainEnterDirect(ctx, req, routing)
//}
//
//func MainEnter(req *ESBRequest) string {
//	//unmarshal
//	//req:=&ESBRequest{}
//	//if err := req.Unmarshal([]byte(reqStr));err!=nil{
//	//	return ""
//	//}
//	s, _ := req.Marshal()
//	return string(s)
//
//}
