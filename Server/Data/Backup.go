package main

//The function will normally blocked, as a backup server or finally serve as a data server
func InitializeBackupServer(port string) error{
	//TODO: Create a znode on /BackupNode/$port of zookeeper
	//TODO: Watch the znode in a dead loop, response for heart beat detection, until master set the special flag on the znode
	//TODO: It also keeps synchronization with the data server in the dead loop

	//TODO: Initialize as a data server if it breaks from the dead loop.

	//TODO: Serve as a data server and blocked
}
