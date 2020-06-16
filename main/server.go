package main

import (
	"fmt"
	"time"
)

func thread1() time.Time{
	resultChan1 := make(chan time.Time)
	go func(){
		sum := 0
		for i:= 0;i<1000000000;i++{
			sum += i
			sum -= i
			sum += i
			sum -= i
		}
		resultChan1 <- time.Now()
	}()
	return <-resultChan1
}
func thread2() time.Time{
	resultChan2 := make(chan time.Time)
	go func(){
		sum := ""
		for i:= "A"; len(i) < 10; i += "1"{
			sum += i
		}
		resultChan2 <- time.Now()
	}()
	return <-resultChan2
}

func main() {
	fmt.Println("go func test")

	fmt.Println("\n")
	time1 := thread1()
	time2 := thread2()
	fmt.Println(time1)
	fmt.Println(time2)
}

