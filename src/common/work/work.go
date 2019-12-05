package work

import (
	"context"

	. "common"
	//========================
	//register
	_ "common/work/calc"
	_ "common/work/call"
	"common/work/edi"
	_ "common/work/sample"
	_ "common/work/test"
)

func LoadAllCfg(content []byte)error{
	err:=edi.LoadCfg(content)

	return err
}

type CmdQueue struct {
	Broker          string `yaml:"broker"`
	Backend         string `yaml:"backend"`
	ResultsExpireIn int    `yaml:"resultsExpireIn"`
	WaitTimeout int `yaml:"waitTimeout"`
}

func MainEnterDirect(ctx context.Context, req* ESBRequest,routing []string) (string, error) {

	//impl
	if work, ok := RegisterWorkMap[req.Topic]; ok {
		//init,防止状态带入
		work.Init()
		//parse
		err := work.Parse(req.Payload)
		if err != nil {
			return "", err
		}

		result, err := work.Do(ctx)
		if err != nil {
			return "", err
		}

		//marshal
		res := &ESBResponse{}
		resJson, err := res.AssignForReq(req, result).Marshal()
		if err != nil {
			return "", err
		}
		return string(resJson), nil
	}
	return "", NoImplFunError
}

//TODO Json Marshal/Unmarshal对性能损害较大(reflace),因改用自解析或其他第三方解析库(msgPack?)
func MainEnter(ctx context.Context, reqStr string,routing []string) (string, error) {
	//unmarshal
	req:=&ESBRequest{}
	if err := req.Unmarshal([]byte(reqStr));err!=nil{
		return "", err
	}
	return MainEnterDirect(ctx,req,routing)
}
