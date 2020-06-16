package main

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

const (
	heartBeatDetectionFailureBound = 3
	aliveResp = "Alive"
	aliveReq = "Is alive?"
	heartBeatTimeInterval = 2
)

type DataNodeManager struct{
	conn *zk.Conn
	hashRing *HashRing
	//nodeMeta map[uuid.UUID] *DataNode //TODO: UUID may be useless
	//nodeMeta map[string] masterDataPb.MasterDataClient //TODO: may be more useful than uuid
	dataNodesSet map[string] int //The set of data nodes. The int value is use to do heart beat detection.
	nodeNum int
}

/*
type DataNode struct{
	port   string
	client masterDataPb.MasterDataClient
}
*/

func NewDataNodeManager(conn *zk.Conn) *DataNodeManager{
	hashRing := NewHashRing(100)
	//nodeMeta := map[uuid.UUID]*DataNode{}
	nodeNum := 0
	return &DataNodeManager{
		hashRing: hashRing,
		//nodeMeta: nodeMeta,
		nodeNum: nodeNum,
		dataNodesSet: map[string]int{},
		conn: conn,
	}
}

/*
func NewDataNode(port string) *DataNode{
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	cli := masterDataPb.NewMasterDataClient(conn)
	return &DataNode{
		port: port,
		client: cli,
	}
}
*/

func (dataNodeManager DataNodeManager)DeleteDataNode(port string) error{
	delete(dataNodeManager.dataNodesSet, port)
	dataNodeManager.hashRing.RemoveNode(port)
	dataNodeManager.nodeNum -= 1
	return nil
}

func (dataNodeManager DataNodeManager)RegisterDataNode(port string) error{
	//TODO: lock the master.exe node and corresponding data node when register is undone

	/*
	//Register the Node at NodeMeta
	UUID := uuid.New()
	dataNodeManager.nodeMeta[UUID] = NewDataNode(port)
	*/
	dataNodeManager.dataNodesSet[port] = 0

	//Register the Node at hashRing
	dataNodeManager.hashRing.AddNode(port, 1)

	dataNodeManager.nodeNum += 1

	go dataNodeManager.HeartBeatDetection(port)

	//TODO: do the data migration
	return nil
}

func (dataNodeManager DataNodeManager)HeartBeatDetection(port string) error{
	nodePath := fmt.Sprintf("%v/%v", dataNodesPath, port)
	zkNodeExist, _, _ := dataNodeManager.conn.Exists(nodePath)
	_, dataNodeExist := dataNodeManager.dataNodesSet[port]
	if !(zkNodeExist && dataNodeExist){
		log.Printf("Error: no node to detect heart beat\n")
		return errors.New("no zk node or no data node in data node set")
	}

	//Do the heart beat detection
	for{
		value, s, err := dataNodeManager.conn.Get(nodePath)
		if err != nil{
			return err
		}

		//Heart beat detection failed
		if string(value[:]) != aliveResp{
			if dataNodeManager.dataNodesSet[port] < heartBeatDetectionFailureBound - 1{
				dataNodeManager.dataNodesSet[port] += 1
				log.Printf("data node on port %v failed %d times, node value: %v\n", port, dataNodeManager.dataNodesSet[port], string(value[:]))
			} else{
				log.Printf("data node on port %v utterly failed, delete it\n", port)
				if err := dataNodeManager.conn.Delete(nodePath, s.Version); err != nil{
					return err
				}
				delete(dataNodeManager.dataNodesSet, port)
				break
			}
		} else{
			dataNodeManager.conn.Set(nodePath, []byte(aliveReq),s.Version)
		}
		time.Sleep(time.Duration(heartBeatTimeInterval)*time.Second)
	}

	return nil
}

func (dataNodeManager DataNodeManager) HandleDataNodesChanges(ports []string) error{
	//log.Printf("Node num: %d, ports: %v\n", dataNodeManager.nodeNum, ports)
	//Node delete
	if len(ports) < dataNodeManager.nodeNum{
		newPortsMap := map[string] bool{}
		for _, port := range ports{
			newPortsMap[port] = true
		}
		//Find the port in dataNodesSet but not in the new ports
		for oldPort, _ := range dataNodeManager.dataNodesSet{
			if _, exist := newPortsMap[oldPort]; !exist{
				//Delete the node
				log.Printf("Delete node on port: %v\n", oldPort)
				err := dataNodeManager.DeleteDataNode(oldPort)
				if err != nil{
					return err
				}
			}
		}
	} else if len(ports) > dataNodeManager.nodeNum{
		//Node register: register unexist nodes
		for _, port := range ports{
			if _, exist := dataNodeManager.dataNodesSet[port]; !exist{
				log.Printf("Register node on port: %v\n", port)
				err := dataNodeManager.RegisterDataNode(port)
				if err != nil{
					return err
				}
			}
		}
	} else{
		return errors.New("no register and no delete, why")
	}
	return nil
}


func (dataNodeManager DataNodeManager)FindDataNode(key string) (string, error){
	return dataNodeManager.hashRing.GetNode(key), nil
}