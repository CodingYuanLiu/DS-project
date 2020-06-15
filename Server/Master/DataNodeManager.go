package main

const (
	dataPort1 = ":7777" //TODO: zookeeper nameservice
	dataPort2 = ":2200"
)

type DataNodeManager struct{
	hashRing *HashRing
	//nodeMeta map[uuid.UUID] *DataNode //TODO: UUID may be useless
	nodeNum int
}

/*
type DataNode struct{
	port   string
	client masterDataPb.MasterDataClient
}
*/

func NewDataNodeManager() *DataNodeManager{
	hashRing := NewHashRing(100)
	//nodeMeta := map[uuid.UUID]*DataNode{}
	nodeNum := 0
	return &DataNodeManager{
		hashRing: hashRing,
		//nodeMeta: nodeMeta,
		nodeNum: nodeNum,
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

func (dataNodeManager DataNodeManager)RegisterDataNode(port string) error{
	//TODO: lock the master.exe node and corresponding data node when register is undone

	/*
	//Register the Node at NodeMeta
	UUID := uuid.New()
	dataNodeManager.nodeMeta[UUID] = NewDataNode(port)
	*/

	//Register the Node at hashRing
	dataNodeManager.hashRing.AddNode(port, 1)

	dataNodeManager.nodeNum += 1
	//TODO: do the data migration
	return nil
}

func (dataNodeManager DataNodeManager)FindDataNode(key string) (string, error){
	return dataNodeManager.hashRing.GetNode(key), nil
}