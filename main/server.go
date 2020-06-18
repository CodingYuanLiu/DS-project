package main

import (
	clientMasterPb "FinalProject/proto/ClientMaster"
	masterDataPb "FinalProject/proto/MasterData"
	"FinalProject/utils"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	lockPath = "/testlock"
	path1 = "testlock/path1"
	masterPort = ":7000"
)

type Master struct{
	clientMasterPb.UnimplementedClientMasterServer
	masterDataPb.UnimplementedMasterDataServer
}

func testUtilDebug(){
	utils.Debug("hello %d, %d", 1, 2)
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

func test2RpcServers(){
	lis,err := net.Listen("tcp", masterPort)
	if err != nil{
		log.Fatalf("failed to listen %v\n", err)
	}

	zkConn := ConnectZookeeper()
	defer zkConn.Close()

	masterServer := Master{
	}


	s := grpc.NewServer()
	clientMasterPb.RegisterClientMasterServer(s, &masterServer)
	masterDataPb.RegisterMasterDataServer(s, &masterServer)
	fmt.Println("Register complete, ready to serve...")
	if err := s.Serve(lis); err != nil{
		log.Fatalf("err to serve: %v", err)
	}
}

func main(){
	test2RpcServers()
}
