package main

import (
	"FinalProject/lock"
	"FinalProject/utils"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

const (
	lockPath = "/testlock"
	path1 = "testlock/path1"
)
func thread1(c chan string, sum *int) {

	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("New lock for lock1")
	l := lock.NewLock(zkConn, lockPath)
	fmt.Printf("cli 1 id: %s\n", l.ID)
	for i:=0; i < 100000; i++{
		err = l.Lock()
		fmt.Printf("acquire lock for cli 1\n")
		if err != nil{
			fmt.Printf("acquire lock1 err: %v\n", err)
			return
		}
		*sum += 1
		fmt.Println(*sum)
		err = l.Unlock()
		fmt.Printf("release lock for cli 1\n")
		if err != nil{
			return
		}
	}
	c <- "t1 release"
}
func thread2(c chan string, sum *int) {

	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("New lock for lock2")
	l := lock.NewLock(zkConn, lockPath)
	fmt.Printf("cli 2 id: %s\n", l.ID)
	for i:=0; i < 100000; i++{
		err = l.Lock()
		fmt.Printf("acquire lock for cli 2\n")
		if err != nil{
			fmt.Printf("acquire lock2 err: %v\n", err)
			return
		}
		*sum += 1
		fmt.Println(*sum)
		err = l.Unlock()
		fmt.Printf("release lock for cli 2\n")
		if err != nil{
			return
		}
	}
	c <- "t2 release"
}

func testLock() {
	fmt.Println("go func test")
	c := make(chan string, 2)
	sum := 0
	go thread1(c, &sum)
	go thread2(c, &sum)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(sum)
	close(c)

}

func thread11(c chan string) {
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("New lock for lock2")
	l := lock.NewLock(zkConn, lockPath)
	l.Lock()
	l.Lock()
	c <- "t1 finished"
}

func thread21(c chan string) {
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("New lock for lock2")
	l := lock.NewLock(zkConn, lockPath)
	l.Unlock()

	c <- "t2 finished"
}

func testRespectiveLock(){
	c := make(chan string, 2)
	go thread11(c)
	time.Sleep(time.Second)
	go thread21(c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}

func testRwLockThreadRead(c chan string, ID int){
	port := ":1234"
	rwLock := lock.NewRwLock(port)
	err := rwLock.LockReader()
	if err != nil{
		log.Println(err)
	}
	for i := 1; i < 5; i++{
		fmt.Printf("thread %d read\n", ID)
		time.Sleep(time.Second)
	}
	err = rwLock.UnlockReader()
	if err != nil{
		log.Printf("release lock err: %v\n", err)
	}
	c<- fmt.Sprintf("thread %d finished", ID)
}

func testRwLockThreadWrite(c chan string, ID int){
	port := ":1234"
	rwLock := lock.NewRwLock(port)
	err := rwLock.LockWriter()
	if err != nil{
		log.Println(err)
	}
	for i := 1; i < 5; i++{
		fmt.Printf("thread %d write\n", ID)
		time.Sleep(time.Second)
	}
	err = rwLock.UnlockWriter()
	if err != nil{
		log.Printf("release lock err: %v\n", err)
	}
	c<- fmt.Sprintf("thread %d finished", ID)
}

func testRwLock(){
	threadNum := 5
	c := make(chan string, threadNum)
	go testRwLockThreadRead(c, 0)
	go testRwLockThreadRead(c, 1)
	go testRwLockThreadWrite(c, 2)
	go testRwLockThreadWrite(c, 3)
	go testRwLockThreadRead(c, 4)

	for i := 0; i < threadNum; i++{
		fmt.Println(<-c)
	}

}

func testNodeSeq(){
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	zkConn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}
	path := "/test/test"
	path0 := "/test"
	for i := 0; i < 10; i++{
		_, err = zkConn.CreateProtectedEphemeralSequential(path, []byte{}, zk.WorldACL(zk.PermAll))
		if err != nil{
			log.Fatal(err)
		}
		v, _, err := zkConn.Children(path0)
		if err != nil{
			log.Fatalf("Get children err: %v\n", err)
		}
		fmt.Println(v)
	}
}

func testUtilDebug(){
	utils.Debug("hello %d, %d", 1, 2)
}
func main(){
	testRwLock()
}
