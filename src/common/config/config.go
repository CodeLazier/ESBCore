package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	REG_CALL = "/nt/esb/conf/work/regcall"
)

func NewEtcdClient() (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SetRegCall(client *clientv3.Client, reg []string) error {
	_, err := client.Put(context.Background(), REG_CALL, strings.Join(reg, ","))
	if err != nil {
		return err
	}
	return nil
}

func GetRegCall(client *clientv3.Client) ([]string, error) {
	resp, err := client.Get(context.Background(), REG_CALL)
	if err != nil {
		return nil, err
	}

	r := make([]string, 0)
	for _, kvs := range resp.Kvs {
		r = append(r, string(kvs.Value))
	}

	return r, nil
}

func WatchRegCall(client *clientv3.Client) {
	for {
		watch := client.Watch(context.Background(), REG_CALL)
		for res := range watch {
			for _, v := range res.Events {
				fmt.Printf("%s %q : %q \n", v.Type, v.Kv.Key, v.Kv.Value)
			}
		}
	}
}

/*
 client, err := etcd_client.New(etcd_client.Config{
        Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        fmt.Printf("connect failed ,err ", err)
        return
    }

    defer client.Close()

    background := context.Background()
    client.Put(background, "/logagent/conf/", "123456")
    if err != nil {
        fmt.Println("err :", err)
        return
    }

    fmt.Println("connec success !!")
    for {
        watch := client.Watch(context.Background(), "/logagent/conf/")
        for wresp := range watch {
            for _, v := range wresp.Events {
                fmt.Printf("%s %q : %q \n", v.Type,v.Kv.Key,v.Kv.Value)
            }
        }
    }
*/
