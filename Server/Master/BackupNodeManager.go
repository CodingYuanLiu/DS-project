package main

import (
	masterDataPb "FinalProject/proto/MasterData"
	"FinalProject/utils"
	"context"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type BackupNodeManager struct{
	conn *zk.Conn
	backupNodeSet map[string] map[string] *BackupNode
}

type BackupNode struct{
	heartBeatFailures int
}

func NewBackupNodeManager(conn *zk.Conn) (*BackupNodeManager,error){
	exist, _, err := conn.Exists(backupNodesPath)
	if err != nil{
		utils.Error("Check the existence of backup root node error: %v\n", err)
		return nil, err
	} else if !exist{
		utils.Debug("Create backup root node")
		_, errC := conn.Create(backupNodesPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if errC != nil{
			utils.Error("Create backup root node error in NewBackupNodeManager: %v\n", err)
			return nil, errC
		}
	}
	return &BackupNodeManager{
		conn: conn,
		backupNodeSet: map[string] map[string] *BackupNode {},
	}, nil
}

//create and watch /BackupNode/$port to register new backup nodes
func (backupNodeManager *BackupNodeManager) WatchNewBackupNodes(port string) error{
	utils.Debug("Watch backup node of %s via zookeeper\n", port)
	backupRootPortPath := fmt.Sprintf("%s/%s", backupNodesPath, port)
	conn := backupNodeManager.conn
	exist, _, err := conn.Exists(backupRootPortPath)
	if err != nil{
		utils.Error("check znode on backupRootPortPath %s in WatchNewBackupNodes err: %v\n", backupRootPortPath, err)
		return err
	} else if !exist{
		_, err := conn.Create(backupRootPortPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil{
			utils.Error("create znode on backupRootPortPath %s in WatchNewBackupNodes err: %v", backupRootPortPath, err)
			return err
		}
	}
	for{
		_, _, getCh, err := conn.ChildrenW(backupRootPortPath)
		if err != nil {
			utils.Error("watch children error: %v\n", err)
		}
		select {
		case chEvent := <- getCh:
			{
				if chEvent.Type == zk.EventNodeChildrenChanged {
					utils.Debug("detect backup node of %s changed on zookeeper\n", port)
					backupNodes , _, err := conn.Children(backupRootPortPath)
					if err != nil{
						return err
					}
					if err := backupNodeManager.HandleBackupNodesChanges(backupNodes, port); err != nil{
						log.Printf("Handle data nodes change error: %v\n", err)
					}
				} else{
					fmt.Printf("other events on path %s\n", chEvent.Path)
				}
			}
		}
	}
}

func (backupNodeManager *BackupNodeManager) HandleBackupNodesChanges(backupPorts []string, dataPort string) error{
	oldBackupNodes, exist := backupNodeManager.backupNodeSet[dataPort]
	utils.Debug("old backups: %v, new backups: %v\n", oldBackupNodes, backupPorts)
	if !exist{
		utils.Debug("first backup")
		if len(backupPorts) != 1{
			utils.Debug("Warning: first backup node register of dataPort %s with more than 1 backup nodes: %v\n", dataPort, backupPorts)
		}
		oldBackupNodes = map[string] *BackupNode{}
		backupNodeManager.backupNodeSet[dataPort] = oldBackupNodes
	}
	//Need to register new backup nodes
	if len(backupPorts) > len(oldBackupNodes){
		for _, backupPort := range backupPorts {
			_, exist := oldBackupNodes[backupPort]
			if !exist{
				err := backupNodeManager.RegisterBackupNode(backupPort, dataPort)
				if err != nil{
					utils.Error("RegisterBackupNode %s of data node %s error: %v\n", backupPort, dataPort, err)
					return err
				}
			}
		}
	} else if len(backupPorts) < len(oldBackupNodes){
		backupNodePortsMap := map[string] bool{} // = set
		for _, backupNode := range backupPorts {
			backupNodePortsMap[backupNode] = true
		}
		for oldPort, _ := range oldBackupNodes{
			if _, exist := backupNodePortsMap[oldPort]; !exist{
				log.Printf("Delete node on dataPort: ")
				err := backupNodeManager.DeleteBackupNode(oldPort, dataPort)
				if err != nil{
					utils.Error("Delete backup node error: %v\n", err)
					return err
				}
			}
		}


	}
	return nil
}

func (backupNodeManager *BackupNodeManager) RegisterBackupNode(backupPort string, dataPort string) error{
	oldBackupNodes := backupNodeManager.backupNodeSet[dataPort]
	oldBackupNodes[backupPort] = &BackupNode{
		heartBeatFailures: 0,
	}
	go backupNodeManager.HeartBeatDetection(backupPort, dataPort)
	err := registerBackupToData(backupPort, dataPort)
	if err != nil{
		utils.Error("Register backup node to data node error: %v\n", err)
		return err
	}
	return nil
}

func registerBackupToData(backupPort string, dataPort string) error{
	cli := GetDataCli(dataPort)
	utils.Debug("Inform data node %s to sync...\n", dataPort)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	resp, err := cli.MasterDataRegisterBackupToData(ctx, &masterDataPb.MasterDataRegisterBackupToDataReq{
		BackupPort: backupPort,
	})
	if err != nil{
		utils.Error("Register backup node to data node via rpc error: %v\n", err)
		return err
	}
	utils.Debug("sync result: %s\n", resp.Message)

	return nil
}

func deleteBackupOfData(backupPort string, dataPort string) error{
	cli := GetDataCli(dataPort)
	utils.Debug("Inform data node %s to delete %s\n", dataPort, backupPort)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	resp, err := cli.MasterDataDeleteBackupOfData(ctx, &masterDataPb.MasterDataDeleteBackupOfDataReq{
		BackupPort: backupPort,
	})
	if err != nil{
		utils.Error("Delete backup node of the data node via rpc error: %v\n", err)
		return err
	}
	utils.Debug("Delete rpc result %s\n", resp.Message)
	return nil
}
func (backupNodeManager *BackupNodeManager) DeleteBackupNode(backupPort string, dataPort string) error{
	oldBackupNodes := backupNodeManager.backupNodeSet[dataPort]
	delete(oldBackupNodes, backupPort)
	//TODO: Delete the backup node in the data node
	err := deleteBackupOfData(backupPort, dataPort)
	if err != nil{
		utils.Error("Delete the backup node of the data error: %v", err)
		return err
	}

	return nil
}

func (backupNodeManager *BackupNodeManager) HeartBeatDetection(backupPort string, dataPort string) error{
	log.Printf("Start to do heart beat detection on backup node: %s/%s\n", dataPort, backupPort)
	conn := backupNodeManager.conn
	nodePath := fmt.Sprintf("%s/%s/%s", backupNodesPath, dataPort, backupPort)

	zkNodeExist, _, err := conn.Exists(nodePath)
	if err != nil{
		utils.Error("Check node existence in heart beat detection error: %v\n", err)
		return err
	} else if !zkNodeExist{
		utils.Error("No zk node to detect heart beat")
		return errors.New("no zk node to detect heart beat")
	}
	backupNodeSetOfPort := backupNodeManager.backupNodeSet[dataPort]
	for{
		value, s, err := conn.Get(nodePath)
		if err != nil{
			utils.Error("Get znode in HeartBeatDetection error: %v", err)
			return err
		}
		if string(value[:]) != aliveResp{
			if backupNodeSetOfPort[backupPort].heartBeatFailures < heartBeatDetectionFailureBound - 1{
				backupNodeSetOfPort[backupPort].heartBeatFailures += 1
				log.Printf("data node on port %s/%s failed %d times",
					dataPort, backupPort, backupNodeSetOfPort[backupPort].heartBeatFailures)
			} else{
				log.Printf("data node on port %s/%s utterly failed, delete it\n", dataPort, backupPort)
				_ = backupNodeManager.DeleteBackupNode(backupPort, dataPort)
				if err := backupNodeManager.conn.Delete(nodePath, s.Version); err != nil{
					utils.Error("Can not delete the node when heart beat detection failed: %v\n", err)
					return err
				}
				break
			}
		} else{
			_, err := conn.Set(nodePath, []byte(aliveReq), s.Version)
			if err != nil{
				utils.Error("heart beat detection: set node error: %v\n", err)
				return err
			}
		}
		time.Sleep(time.Duration(heartBeatTimeInterval) * time.Second)
	}
	return nil
}


//TODO: inform the backup data to "promote", and delete the corresponding znode at zookeeper