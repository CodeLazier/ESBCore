package test

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"

	"common/fundef"
)

func init() {
	p := &Test{}
	for _, v := range []string{
		"NT/Test/DataType",
		"NT/Test/Ping",
		"NT/Test/FetchImage",
	} {
		fundef.RegisterWorkMap[v] = p //同一指针地址,节省内存,提高效率
	}
}

type SubItem struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

type Test struct {
	Num      int               `json:"num"`
	Float    float32           `json:"float"`
	Float64  float64           `json:"float64"`
	Bool     bool              `json:"bool"`
	Str      string            `json:"string"`
	Array    []int             `json:"array"`
	ArrayStr []string          `json:"arrayStr"`
	Map      map[int]string    `json:"map"`
	KV       map[string]string `json:"kv"`
	Items    []SubItem         `json:"items"`
}

func (self *Test) Init() {
	*self = Test{}
}

func (*Test) Parse([]byte) error {
	return nil
}

func (self *Test) fetchImage() (interface{}, error) {
	resp, err := http.Get("http://webeip.newtype.com.tw/object/EIP_logo/newtype_logow.png")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (self *Test) dataType() (interface{}, error) {
	self.Num = rand.Int()
	self.Float = rand.Float32()
	self.Float64 = rand.Float64()
	self.Bool = rand.Intn(2) == 0
	self.Str = "is joke"
	self.Array = make([]int, 0)
	for i := 0; i < 10; i++ {
		self.Array = append(self.Array, rand.Int())
	}
	self.ArrayStr = []string{
		"My fancies are fireflies",
		"specks of living light",
		"twinkling in the dark",
	}
	self.Map = map[int]string{
		0:   "tom",
		20:  "jerry",
		50:  "duck",
		993: "mouse",
	}
	self.KV = map[string]string{
		"tom":  "30",
		"duck": "16",
		"man":  "supper",
	}
	self.Items = []SubItem{
		{
			Name: "xiao ming",
			Age:  23,
			Job:  "worker",
		},
		{
			Name: "zhang liang",
			Age:  19,
			Job:  "student",
		},
		{
			Name: "ding ding",
			Age:  57,
			Job:  "CEO",
		},
	}
	return self, nil
}

func (*Test) ping() (interface{}, error) {
	return "Pong", nil
}

func (self *Test) Do(ctx context.Context, method string) (interface{}, error) {
	switch method {
	case "Ping":
		return self.ping()
	case "DataType":
		return self.dataType()
	case "FetchImage":
		return self.fetchImage()
	}
	return nil, fundef.NoImplFunError
}
