master:
	go build -o build/master.exe Server/Master/Master.go Server/Master/DataNodeManager.go Server/Master/Hashring.go

data:
	go build -o build/data.exe Server/Data/Data.go Server/Data/RPC.go

client:
	go build -o build/client.exe Client/Client.go

all: client server
