package main

const(
	masterPort = ":7000"
	dataNodesPath = "/DataNode" //The root node of the data nodes' ports
	backupNodesPath = "/BackupNode"
)

const (
	heartBeatDetectionFailureBound = 3
	aliveResp = "Alive"
	aliveReq = "Is alive?"
	heartBeatTimeInterval = 2
)