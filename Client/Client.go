package main

import (
	"FinalProject/lock"
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
	rwLocks map[string] lock.RwLock //stores the locks of all the data nodes. key-value : port-lock
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
		rwLocks: map[string] lock.RwLock{},
	}
}

func (cli *Client) GetDataNodePort(key string) (string,error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	masterResp, err := cli.rpcMasterCli.ClientMasterFindDataNode(ctx, &clientMasterPb.ClientMasterFindDataNodeReq{
		Key: key,
	})
	if err != nil{
		return "", err
	}
	return masterResp.Port, nil
}
func (cli *Client) Put(key string, value string) error {
	port, err := cli.GetDataNodePort(key)
	if err != nil{
		log.Fatalf("Get data node from master error: %v\n", err)
		return err
	}

	dataCli := cli.GetDataCli(port)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	rwLock := cli.GetRwLock(port)
	if err := rwLock.LockWriter(); err!= nil{
		log.Printf("[error] lock writer error: %v\n", err)
		return err
	}

	dataResp, err := dataCli.ClientDataPut(ctx, &clientDataPb.ClientDataPutReq{
		Key: key,
		Value: value,
	})
	if err != nil{
		if err := rwLock.UnlockWriter(); err!= nil{
			log.Printf("[error] unlock writer error: %v\n", err)
			return err
		}
		log.Fatalf("Put key-value to data node error: %v\n", err)
		return err
	}

	if err := rwLock.UnlockWriter(); err!= nil{
		log.Printf("[error] unlock writer error: %v\n", err)
		return err
	}
	log.Printf("Put %v:%v to data node (port \"%v\") succeed\nmessage: %v\n",
		key, value, port, dataResp.Message)
	return nil
}

func (cli *Client) Read(key string) (string, error){
	port, err := cli.GetDataNodePort(key)
	if err != nil{
		log.Fatalf("Get data node from master.exe error: %v\n", err)
		return "", err
	}

	dataCli := cli.GetDataCli(port)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	rwLock := cli.GetRwLock(port)
	if err := rwLock.LockReader(); err!= nil{
		return "", err
	}

	dataResp, err := dataCli.ClientDataRead(ctx, &clientDataPb.ClientDataReadReq{
		Key: key,
	})
	if err != nil{
		log.Printf("Put key-value to data node error: %v\n", err)
		if err := rwLock.UnlockReader(); err!= nil{
			return "", err
		}
		return "", err
	}
	if err := rwLock.UnlockReader(); err!= nil{
		return "", err
	}
	value := dataResp.Value
	log.Printf("Read %v:%v from data node (port \"%v\") succeed\nmessage: %v\n",
		key, value, port, dataResp.Message)
	return value, nil
}

func (cli *Client) Delete(key string) error{
	port, err := cli.GetDataNodePort(key)
	if err != nil{
		log.Fatalf("Get data node from master error: %v\n", err)
		return err
	}

	dataCli := cli.GetDataCli(port)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	rwLock := cli.GetRwLock(port)
	if err := rwLock.LockWriter(); err!= nil{
		return err
	}

	dataResp, err := dataCli.ClientDataDelete(ctx, &clientDataPb.ClientDataDeleteReq{
		Key: key,
	})
	if err != nil{
		log.Printf("Put key-value to data node error: %v\n", err)
		if err := rwLock.UnlockWriter(); err!= nil{
			return err
		}
		return err
	}
	if err := rwLock.UnlockWriter(); err!= nil{
		return err
	}
	log.Printf("Delete %v from data node (port \"%v\") succeed\nmessage: %v\n",
		key, port, dataResp.Message)
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

func (cli *Client) GetRwLock(port string) lock.RwLock{
	rwLock, exist := cli.rwLocks[port]
	if !exist{
		rwLock = lock.NewRwLock(port)
		cli.rwLocks[port] = rwLock
	}
	return rwLock
}

/*
func main() {
	cli := NewClient()
	if err := cli.Put("1", "value1"); err != nil{
		log.Fatal(err)
	}
	if err := cli.Put("2", "value2"); err != nil{
		log.Fatal(err)
	}
}
*/