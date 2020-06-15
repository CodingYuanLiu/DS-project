package main

import (
	clientMasterPb "FinalProject/proto/ClientMaster"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)


// Interfaces: PUT, DELETE, READ
const(
	masterPort = ":7000"
)

type Master struct{
	clientMasterPb.UnimplementedClientMasterServer
	//Store the metadata of the nodes
	dataNodeManager *DataNodeManager //key:value => ID:port
}

func (master *Master) ClientMasterPut(ctx context.Context, req *clientMasterPb.ClientMasterPutReq) (*clientMasterPb.ClientMasterPutResp, error){
	log.Printf("Serve client...")

	dataPort, err := master.dataNodeManager.FindDataNode(req.Key)
	if err != nil{
		log.Fatalf("Find data node error: %v\n", err)
	}
	return &clientMasterPb.ClientMasterPutResp{
		Port: dataPort,
	}, nil
}

func main(){
	fmt.Println("Start to run master node...")
	lis,err := net.Listen("tcp", masterPort)
	if err != nil{
		log.Fatalf("failed to listen %v\n", err)
	}
	s := grpc.NewServer()
	masterServer := Master{
		dataNodeManager: NewDataNodeManager(),
	}

	_ = masterServer.dataNodeManager.RegisterDataNode(dataPort1)
	_ = masterServer.dataNodeManager.RegisterDataNode(dataPort2)

	clientMasterPb.RegisterClientMasterServer(s, &masterServer)
	fmt.Println("Register complete, ready to serve...")
	if err := s.Serve(lis); err != nil{
		log.Fatalf("err to serve: %v", err)
	}
}

