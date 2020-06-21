package main

import (
	clientDataPb "FinalProject/proto/ClientData"
	dataDataPb "FinalProject/proto/DataData"
	masterDataPb "FinalProject/proto/MasterData"
	"FinalProject/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type DataServer struct{
	clientDataPb.UnimplementedClientDataServer
	masterDataPb.UnimplementedMasterDataServer
	dataDataPb.UnimplementedDataDataServer
	port string
	database map[string] string
	backupNodes map [string] string //Use as a set
	//dataMasterCli dataMasterPb.DataMasterClient
}

func (dataServer *DataServer) ClientDataPut(ctx context.Context, req *clientDataPb.ClientDataPutReq) (*clientDataPb.ClientDataPutResp, error){

	dataServer.database[req.Key] = req.Value

	log.Printf("put key: %v, value: %v\n", req.Key, req.Value)
	log.Println(dataServer.database)

	err := SyncBackupPut(dataServer.backupNodes, req.Key, req.Value)
	if err != nil{
		utils.Error("SyncBackupPut in ClientDataPut error: %v\n", err)
		return nil, err
	}
	return &clientDataPb.ClientDataPutResp{
		Message: "[Data server]: put succeed",
	}, nil
}

func SyncBackupPut(backupNodes map[string] string, key string, value string) error{
	for backupNode, _ := range backupNodes{
		cli := NewDataDataCli(backupNode)
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		_, err := cli.DataDataSyncPut(ctx, &dataDataPb.DataDataSyncPutReq{
			Key: key,
			Value: value,
		})
		if err != nil{
			utils.Error("Sync backup node %s in SyncBackupPut error: %v\n", backupNode, err)
			return err
		}
		cancel()
	}
	return nil
}

func (dataServer *DataServer) ClientDataRead(ctx context.Context, req *clientDataPb.ClientDataReadReq) (*clientDataPb.ClientDataReadResp, error){
	value, exist := dataServer.database[req.Key]
	if !exist{
		return nil, errors.New("no value in the database")
	}
	utils.Debug("read key: %v, value: %v\n", req.Key, value)
	utils.Debug("database in ClientDataRead: %v\n", dataServer.database)
	//No need to sync to backup nodes
	return &clientDataPb.ClientDataReadResp{
		Value: value,
		Message: "[Data server]: read succeed",
	}, nil
}

func (dataServer *DataServer) ClientDataDelete(ctx context.Context, req *clientDataPb.ClientDataDeleteReq) (*clientDataPb.ClientDataDeleteResp, error){
	_, exist := dataServer.database[req.Key]
	if !exist{

		return nil, errors.New("no value in the database")
	}
	delete(dataServer.database, req.Key)

	err := SyncBackupDelete(dataServer.backupNodes, req.Key)
	if err != nil{
		utils.Error("SyncBackupDelete in ClientDataDelete error: %v\n", err)
	}
	return &clientDataPb.ClientDataDeleteResp{
		Message: "[Data server]: delete succeed",
	}, nil

}

func SyncBackupDelete(backupNodes map[string] string, key string) error{
	for backupNode, _ := range backupNodes{
		cli := NewDataDataCli(backupNode)
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		_, err := cli.DataDataSyncDelete(ctx, &dataDataPb.DataDataSyncDeleteReq{
			Key: key,
		})
		if err != nil{
			utils.Error("Sync backup node %s in SyncBackupDelete error: %v\n", backupNode, err)
			return err
		}
		cancel()
	}
	return nil
}

func (dataServer *DataServer) MasterDataInformReshard(ctx context.Context, req *masterDataPb.MasterDataInformReshardReq) (*masterDataPb.MasterDataInformReshardResp, error){
	var keys []string
	for key, _ := range dataServer.database{
		keys = append(keys, key)
	}
	return &masterDataPb.MasterDataInformReshardResp{
		Keys: utils.StringArrayToByte(keys),
	}, nil
}

