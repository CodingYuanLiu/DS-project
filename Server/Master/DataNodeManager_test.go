package main

import (
	"fmt"
	"testing"
)

func TestDataNodeManager_FindDataNode(t *testing.T) {
	dataNodeManager := NewDataNodeManager()
	_ = dataNodeManager.RegisterDataNode(":7777")
	_ = dataNodeManager.RegisterDataNode(":2200")
	str, _ := dataNodeManager.FindDataNode("1")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("2")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("key1")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("key2")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("asdgzxb1")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("asdgzxb2")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("asdgzxb3")
	fmt.Println(str)
	str, _ = dataNodeManager.FindDataNode("asdgzxb4")
	fmt.Println(str)

}