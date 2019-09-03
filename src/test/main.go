package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "common/work/call/pb"
	"google.golang.org/grpc"
)

var NoError error = errors.New("No provide!")

type CallServer struct {

}

func (self *CallServer) Call(ctx context.Context, res *pb.Request) (*pb.Response, error){
	return &pb.Response{
		Result : fmt.Sprintf(`{"msg":"my work is completed","err","","request":"%s"}`,res.Params),
	},nil
}

func (self *CallServer)ServerStreamingCall(*pb.Request, pb.Call_ServerStreamingCallServer) error{
	return NoError
}

func (self *CallServer)ClientStreamingCall(pb.Call_ClientStreamingCallServer) error{
	return NoError
}

func (self *CallServer)BidirectionalStreamingCall(pb.Call_BidirectionalStreamingCallServer) error{
	return NoError
}


func main(){
	lis, err := net.Listen("tcp", ":10036")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCallServer(grpcServer, &CallServer{})
	if grpcServer.Serve(lis)!=nil{
		fmt.Println("Error")
	}
}
