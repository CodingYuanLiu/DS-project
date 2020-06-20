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

func testByteStringConvert(){
	str := []string{"hello","world","key2","key3"}
	fmt.Printf("Original str:\n%v\n\n", str)
	bytes := utils.StringArrayToByte(str)
	fmt.Printf("Converted bytes:\n%v\n\n", bytes)
	fmt.Printf("Converted str:\n%v\n\n", utils.ByteToStringArray(bytes))
}


func testMapByteConvert(){
	database := map[string]string{
		"key1": ":7777",
		"key2": ":2200",
	}
	fmt.Printf("Original map:\n%v\n\n", database)
	bytes := utils.KeyValueMapToByte(database)
	fmt.Printf("Converted bytes:\n%v\n\n", bytes)
	fmt.Printf("Converted map:\n%v\n\n", utils.ByteToKeyValueMap(bytes))
}

func continuousPrint(){
	for{
		fmt.Println("print...")
		time.Sleep(time.Second)
	}
}
func testNestedGoroutine(){
	go continuousPrint()
	fmt.Println("Function return")
}

func testZKExist(){
	conn := ConnectZookeeper()
	exist, _, err := conn.Exists("/path1231/zxc")
	fmt.Println(exist)
	fmt.Println(err)
}
func main(){
	m := map[int] int{
		2: 4,
		3:9,
	}
	fmt.Println(len(m))
}
