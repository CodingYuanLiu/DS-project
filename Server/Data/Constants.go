package main

import "errors"

const(
	defaultDataPort = ":7777"
	dataNodesPath   = "/DataNode" //The root node of the data nodes' ports
	backupNodesPath = "/BackupNode"
	aliveResp       = "Alive"
	aliveReq        = "Is alive?"
	masterPort      = ":7000"
	defaultIP       = "localhost"
	messagePromote = "Promote"
)

var (
	ErrMessagePromote = errors.New("promote from backup server to data server")
)