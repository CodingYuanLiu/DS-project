package main

import(
	pb "FinalProject/proto/MasterData"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)
const(
	port1 = ":7777"
	someValue = "The world focus on you: data1"
)

type DataServer struct{
	pb.UnimplementedMasterDataServer
}

func (dataServer *DataServer) MasterDataPut(ctx context.Context, req *pb.MasterDataPutReq) (*pb.MasterDataPutResp, error){
	return &pb.MasterDataPutResp{
		Message: someValue,
	}, nil
}

func main() {
	port := port1
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	lis, err := net.Listen("tcp", port)
	fmt.Printf("Data node listening port %v...\n", port)

	if err != nil{
		log.Fatal(err)
	}
	dataServer := grpc.NewServer()
	pb.RegisterMasterDataServer(dataServer, &DataServer{})
	if err = dataServer.Serve(lis); err != nil{
		log.Fatal(err)
	}
}
