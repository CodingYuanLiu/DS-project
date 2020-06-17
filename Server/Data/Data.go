package main

import(
	"FinalProject/lock"
	clientDataPb "FinalProject/proto/ClientData"
	"context"
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
	port1 = ":7777"
	dataNodesPath = "/DataNode" //The root node of the data nodes' ports
	aliveResp = "Alive"
	aliveReq = "Is alive?"
)

type ClientDataServer struct{
	clientDataPb.UnimplementedClientDataServer
	rwLock	lock.RwLock
	database map[string] string
}

func (clientDataServer *ClientDataServer) ClientDataPut(ctx context.Context, req *clientDataPb.ClientDataPutReq) (*clientDataPb.ClientDataPutResp, error){
	err := clientDataServer.rwLock.LockWriter()
	if err != nil{
		return nil, err
	}

	clientDataServer.database[req.Key] = req.Value

	err = clientDataServer.rwLock.UnlockWriter()
	if err != nil{
		return nil, err
	}

	log.Printf("put key: %v, value: %v\n", req.Key, req.Value)
	log.Println(clientDataServer.database)
	return &clientDataPb.ClientDataPutResp{
		Message: "[Data server]: put succeed",
	}, nil
}

func (clientDataServer *ClientDataServer) ClientDataRead(ctx context.Context, req *clientDataPb.ClientDataReadReq) (*clientDataPb.ClientDataReadResp, error){
	err := clientDataServer.rwLock.LockReader()
	if err != nil{
		return nil, err
	}
	value, exist := clientDataServer.database[req.Key]

	if !exist{
		err := clientDataServer.rwLock.UnlockReader()
		if err != nil{
			return nil, err
		}
		return nil, errors.New("no value in the database")
	}

	err = clientDataServer.rwLock.UnlockReader()
	if err != nil{
		return nil, err
	}

	log.Printf("read key: %v, value: %v\n", req.Key, value)
	log.Println(clientDataServer.database)
	return &clientDataPb.ClientDataReadResp{
		Value: value,
		Message: "[Data server]: read succeed",
	}, nil
}

func (clientDataServer *ClientDataServer) ClientDataDelete(ctx context.Context, req *clientDataPb.ClientDataDeleteReq) (*clientDataPb.ClientDataDeleteResp, error){
	err := clientDataServer.rwLock.LockWriter()
	if err != nil{
		return nil, err
	}
	_, exist := clientDataServer.database[req.Key]
	if !exist{
		err = clientDataServer.rwLock.UnlockWriter()
		if err != nil{
			return nil, err
		}
		return nil, errors.New("no value in the database")
	}
	delete(clientDataServer.database, req.Key)
	err = clientDataServer.rwLock.UnlockWriter()
	if err != nil{
		return nil, err
	}
	return &clientDataPb.ClientDataDeleteResp{
		Message: "[Data server]: delete succeed",
	}, nil

}

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
	return nil
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
	zkConn := ConnectZookeeper()

	err = ZkRegisterDataNodePort(zkConn, port)
	go HeartBeatResponse(zkConn, port)
	if err != nil{
		log.Fatalf("Register data node to zookeeper error: %v\n", err)
	}

	clientDataPb.RegisterClientDataServer(dataServer, &ClientDataServer{
		database: map[string]string{},
		rwLock: lock.NewRwLock(),
	})
	if err = dataServer.Serve(lis); err != nil{
		log.Fatal(err)
	}
}
