package server

import (
	"common/task/controller"
	"common/task/message"
	"common/task/yerrors"
	"time"
)

type Client struct {
	isClone bool
	server  *Server
	ctl     controller.TaskCtl

	// ctl name
	RetryCount string
}

func NewClient(server *Server) Client {
	return Client{
		server:     server,
		ctl:        controller.NewTaskCtl(),
		RetryCount: "RetryCount",
	}
}
func (c *Client) Clone() *Client {
	if c.isClone {
		return c
	} else {
		return &Client{
			isClone:    true,
			server:     c.server,
			ctl:        c.ctl,
			RetryCount: c.RetryCount,
		}
	}
}
func (c *Client) SetTaskCtl(name string, value interface{}) *Client {
	cloneC := c.Clone()
	switch name {
	case cloneC.RetryCount:
		cloneC.ctl.RetryCount = value.(int)
	}
	return cloneC
}

// return: taskId, err
func (c *Client) Send(groupName string, workerName string, args ...interface{}) (string, error) {
	return c.server.Send(groupName, workerName, c.ctl, args...)
}

// taskId:
// timeout:
// sleepDuration:
func (c *Client) GetResult(taskId string, timeout time.Duration, sleepTime time.Duration) (message.Result, error) {
	if c.server.backend == nil {
		return message.Result{}, yerrors.ErrNilResult{}
	}
	n := time.Now()
	for {
		if time.Now().Sub(n) >= timeout {
			return message.Result{}, yerrors.ErrTimeOut{}
		}
		r, err := c.server.GetResult(taskId)
		if err == nil && r.IsFinish() {
			return r, nil
		}
		time.Sleep(sleepTime)
	}
}