func (dataServer *DataServer) MasterDataReshardDestination(ctx context.Context, req *masterDataPb.MasterDataReshardDestinationReq) (*masterDataPb.MasterDataReshardDestinationResp, error){
	mapDestination := utils.ByteToKeyValueMap(req.KeyDestination)
	utils.Debug("Receive reshard instructions: %v\n", mapDestination)

	err := dataServer.doDataMigration(mapDestination)
	if err != nil{
		utils.Error("doDataMigration error in MasterDataReshardDestination: %v\n", err)
		return nil, err
	}
	return &masterDataPb.MasterDataReshardDestinationResp{
		Message: fmt.Sprintf("data server on port %s migration complete", dataServer.port),
	}, nil
}

func (dataServer *DataServer) DataDataMigrate(ctx context.Context, req *dataDataPb.DataDataMigrateReq) (*dataDataPb.DataDataMigrateResp, error){
	mapData := utils.ByteToKeyValueMap(req.KeyValues)
	utils.Debug("migrated data: %v\n", mapData)
	dataServer.database = mapData
	return &dataDataPb.DataDataMigrateResp{
		Message: "migration complete",
	}, nil
}

func (dataServer *DataServer) MasterDataRegisterBackupToData(ctx context.Context, req *masterDataPb.MasterDataRegisterBackupToDataReq) (*masterDataPb.MasterDataRegisterBackupToDataResp, error){
	backupPort := req.BackupPort
	log.Printf("Register backup node on port: %s\n", backupPort)
	dataServer.backupNodes[backupPort] = "true"
	cli := NewDataDataCli(backupPort)
	newCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	_, err := cli.DataDataSyncAll(newCtx, &dataDataPb.DataDataSyncAllReq{
		KeyValues: utils.KeyValueMapToByte(dataServer.database),
	})
	if err != nil{
		utils.Error("Sync database to %s error: %v\n", backupPort, err)
		return &masterDataPb.MasterDataRegisterBackupToDataResp{
			Message: "sync error",
		}, err
	}
	defer cancel()
	return &masterDataPb.MasterDataRegisterBackupToDataResp{
		Message: "register ok",
	}, err
}

func (dataServer *DataServer) MasterDataDeleteBackupOfData(ctx context.Context, req *masterDataPb.MasterDataDeleteBackupOfDataReq) (*masterDataPb.MasterDataDeleteBackupOfDataResp, error){
	backupPort := req.BackupPort
	log.Printf("delete backup node on port: %s\n", backupPort)
	delete(dataServer.backupNodes, backupPort)
	return &masterDataPb.MasterDataDeleteBackupOfDataResp{
		Message: "delete ok",
	}, nil
}

func (backupServer *BackupServer) DataDataSyncAll(ctx context.Context, req *dataDataPb.DataDataSyncAllReq) (*dataDataPb.DataDataSyncAllResp, error){
	keyValueMap := utils.ByteToKeyValueMap(req.KeyValues)

	backupServer.database = keyValueMap

	utils.Debug("Synced database: %v\n", backupServer.database)

	return &dataDataPb.DataDataSyncAllResp{
		Message: "Sync complete",
	}, nil
}

func (backupServer *BackupServer) DataDataSyncPut(ctx context.Context, req *dataDataPb.DataDataSyncPutReq) (*dataDataPb.DataDataSyncPutResp, error) {
	backupServer.database[req.Key] = req.Value
	utils.Debug("database after sync put: %v\n", backupServer.database)
	return &dataDataPb.DataDataSyncPutResp{
		Message: "sync put: ok",
	}, nil
}

func (backupServer *BackupServer) DataDataSyncDelete(ctx context.Context, req *dataDataPb.DataDataSyncDeleteReq) (*dataDataPb.DataDataSyncDeleteResp, error){
	delete(backupServer.database, req.Key)
	utils.Debug("database after sync delete: %v\n", backupServer.database)
	return &dataDataPb.DataDataSyncDeleteResp{
		Message: "sync delete: ok",
	}, nil
}

func NewDataDataCli(port string) dataDataPb.DataDataClient{
	dataAddress := fmt.Sprintf("%s%s", defaultIP, port)
	conn, err := grpc.Dial(dataAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect to master: %v", err)
	}
	dataDataCli := dataDataPb.NewDataDataClient(conn)
	return dataDataCli
}

