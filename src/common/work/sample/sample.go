package sample

import (
	"context"

	. "common"
)

func init() {
	p := &Sample{}
	for _, v := range []string{
		"NT/Sample/",
	} {
		RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type Sample struct {
}

func (self *Sample) Init() {
	*self = Sample{}
}

func (*Sample) Parse(interface{}) error {
	return nil
}

func (*Sample) Do(ctx context.Context) (interface{}, error) {
	return "Pong", nil
}
