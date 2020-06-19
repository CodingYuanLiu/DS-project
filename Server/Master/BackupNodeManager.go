package main

import "github.com/samuel/go-zookeeper/zk"

type BackupNodeManager struct{
	conn zk.Conn
	backupNodeSet map[string] []string
}

//TODO: watch /BackupNode to initialize the root path of a data node's backup nodes??? (may can be inform directly by dataNodeManager)

//TODO: watch /BackupNode/$port to register new backup nodes

//TODO: do heart beat detection of the backup servers

//TODO: inform the backup data to "promote", and delete the corresponding znode at zookeeper