package main

import(
	clientDataPb "FinalProject/proto/ClientData"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)
const(
	port1 = ":7777"
)

type ClientDataServer struct{
	clientDataPb.UnimplementedClientDataServer
	database map[string] string
}

func (clientDataServer *ClientDataServer) ClientDataPut(ctx context.Context, req *clientDataPb.ClientDataPutReq) (*clientDataPb.ClientDataPutResp, error){
	clientDataServer.database[req.Key] = req.Value
	log.Printf("put key: %v, value: %v\n", req.Key, req.Value)
	log.Println(clientDataServer.database)
	return &clientDataPb.ClientDataPutResp{
		Message: "[Data server]: put succeed",
	}, nil
}

func (clientDataServer *ClientDataServer) ClientDataRead(ctx context.Context, req *clientDataPb.ClientDataReadReq) (*clientDataPb.ClientDataReadResp, error){
	value, exist := clientDataServer.database[req.Key]

	if !exist{
		return nil, errors.New("no value in the database")
	}

	log.Printf("read key: %v, value: %v\n", req.Key, value)
	log.Println(clientDataServer.database)
	return &clientDataPb.ClientDataReadResp{
		Value: value,
		Message: "[Data server]: read succeed",
	}, nil
}

func (clientDataServer *ClientDataServer) ClientDataDelete(ctx context.Context, req *clientDataPb.ClientDataDeleteReq) (*clientDataPb.ClientDataDeleteResp, error){
	_, exist := clientDataServer.database[req.Key]
	if !exist{
		return nil, errors.New("no value in the database")
	}
	delete(clientDataServer.database, req.Key)
	return &clientDataPb.ClientDataDeleteResp{
		Message: "[Data server]: delete succeed",
	}, nil

}

func main() {
	port := port1
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	lis, err := net.Listen("tcp", port)
	fmt.Printf("Data node listening port %v...\n", port)

	if err != nil{
		log.Fatal(err)
	}
	dataServer := grpc.NewServer()

	clientDataPb.RegisterClientDataServer(dataServer, &ClientDataServer{
		database: map[string]string{},
	})
	if err = dataServer.Serve(lis); err != nil{
		log.Fatal(err)
	}
}
