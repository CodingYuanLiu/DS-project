package main


import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestBasicOperations(t *testing.T){
	cli := NewClient()
	if err := cli.Put("testkey1", "value1"); err != nil{
		t.Error(err)
	}
	if value, err := cli.Read("testkey1"); err != nil || value != "value1"{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2"); err != nil{
		t.Error(err)
	}
	if value, err := cli.Read("testkey2"); err != nil || value != "value2"{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2.2"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("testkey2"); err != nil || value != "value2.2"{
		t.Error(err)
	}

	if err := cli.Delete("testkey2"); err != nil{
		t.Error(err)
	}

	if _, err := cli.Read("testkey2"); err == nil{
		t.Error(errors.New("this should return an error"))
	}

	if err := cli.Put("testkey2", "value2"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("testkey2"); err != nil || value != "value2"{
		t.Error(err)
	}
}


func testClient(t *testing.T, ch chan string, id int) {
	cli := NewClient()
	if err := cli.Put("testkey1", "value1"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("testkey1"); err != nil && value != "value1"{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("testkey2"); err != nil || value != "value2" && value != "value2.2"{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2.2"); err != nil{
		t.Error(err)
	}

	if value, err := cli.Read("testkey2"); err != nil || value != "value2" && value != "value2.2" {
		t.Error(err)
	}

	ch <- "test finish"
}
func TestConcurrentClient(t *testing.T){
	threadNum := 10
	ch := make(chan string)
	for i := 0; i < threadNum; i++{
		go testClient(t, ch, i)
	}
	for i := 0; i < threadNum; i++{
		fmt.Println(<-ch)
	}
}

func testLockScalabilityThread(t *testing.T, ch chan string, ID int){
	cli := NewClient()
	if err := cli.Put("testkey1", "value1"); err != nil{
		t.Error(err)
	}

	if err := cli.Put("testkey2", "value2"); err != nil{
		t.Error(err)
	}
	for i := 0; i < 6; i++{
		if value, err := cli.Read("testkey1"); err != nil || value != "value1" {
			t.Error(err)
		}
		if value, err := cli.Read("testkey2"); err != nil || value != "value2" {
			t.Error(err)
		}
		fmt.Printf("Read repeatly for %d times\n", i)
		time.Sleep(time.Millisecond * 500)
	}
	ch <- fmt.Sprintf("thread %d done\n", ID)
}

func TestLockScalabilityConcurrently(t *testing.T){
	threadNum := 2
	ch := make(chan string)
	for i := 0; i < threadNum; i++{
		go testLockScalabilityThread(t, ch, i)
	}
	for i := 0; i < threadNum; i++{
		fmt.Println(<-ch)
	}
}

func TestAvailability(t *testing.T){
	//The test should be run after previous put operations
	cli := NewClient()
	if value, err := cli.Read("testkey1"); err != nil && value != "value1"{
		t.Error(err)
	}

	if value, err := cli.Read("testkey2"); err != nil || value != "value2" && value != "value2.2" {
		t.Error(err)
	}

}