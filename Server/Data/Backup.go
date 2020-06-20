package main

import (
	dataDataPb "FinalProject/proto/DataData"
	"FinalProject/utils"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
)

type BackupServer struct{
	dataDataPb.UnimplementedDataDataServer
	port string
	database map[string] string
	backupNodes map[string] string
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func BackupHeartBeatResponse(conn *zk.Conn, dataPort string, backupPort string) error{
	path := fmt.Sprintf("%s/%v/%v", backupNodesPath, dataPort, backupPort)
	exist, _, err := conn.Exists(path)
	if !exist{
		return errors.New("the data node does not exist in the zookeeper")
	} else if err != nil{
		return err
	}

	if err := HeartBeatResponseLoop(conn, path); err != nil{
		if err == ErrMessagePromote{
			log.Printf("Try to promote to data server %s\n", dataPort)
			//TODO: Promote from a backup server to a data server
		} else{
			utils.Error("HeartBeatResponse error: %v\n", err)
			return err
		}
	}
	return nil
}

//The function will normally blocked, as a backup server or finally serve as a data server
func InitializeBackupServer(dataPort string, dataServer *DataServer) error{
	backupPortInt, err := GetFreePort()
	if err != nil{
		utils.Error("Get free port error: %v\n", err)
		return err
	}
	backupPort := fmt.Sprintf(":%d", backupPortInt)
	log.Printf("Initialize backup server of %s on port %s...\n", dataPort, backupPort)

	//TODO: It keeps synchronization with the data server in the dead loop
	backupServer := &BackupServer{
		port: dataPort,
		database: map[string] string{},
		backupNodes: map[string] string{},
	}

	go RegisterSyncServer(backupPort, backupServer)

	//Create the backup node on the zookeeper
	zkConn := ConnectZookeeper()
	err = ZkRegisterBackupNodePort(zkConn, dataPort, backupPort)
	if err != nil{
		utils.Error("ZkRegisterBackupNodePort error: %v\n", err)
	}

	//Watch the znode in a dead loop, response for heart beat detection, until master set the special flag on the znode
	go BackupHeartBeatResponse(zkConn, dataPort, backupPort)


	//TODO: Initialize as a data server if it breaks from the dead loop.
	//TODO: Don't forget the delete the corresponding znode from the zookeeper
	//TODO: Serve as a data server and blocked

	for{
		;
	}
	return nil
}

//Create the zookeeper znode of the backup node
func ZkRegisterBackupNodePort(zkConn *zk.Conn, dataPort string, backupPort string) error {
	backupRootPortPath := fmt.Sprintf("%s/%s", backupNodesPath, dataPort)
	//Check the root path of the data node's backup node
	exist, _, err := zkConn.Exists(backupRootPortPath)
	if err != nil{
		utils.Error("Check root path error: %v\n", err)
		return err
	} else if !exist{
		utils.Error("The backup root path does not exist: %v\n", err)
		return err
	}

	backupPath := fmt.Sprintf("%s/%s", backupRootPortPath, backupPort)
	_, err = zkConn.Create(backupPath, []byte(aliveResp), 0, zk.WorldACL(zk.PermAll))
	if err != nil{
		utils.Error("Create backup znode error:%v\n", err)
		return err
	}
	return nil
}

func RegisterSyncServer(backupPort string, backupServer *BackupServer){

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", backupPort)
	if err != nil{
		utils.Error("fail to listen the backup port: %s\n", backupPort)
		log.Fatal(err)
	}

	dataDataPb.RegisterDataDataServer(grpcServer, backupServer)
	if err = grpcServer.Serve(lis); err != nil{
		log.Fatal(err)
	}
}
