package calc

import (
	"context"
	"encoding/json"
	"errors"

	"common/fundef"
)

func init() {
	p := &CalcParams{}
	for _, v := range []string{
		"NT/Sample/Calc/Add",
		"NT/Sample/Calc/Sub",
		"NT/Sample/Calc/Mul",
		"NT/Sample/Calc/Div",
	} {
		fundef.RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type CalcParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (self *CalcParams) Init() {
	*self = CalcParams{}
}

func (self *CalcParams) Parse(jsonData []byte) error {
	if err := json.Unmarshal(jsonData, self); err != nil {
		return err
	}
	return nil
}

func (self *CalcParams) Do(ctx context.Context, method string) (interface{}, error) {
	switch method {
	case "Add":
		return self.A + self.B, nil
	case "Sub":
		return self.A - self.B, nil
	case "Mul":
		return self.A * self.B, nil
	case "Div":
		return self.A / self.B, nil
	}
	return nil, errors.New("No support is operator!")
}
