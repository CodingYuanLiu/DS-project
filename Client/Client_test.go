package main

import (
	"fmt"
	"testing"
)

func testClient(t *testing.T, ch chan string, id int) {
	cli := NewClient()
	if err := cli.Put("testkey1", "value1"); err != nil{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2"); err != nil{
		t.Error(err)
	}
	if value, err := cli.Read("testkey2"); err != nil || value != "value2"{
		t.Error(err)
	}
	if value, err := cli.Read("testkey2"); err != nil || value != "value2" && value != "value2.2"{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2.2"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("testkey2"); err != nil || value != "value2.2" {
		t.Error(err)
	}

	ch <- "test finish"
}
func TestConcurrentClient(t *testing.T){
	threadNum := 3
	ch := make(chan string)
	for i := 0; i < threadNum; i++{
		go testClient(t, ch, i)
	}
	for i := 0; i < threadNum; i++{
		fmt.Println(<-ch)
	}
}

func TestScalability(t *testing.T){
	//The intended test sequence: testClient(), new node register, testScalability
	cli := NewClient()
	if value, err := cli.Read("testkey2"); err != nil || value != "value2" && value != "value2.2"{
		t.Error(err)
	}
	if value, err := cli.Read("testkey2"); err != nil || value != "value2" && value != "value2.2"{
		t.Error(err)
	}
}