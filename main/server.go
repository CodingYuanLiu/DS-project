package main

import (
	"FinalProject/lock"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
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

func main() {
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

