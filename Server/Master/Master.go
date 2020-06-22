package main

import (
	"FinalProject/lock"
	clientMasterPb "FinalProject/proto/ClientMaster"
	"FinalProject/utils"

	"context"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
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

func (master *Master) WatchNewDataNode(path string) error{
	err := master.dataNodeManager.WatchNewDataNode(path)
	if err != nil{
		log.Printf("Master: watch new node error: %v\n", err)
		return err
	}
	return nil
}


func InitializeRootZnodes(conn *zk.Conn) error{
	rootZnodes := []string{backupNodesPath, dataNodesPath, lock.ReaderNumRootPath, lock.LockPath, lock.GlobalLockPath,
		fmt.Sprintf("%s%s", lock.ReaderNumRootPath, lock.GlobalReaderPath),
		fmt.Sprintf("%s%s", lock.LockPath, lock.ReaderLockPath),
		fmt.Sprintf("%s%s", lock.LockPath, lock.WriterLockPath),
	}
	for _, znodePath := range rootZnodes{
		exist, _, err := conn.Exists(znodePath)
		if err != nil{
			utils.Error("Check existence in InitializeRootZnodes error: %v\n", err)
			return err
		}
		if !exist{
			_, err := conn.Create(znodePath, []byte("0"), 0, zk.WorldACL(zk.PermAll))
			if err != nil{
				utils.Error("Create new root znodes error: %v\n", err)
				return err
			}
		}
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
	if err := InitializeRootZnodes(zkConn); err != nil{
		log.Fatal(err)
	}

	dataNodeManager, err := NewDataNodeManager(zkConn)
	if err != nil{
		log.Fatalf("fail to initialize data node manager: %v", err)
	}
	masterServer := Master{
		dataNodeManager: dataNodeManager,
	}

	go masterServer.WatchNewDataNode(dataNodesPath)

	s := grpc.NewServer()
	clientMasterPb.RegisterClientMasterServer(s, &masterServer)

	fmt.Println("Register complete, ready to serve...")
	if err := s.Serve(lis); err != nil{
		log.Fatalf("err to serve: %v", err)
	}
}

