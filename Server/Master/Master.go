package main

import (
	clientMasterPb "FinalProject/proto/ClientMaster"
	"context"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)


// Interfaces: PUT, DELETE, READ
const(
	masterPort = ":7000"
	dataNodesPath = "/DataNode" //The root node of the data nodes' ports
)

type Master struct{
	clientMasterPb.UnimplementedClientMasterServer
	//Store the metadata of the nodes
	dataNodeManager *DataNodeManager //key:value => ID:port
}

func (master *Master) ClientMasterFindDataNode(ctx context.Context, req *clientMasterPb.ClientMasterFindDataNodeReq) (*clientMasterPb.ClientMasterFindDataNodeResp, error){
	log.Printf("Serve client...")

	dataPort, err := master.dataNodeManager.FindDataNode(req.Key)
	if err != nil{
		log.Printf("Find data node error: %v\n", err)
		return nil, err
	}
	return &clientMasterPb.ClientMasterFindDataNodeResp{
		Port: dataPort,
	}, nil
}

func ConnectZookeeper() *zk.Conn{
	fmt.Println("Connect zookeeper...")
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Zookeeper connected")
	return zkConn
}

func (master *Master) WatchNewNode(conn *zk.Conn, path string) error{
	err := master.dataNodeManager.WatchNewNode(conn, path)
	if err != nil{
		log.Printf("Master: watch new node error: %v\n", err)
		return err
	}
	return nil
}

func main(){
	fmt.Println("Start to run master node...")
	lis,err := net.Listen("tcp", masterPort)
	if err != nil{
		log.Fatalf("failed to listen %v\n", err)
	}

	zkConn := ConnectZookeeper()
	defer zkConn.Close()

	masterServer := Master{
		dataNodeManager: NewDataNodeManager(zkConn),
	}

	go masterServer.WatchNewNode(zkConn, dataNodesPath)

	s := grpc.NewServer()
	clientMasterPb.RegisterClientMasterServer(s, &masterServer)
	fmt.Println("Register complete, ready to serve...")
	if err := s.Serve(lis); err != nil{
		log.Fatalf("err to serve: %v", err)
	}
}

