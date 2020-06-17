package main

import (
	"fmt"
	"testing"
)

func testClient(t *testing.T, ch chan string, id int) {
	cli := NewClient()

	if err := cli.Put("1", "value1"); err != nil{
		t.Error(err)
	}
	if err := cli.Put("2", "value2"); err != nil{
		t.Error(err)
	}
	if value, err := cli.Read("2"); err != nil || value != "value2"{
		t.Error(err)
	}
	if value, err := cli.Read("2"); err != nil || value != "value2" && value != "value2.2"{
		t.Error(err)
	}

	if err := cli.Put("2", "value2.2"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("2"); err != nil || value != "value2.2"{
		t.Error(err)
	}
	if err := cli.Delete("1"); err != nil{
		t.Log(err)
	}
	ch <- "test finish"
}

func TestConcurrentClient(t *testing.T){
	ch := make(chan string)
	go testClient(t, ch, 0)
	go testClient(t, ch, 1)
	go testClient(t, ch, 2)
	go testClient(t, ch, 3)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}