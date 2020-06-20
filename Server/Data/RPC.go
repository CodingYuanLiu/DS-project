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

	//TODO: Sync to backup nodes
	return &clientDataPb.ClientDataPutResp{
		Message: "[Data server]: put succeed",
	}, nil
}

func (dataServer *DataServer) ClientDataRead(ctx context.Context, req *clientDataPb.ClientDataReadReq) (*clientDataPb.ClientDataReadResp, error){
	value, exist := dataServer.database[req.Key]
	if !exist{
		return nil, errors.New("no value in the database")
	}
	utils.Debug("read key: %v, value: %v\n", req.Key, value)
	utils.Debug("database now: %v\n", dataServer.database)
	//TODO: Sync to backup nodes
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

	//TODO: Sync to backup nodes
	return &clientDataPb.ClientDataDeleteResp{
		Message: "[Data server]: delete succeed",
	}, nil

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
	utils.Debug("Register backup node on port: %s\n", backupPort)
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
	utils.Debug("delete backup node on port: %s\n", backupPort)
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
	//TODO: Need to migrate to backup nodes????
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
		}
		utils.Debug("Migrate data succeed, rpc response: %v\n", resp.Message)
	}
	return nil
}