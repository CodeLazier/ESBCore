package sample

import (
	"context"

	"common/fundef"
)

func init() {
	p := &Sample{}
	for _, v := range []string{
		"NT/Sample/",
	} {
		fundef.RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type Sample struct {
}

func (self *Sample) Init() {
	*self = Sample{}
}

func (*Sample) Parse([]byte) error {
	return nil
}

func (*Sample) Do(ctx context.Context, method string) (interface{}, error) {
	return "Pong", nil
}
