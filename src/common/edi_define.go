package common

import (
	"github.com/mitchellh/mapstructure"
)

type EDIRequest struct {
	Method string      `json:"method"`
	Params string  `json:"params"` //json 格式
}

func (s *EDIRequest) Unmarshal(in interface{}) error {
	if err := mapstructure.Decode(in, s); err != nil {
		return err
	}
	return nil
}
