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

func (master *Master) WatchNewNode(conn *zk.Conn, path string) error {
	exist, _, err := conn.Exists(path)
	if err != nil{
		fmt.Println(err)
	}
	if !exist{
		_, err = conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
		log.Printf("Register Root data node in zookeeper: %s, node is created\n", path)
	}

	for {
		_, _, getCh, err := conn.ChildrenW(path)
		if err != nil {
			fmt.Printf("watch children error: %v\n", err)
		}

		select {
			case chEvent := <- getCh:
			{
				fmt.Printf("%+v\n", chEvent)
				if chEvent.Type == zk.EventNodeChildrenChanged {
					fmt.Printf("detect data node changed on zookeeper\n")
					v,_, err := conn.Children(path)
					if err != nil{
						return err
					}
					fmt.Printf("value of path[%s]=[%s].\n", path, v)
					master.dataNodeManager.HandleDataNodesChanges(v)
				} else{
					fmt.Printf("other events on path %s\n", chEvent.Path)
				}
			}
		}
	}
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

