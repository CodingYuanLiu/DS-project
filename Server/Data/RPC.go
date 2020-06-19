package main

import (
	clientDataPb "FinalProject/proto/ClientData"
	dataMasterPb "FinalProject/proto/DataMaster"
	masterDataPb "FinalProject/proto/MasterData"
	"FinalProject/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

type DataServer struct{
	clientDataPb.UnimplementedClientDataServer
	masterDataPb.UnimplementedMasterDataServer

	port string
	database map[string] string
	dataMasterCli dataMasterPb.DataMasterClient
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
	return &masterDataPb.MasterDataReshardDestinationResp{
		Message: fmt.Sprintf("data server on port %s migration complete", dataServer.port),
	}, nil
}

func NewDataMasterCli() dataMasterPb.DataMasterClient{
	masterAddress := fmt.Sprintf("%s%s", defaultIP, masterPort)
	conn, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect to master: %v", err)
	}
	dataMasterCli := dataMasterPb.NewDataMasterClient(conn)
	return dataMasterCli
}