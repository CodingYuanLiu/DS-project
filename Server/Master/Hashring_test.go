package main

import (
	"log"
	"testing"
)

const (
	node1 = "192.168.1.1"
	node2 = "192.168.1.2"
	node3 = "192.168.1.3"
)

func getNodesCount(nodes nodesArray) (int, int, int) {
	node1Count := 0
	node2Count := 0
	node3Count := 0

	for _, node := range nodes {
		if node.nodeKey == node1 {
			node1Count += 1
		}
		if node.nodeKey == node2 {
			node2Count += 1

		}
		if node.nodeKey == node3 {
			node3Count += 1

		}
	}
	return node1Count, node2Count, node3Count
}

func TestHash(t *testing.T) {

	nodeWeight := make(map[string]int)
	nodeWeight[node1] = 2
	nodeWeight[node2] = 2
	nodeWeight[node3] = 3
	vitualSpots := 100
	hash := NewHashRing(vitualSpots)
	t.Logf("empty get: %s\n", hash.GetNode("1"))

	hash.AddNodes(nodeWeight)
	if hash.GetNode("1") != node3 {
		t.Fatalf("expetcd %v got %v", node3, hash.GetNode("1"))
	}
	if hash.GetNode("2") != node3 {
		t.Fatalf("expetcd %v got %v", node3, hash.GetNode("2"))
	}
	if hash.GetNode("3") != node2 {
		t.Fatalf("expetcd %v got %v", node2, hash.GetNode("3"))
	}
	c1, c2, c3 := getNodesCount(hash.nodes)
	t.Logf("len of nodes is %v after AddNodes node1:%v, node2:%v, node3:%v", len(hash.nodes), c1, c2, c3)

	hash.RemoveNode(node3)
	if hash.GetNode("1") != node1 {
		t.Fatalf("expetcd %v got %v", node1, hash.GetNode("1"))
	}
	if hash.GetNode("2") != node2 {
		t.Fatalf("expetcd %v got %v", node1, hash.GetNode("2"))
	}
	if hash.GetNode("3") != node2 {
		t.Fatalf("expetcd %v got %v", node2, hash.GetNode("3"))
	}
	c1, c2, c3 = getNodesCount(hash.nodes)
	t.Logf("len of nodes is %v after RemoveNode node1:%v, node2:%v, node3:%v", len(hash.nodes), c1, c2, c3)

	hash.AddNode(node3, 3)
	if hash.GetNode("1") != node3 {
		t.Fatalf("expetcd %v got %v", node3, hash.GetNode("1"))
	}
	if hash.GetNode("2") != node3 {
		t.Fatalf("expetcd %v got %v", node3, hash.GetNode("2"))
	}
	if hash.GetNode("3") != node2 {
		t.Fatalf("expetcd %v got %v", node2, hash.GetNode("3"))
	}
	c1, c2, c3 = getNodesCount(hash.nodes)
	t.Logf("len of nodes is %v after AddNode node1:%v, node2:%v, node3:%v", len(hash.nodes), c1, c2, c3)

}

func TestHash_My(t *testing.T){
	nodeWeight := make(map[string]int)
	n1 := ":7777"
	n2 := ":2200"
	nodeWeight[n1] = 1
	vitualSpots := 100
	hash := NewHashRing(vitualSpots)
	t.Logf("empty get: %s\n", hash.GetNode("key1"))
	hash.AddNodes(nodeWeight)
	log.Printf("get 1: %v\n", hash.GetNode("testkey1"))
	log.Printf("get 2: %v\n", hash.GetNode("testkey2"))
	log.Printf("get 3: %v\n", hash.GetNode("testkeya"))

	hash.AddNode(n2, 1)
	log.Printf("Add node\n")
	log.Printf("get 1: %v\n", hash.GetNode("testkey1"))
	log.Printf("get 2: %v\n", hash.GetNode("testkey2"))
	log.Printf("get 3: %v\n", hash.GetNode("testkeya"))

	hash.RemoveNode(n2)
	log.Printf("Remove node\n")
	log.Printf("get 1: %v\n", hash.GetNode("testkey1"))
	log.Printf("get 2: %v\n", hash.GetNode("testkey2"))
	log.Printf("get 3: %v\n", hash.GetNode("testkeya"))

}