func (dataServer *DataServer) doDataMigration(mapDestination map[string] string) error {
	migrateKeyValues := map[string] string{}
	dest := ""

	for key, destPort := range mapDestination{
		if destPort != dataServer.port{
			migrateKeyValues[key] = dataServer.database[key]
			if dest == ""{
				dest = destPort
			} else if dest != destPort{
				utils.Error("Data migration error: multiple destination port: %s and %s\n", dest, destPort)
				return errors.New("data migration error: multiple destination port")
			}

		}
	}

	//Need to migrate
	if dest != ""{
		cli := NewDataDataCli(dest)
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		resp, err := cli.DataDataMigrate(ctx, &dataDataPb.DataDataMigrateReq{
			KeyValues: utils.KeyValueMapToByte(migrateKeyValues),
		})
		if err != nil{
			utils.Error("Migrate data rpc request error: %v\n", err)
			return err
		}
		for migrateKey,_ := range migrateKeyValues{
			delete(dataServer.database, migrateKey)
			err := SyncBackupDelete(dataServer.backupNodes, migrateKey)
			if err != nil{
				utils.Error("Delete the data on backup nodes in doDataMigration error: %v\n", err)
				return err
			}
		}
		utils.Debug("Migrate data succeed, rpc response: %v\n", resp.Message)
	}
	return nil
}


//Promote the backup server as a rpc service.
//It initialize the dataserver, continue the heart beat detection, and delete the znode to stop the previous heart beat detection
//However, the previous grpc server(backup server) is not stopped. But as its port has been deleted from the backupNodeManager,
//there will be no request to the backup server.
func (backupServer *BackupServer) MasterBackupInformPromotion (ctx context.Context, req *masterDataPb.MasterBackupInformPromotionReq) (*masterDataPb.MasterBackupInformPromotionResp, error){
	log.Println("Try to promote to data server...")
	backupNodesList := utils.ByteToStringArray(req.BackupNodes)
	backupNodesSet := map[string] string{}
	for _, backupNode := range backupNodesList{
		if backupNode == "" { //Go feature: empty backupNodesList still have a length 1 with a "" string
			continue
		}
		backupNodesSet[backupNode] = "true"
	}

	utils.Debug("promoted data node's backup nodes: %v\n", backupNodesSet)
	dataServer := &DataServer{
		port: backupServer.dataPort,
		database: backupServer.database,
		backupNodes: backupNodesSet,
	}

	//continue the heartbeat response
	zkConn := ConnectZookeeper()
	go HeartBeatResponse(zkConn, backupServer.dataPort)


	//delete the backup's znode
	backupPath := fmt.Sprintf("%s/%s/%s", backupNodesPath, backupServer.dataPort, backupServer.backupPort)
	exist, s, err := zkConn.Exists(backupPath)
	if err != nil{
		utils.Error("Check the existence of the backup's znode error: %v\n", err)
		return nil, err
	} else if !exist{
		utils.Error("The znode of path %s is already deleted\n", backupPath)
		return nil, errors.New("the znode of the promoted backup server is already deleted")
	}
	err = zkConn.Delete(backupPath, s.Version)
	if err != nil{
		utils.Error("Can not delete the znode on path %s of the backup server: %v", backupPath, err)
		return nil, err
	}

	//Initialize dataServer rpc
	go func(dataPort string) error{
		lis, err := net.Listen("tcp", backupServer.dataPort)
		if err != nil{
			utils.Error("Backup server can not listen the data server port %s: %v\n", backupServer.dataPort, err)
			return err
		}
		grpcServer := grpc.NewServer()
		clientDataPb.RegisterClientDataServer(grpcServer, dataServer)
		masterDataPb.RegisterMasterDataServer(grpcServer, dataServer)
		dataDataPb.RegisterDataDataServer(grpcServer, dataServer)
		log.Println("Now serve as a data server")
		if err = grpcServer.Serve(lis); err != nil{
			utils.Error("Can not serve the listener on port %s: %v\n", dataPort, err)
			return err
		}
		return nil
	}(backupServer.dataPort)
	return &masterDataPb.MasterBackupInformPromotionResp{
		Message: "promote: ok",
	}, nil
}