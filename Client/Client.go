package main

import (
	pb "FinalProject/proto/ClientMaster"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

const(
	masterAddress = "localhost:7000"
)

func main() {
	conn, err := grpc.Dial(masterAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewClientMasterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.ClientMasterPut(ctx, &pb.ClientMasterPutReq{
		Key: "World",
	})

	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Response message: %v\n", resp.Port)

	resp, err = client.ClientMasterPut(ctx, &pb.ClientMasterPutReq{
		Key: "World123",
	})

	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Response message: %v\n", resp.Port)

	resp, err = client.ClientMasterPut(ctx, &pb.ClientMasterPutReq{
		Key: "World324",
	})

	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Response message: %v\n", resp.Port)

	resp, err = client.ClientMasterPut(ctx, &pb.ClientMasterPutReq{
		Key: "World666",
	})

	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Response message: %v\n", resp.Port)
}
