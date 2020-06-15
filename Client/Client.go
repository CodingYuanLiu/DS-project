package main

import (
	clientDataPb "FinalProject/proto/ClientData"
	clientMasterPb "FinalProject/proto/ClientMaster"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

const(
	masterAddress = "localhost:7000"
)

//Client interface
type Client struct{
	rpcMasterCli  clientMasterPb.ClientMasterClient
	rpcDataCli  map[string] clientDataPb.ClientDataClient //key-value : port-client
}

func NewClient() *Client{
	conn, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := clientMasterPb.NewClientMasterClient(conn)
	return &Client{
		rpcMasterCli: client,
		rpcDataCli: map[string] clientDataPb.ClientDataClient{},
	}
}

func (cli *Client) Put(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	masterResp, err := cli.rpcMasterCli.ClientMasterPut(ctx, &clientMasterPb.ClientMasterPutReq{
		Key: "World",
	})
	if err != nil{
		log.Fatalf("Get data node from master.exe error: %v\n", err)
		return err
	}
	dataCli := cli.GetDataCli(masterResp.Port)
	dataResp, err := dataCli.ClientDataPut(ctx, &clientDataPb.ClientDataPutReq{
		Key: key,
		Value: value,
	})
	if err != nil{
		log.Fatalf("Put key-value to data node error: %v\n", err)
		return err
	}
	log.Printf("Put %v:%v to data node (port \"%v\") succeed\nmessage: %v\n",
		key, value, masterResp.Port, dataResp.Message)
	return nil

}

func (cli *Client) GetDataCli(port string) clientDataPb.ClientDataClient{
	if client, exist := cli.rpcDataCli["port"]; exist {
		return client
	}
	log.Printf("New connection to data node at port %v\n", port)
	defaultIP := "localhost"
	dataNodeAddr := fmt.Sprintf("%v%v", defaultIP, port)
	conn, err := grpc.Dial(dataNodeAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := clientDataPb.NewClientDataClient(conn)
	cli.rpcDataCli["port"] = client

	return client
}

func main() {
	cli := NewClient()
	if err := cli.Put("1", "value1"); err != nil{
		log.Fatal(err)
	}
	if err := cli.Put("2", "value2"); err != nil{
		log.Fatal(err)
	}
}
