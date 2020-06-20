package main

import (
	"FinalProject/lock"
	masterDataPb "FinalProject/proto/MasterData"
	"FinalProject/utils"
	"context"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"time"
)



type DataNodeManager struct{
	conn *zk.Conn
	hashRing *HashRing
	dataNodesSet map[string] *DataNode //The set of data nodes. The int value is use to do heart beat detection.
	nodeNum      int
	globalRwLock lock.GlobalRwLock
	backupNodeManager *BackupNodeManager
}


type DataNode struct{
	heartBeatFailures   int
	client masterDataPb.MasterDataClient
}

func GetDataCli(port string) masterDataPb.MasterDataClient{
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	cli := masterDataPb.NewMasterDataClient(conn)
	return cli
}

func NewDataNode(port string) *DataNode{
	cli := GetDataCli(port)
	return &DataNode{
		heartBeatFailures: 0,
		client: cli,
	}
}

func NewDataNodeManager(conn *zk.Conn) (*DataNodeManager, error){
	hashRing := NewHashRing(100)
	//nodeMeta := map[uuid.UUID]*DataNode{}
	nodeNum := 0
	_, err := conn.Create(lock.ReaderNumRootPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNodeExists{
		log.Printf("Create readers root node: node already exists")
	} else if err != nil{
		return nil, err
	}
	backupNodeManager, err := NewBackupNodeManager(conn)
	if err != nil{
		utils.Error("NewBackupNodeManager in NewDataNodeManager error: %v\n", err)
		return nil, err
	}
	return &DataNodeManager{
		hashRing: hashRing,
		//nodeMeta: nodeMeta,
		nodeNum: nodeNum,
		dataNodesSet: map[string] *DataNode{},
		conn: conn,
		globalRwLock: lock.NewGlobalRwLock(),
		backupNodeManager: backupNodeManager,
	}, nil
}

func (dataNodeManager DataNodeManager)DeleteDataNode(port string) error{
	//TODO: lock the master node and data node when deletion is undone
	delete(dataNodeManager.dataNodesSet, port)
	dataNodeManager.hashRing.RemoveNode(port)
	dataNodeManager.nodeNum -= 1
	//TODO: do the migration (backup deployment)
	//TODO: unlock

	return nil
}

func (dataNodeManager DataNodeManager)RegisterDataNode(port string) error{
	err := dataNodeManager.globalRwLock.LockWriter()
	if err != nil{
		utils.Error("LockWriter error in RegisterDataNode: %v\n", err)
		return err
	}
	//Register the Node at hashRing
	dataNodeManager.hashRing.AddNode(port, 1)
	err = dataNodeManager.DataReshard()
	if err != nil{
		utils.Error("DataReshard in RegisterDataNode error: %v\n", err)
	}

	dataNodeManager.nodeNum += 1
	dataNodeManager.dataNodesSet[port] = NewDataNode(port)

	//Initialization about backup
	//dataNodeManager.InitializeBackupZnode(port)
	go dataNodeManager.backupNodeManager.WatchNewBackupNodes(port)

	go dataNodeManager.HeartBeatDetection(port)
	err = dataNodeManager.globalRwLock.UnlockWriter()
	if err != nil{
		utils.Error("LockWriter error in RegisterDataNode: %v\n", err)
		return err
	}
	return nil
}

func (dataNodeManager *DataNodeManager)InitializeBackupZnode(port string) error{
	path := fmt.Sprintf("%s/%s", backupNodesPath, port)
	exist, _, err := dataNodeManager.conn.Exists(path)
	if exist{
		utils.Debug("Backup root node of port %s is already created\n", port)
		return nil
	} else if err != nil && err != zk.ErrNoNode{
		utils.Error("Check backup root node of port %s existence error: %v\n", port, err)
		return err
	}
	_, err = dataNodeManager.conn.Create(path, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil{
		utils.Error("Can not create backup root node of port %s: %v ", port, err)
		return err
	}

	//TODO: Watch the node to register new backup nodes

	return nil
}
func (dataNodeManager DataNodeManager)DataReshard() error{
	for _, dataNode := range dataNodeManager.dataNodesSet{
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		informResp, err := dataNode.client.MasterDataInformReshard(ctx, &masterDataPb.MasterDataInformReshardReq{
		})
		if err != nil{
			utils.Error("inform client reshard error: %v\n", err)
			return err
		}
		newKeysDestinations := dataNodeManager.RehashKeys(utils.ByteToStringArray(informResp.Keys))
		ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		reshardResp, err := dataNode.client.MasterDataReshardDestination(ctx, &masterDataPb.MasterDataReshardDestinationReq{
			KeyDestination: utils.KeyValueMapToByte(newKeysDestinations),
		})
		if err != nil{
			utils.Error("reshard error: %v\n", err)
			return err
		}
		utils.Debug("Data server reshard response: %s\n", reshardResp.Message)
		cancel()
	}
	return nil
}

func (dataNodeManager DataNodeManager) RehashKeys(keys []string) map[string] string{
	RehashResult := map[string] string{}
	for _, key := range keys{
		RehashResult[key] = dataNodeManager.hashRing.GetNode(key)
	}
	return RehashResult
}

func (dataNodeManager DataNodeManager)HeartBeatDetection(port string) error{
	log.Printf("Start to do heart beat detection on port: %v\n", port)
	nodePath := fmt.Sprintf("%v/%v", dataNodesPath, port)
	zkNodeExist, _, _ := dataNodeManager.conn.Exists(nodePath)
	_, dataNodeExist := dataNodeManager.dataNodesSet[port]
	if !(zkNodeExist && dataNodeExist){
		utils.Error("No node to detect heart beat\n")
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
			if dataNodeManager.dataNodesSet[port].heartBeatFailures < heartBeatDetectionFailureBound - 1{
				dataNodeManager.dataNodesSet[port].heartBeatFailures += 1
				log.Printf("data node on port %v failed %d times, node value: %v\n", port, dataNodeManager.dataNodesSet[port].heartBeatFailures, string(value[:]))
			} else{
				log.Printf("data node on port %v utterly failed, delete it\n", port)
				_ = dataNodeManager.DeleteDataNode(port)
				//Delete the zk node
				if err := dataNodeManager.conn.Delete(nodePath, s.Version); err != nil{
					return err
				}
				break
			}
		} else{
			_, err := dataNodeManager.conn.Set(nodePath, []byte(aliveReq),s.Version)
			if err != nil{
				utils.Error("heart beat detection: set node error: %v\n", err)
				return err
			}
		}
		time.Sleep(time.Duration(heartBeatTimeInterval) * time.Second)
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
	}
	return nil
}

func (dataNodeManager *DataNodeManager) WatchNewDataNode(path string) error {
	conn := dataNodeManager.conn
	exist, _, err := conn.Exists(path)
	if err != nil{
		fmt.Println(err)
	}
	if !exist{
		_, err = conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
		log.Printf("Register Root data node in zookeeper: %s, node is created\n", path)
	}

	for {
		_, _, getCh, err := conn.ChildrenW(path)
		if err != nil {
			fmt.Printf("watch children error: %v\n", err)
		}

		select {
		case chEvent := <- getCh:
			{
				fmt.Printf("%+v\n", chEvent)
				if chEvent.Type == zk.EventNodeChildrenChanged {
					fmt.Printf("detect data node changed on zookeeper\n")
					v,_, err := conn.Children(path)
					if err != nil{
						return err
					}
					utils.Debug("value of path[%s]=[%s].\n", path, v)
					if err := dataNodeManager.HandleDataNodesChanges(v); err != nil{
						log.Printf("Handle data nodes change error: %v\n", err)
					}
				} else{
					fmt.Printf("other events on path %s\n", chEvent.Path)
				}
			}
		}
	}
}

func (dataNodeManager DataNodeManager)FindDataNode(key string) (string, error){
	dataNode := dataNodeManager.hashRing.GetNode(key)
	var err error = nil
	if dataNode == ""{
		err = errors.New("find no data string")
	}
	return dataNode, err
}