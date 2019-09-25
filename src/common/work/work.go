package work

import (
	"context"
	"strings"

	"common/fundef"
	//========================
	//register
	_ "common/work/calc"
	_ "common/work/call"
	_ "common/work/edi"
	_ "common/work/sample"
	_ "common/work/test"
)

type CmdQueue struct {
	Broker          string `yaml:"broker"`
	Backend         string `yaml:"backend"`
	ResultsExpireIn int    `yaml:"resultsExpireIn,omitempty"`
}

//TODO Json Marshal/Unmarshal对性能损害较大(reflace),因改用自解析或其他第三方解析库(msgPack?)
func MainEnter(ctx context.Context, reqStr string,routing []string) (string, error) {
	//unmarshal
	req, err := fundef.UnmarshalForRequestParams([]byte(reqStr))
	if err != nil {
		return "", err
	}

	//impl
	if work, ok := fundef.RegisterWorkMap[req.Topic]; ok {

		//init,防止状态带入
		work.Init()
		//parse
		err := work.Parse([]byte(req.Params))
		if err != nil {
			return "", err
		}

		i:=strings.LastIndexByte(req.Topic,'/')
		var ts string
		if i!=-1{
			ts=req.Topic[i+1:]
		}
		//ts := strings.SplitAfter(req.Topic, "/")
		//do
		result, err := work.Do(ctx, ts /*ts[len(ts)-1]*/)
		if err != nil {
			return "", err
		}

		//marshal
		res := &fundef.ResponseResult{}
		resJson, err := res.SetResult(req.ID, req.Topic, result).Marshal()
		if err != nil {
			return "", err
		}
		return string(resJson), nil
	}
	return "", fundef.NoImplFunError
}
