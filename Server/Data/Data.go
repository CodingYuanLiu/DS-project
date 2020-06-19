package main

import (
	"FinalProject/lock"
	clientDataPb "FinalProject/proto/ClientData"
	dataDataPb "FinalProject/proto/DataData"
	masterDataPb "FinalProject/proto/MasterData"
	"FinalProject/utils"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)
const(
	defaultDataPort = ":7777"
	dataNodesPath   = "/DataNode" //The root node of the data nodes' ports
	aliveResp       = "Alive"
	aliveReq        = "Is alive?"
	masterPort      = ":7000"
	defaultIP       = "localhost"
)

func ConnectZookeeper() *zk.Conn{
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Zookeeper connected")
	return zkConn
}

func ZkRegisterDataNodePort(conn *zk.Conn, port string) error{
	path := fmt.Sprintf("%v/%v", dataNodesPath, port)
	exist, _, err := conn.Exists(path)
	if err != nil{
		return err
 	} else if exist{
 		return errors.New("the data node already exists in the zookeeper")
	}

	_, err = conn.Create(path, []byte("Alive"), 0, zk.WorldACL(zk.PermAll))
	if err != nil{
		return err
	}

	//create node to record the number of readers
	readerPath := fmt.Sprintf("%v/%v", lock.ReaderNumRootPath, port)
	exist, s, err := conn.Exists(readerPath)
	if err != nil{
		return err
	} else if exist{
		_, err := conn.Set(readerPath, []byte("0"), s.Version)
		if err != nil{
			log.Printf("Set readers error")
			return err
		}
		return nil
	}
	_, err = conn.Create(readerPath, []byte("0"), 0, zk.WorldACL(zk.PermAll))

	return err
}

func HeartBeatResponse(conn *zk.Conn, port string) error{
	path := fmt.Sprintf("%v/%v", dataNodesPath, port)
	exist, _, err := conn.Exists(path)
	if !exist{
		return errors.New("the data node does not exist in the zookeeper")
	} else if err != nil{
		return err
	}

	for {
		_, s, getCh, err := conn.GetW(path)
		if err != nil {
			fmt.Printf("watch node error: %v\n", err)
			return err
		}

		select {
		case chEvent := <- getCh:
			{
				if chEvent.Type == zk.EventNodeDataChanged {
					log.Printf("heart beat on port %v suceed\n", port)
					_, err := conn.Set(path, []byte(aliveResp), s.Version + 1)
					if err != nil{
						log.Printf("Error to set heart beat detection response: %v\n", err)
					}
				}
			}
		}
	}
}

func InitializeDataServer(lis net.Listener, port string){
	grpcServer := grpc.NewServer()
	zkConn := ConnectZookeeper()

	err := ZkRegisterDataNodePort(zkConn, port)
	if err != nil{
		log.Fatalf("Register data node to zookeeper error: %v\n", err)
	}

	go HeartBeatResponse(zkConn, port)
	dataServer := &DataServer{
		database: map[string]string{},
		//dataMasterCli: NewDataMasterCli(),
		port: port,
	}

	//TEST
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		resp, _ := dataServer.dataMasterCli.DataMasterReshardComplete(ctx, &dataMasterPb.DataMasterReshardCompleteReq{
			Message: "reshard complete",
		})
		utils.Debug("Response from master via rpc: %v\n", resp.Message)
	*/
	//END TEST

	//Register rpc servers
	clientDataPb.RegisterClientDataServer(grpcServer, dataServer)
	masterDataPb.RegisterMasterDataServer(grpcServer, dataServer)
	dataDataPb.RegisterDataDataServer(grpcServer, dataServer)

	if err = grpcServer.Serve(lis); err != nil{
		log.Fatal(err)
	}
}

func main() {
	port := defaultDataPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	lis, err := net.Listen("tcp", port)
	fmt.Printf("Data node listening port %v...\n", port)

	if err != nil{
		//TODO: err so register backup server
		if err := InitializeBackupServer(port); err != nil{
			utils.Error("Initialize backup server error: %v\n", err)
			log.Fatal(err)
		}
		//Should never reach here
		log.Fatal(err)
	}

	InitializeDataServer(lis, port)
}
