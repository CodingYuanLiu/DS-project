package main

import (
	pb "FinalProject/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)
const(
	port =":50051"
)

type server struct{
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error){
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{
		Message: "Hello " + in.Name,
	},nil
}

func (s *server) SayHelloTwice(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error){
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{
		Message: "Hello " + in.Name + " TWICE",
	}, nil
}

func main(){
	lis,err := net.Listen("tcp", port)
	if err != nil{
		log.Fatalf("failed to listen %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil{
		log.Fatalf("err to serve: %v", err)
	}
}