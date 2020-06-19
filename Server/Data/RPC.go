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
	//dataMasterCli dataMasterPb.DataMasterClient
}

func (dataServer *DataServer) ClientDataPut(ctx context.Context, req *clientDataPb.ClientDataPutReq) (*clientDataPb.ClientDataPutResp, error){

	dataServer.database[req.Key] = req.Value

	log.Printf("put key: %v, value: %v\n", req.Key, req.Value)
	log.Println(dataServer.database)
	return &clientDataPb.ClientDataPutResp{
		Message: "[Data server]: put succeed",
	}, nil
}

func (dataServer *DataServer) ClientDataRead(ctx context.Context, req *clientDataPb.ClientDataReadReq) (*clientDataPb.ClientDataReadResp, error){
	value, exist := dataServer.database[req.Key]
	if !exist{
		return nil, errors.New("no value in the database")
	}
	log.Printf("read key: %v, value: %v\n", req.Key, value)
	log.Println(dataServer.database)
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
	//TODO: data migration
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
		utils.Debug("Migrate data succeed, rpc response: %v\n", resp.Message)
	}
	return nil
}