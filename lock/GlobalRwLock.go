package lock

import (
	"FinalProject/utils"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strconv"
	"time"
)


type GlobalRwLock struct{
	readerLock *Lock
	writerLock *Lock
	conn *zk.Conn
}


func NewGlobalRwLock() GlobalRwLock {
	readerLockPath := fmt.Sprintf("%s%s", globalLockPath, readerLockPath)
	writerLockPath := fmt.Sprintf("%s%s", globalLockPath, writerLockPath)
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	conn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}

	return GlobalRwLock{
		readerLock: NewLock(conn, readerLockPath),
		writerLock: NewLock(conn, writerLockPath),
		conn: conn,
	}
}

func GetGlobalReaderNum(conn *zk.Conn) (int, error){
	path := fmt.Sprintf("%s%s", ReaderNumRootPath, globalReaderPath)
	v, _, err := conn.Get(path)
	if err != nil{
		utils.Error("reader path: %s\n", path)
		return -1, err
	}
	return strconv.Atoi(string(v[:]))//Byte[] to string to int
}

func SetGlobalReaderNum(conn *zk.Conn, reader int) error{
	path := fmt.Sprintf("%s%s", ReaderNumRootPath, globalReaderPath)
	exist, s, err := conn.Exists(path)
	if !exist{
		return errors.New("the node does not exist")
	} else if err != nil{
		return err
	}
	_, err = conn.Set(path, []byte(strconv.Itoa(reader)), s.Version)
	if err != nil{
		log.Printf("[error] set reader num %d error: %v, version: %d\n", reader, err, s.Version)
	}
	return err
}

func (l *GlobalRwLock) LockReader() error{
	if err := l.readerLock.Lock(); err != nil{
		return err
	}
	reader, err := GetGlobalReaderNum(l.conn)
	if err != nil{
		log.Printf("[error] GetGlobalReaderNum in LockReader error: %v\n", err)
		return err
	}
	reader += 1
	//log.Printf("update reader: %d\n", reader)
	if reader == 1{
		log.Printf("[DEBUG] reader == 1, lock writer with ID: %s\n", l.writerLock.ID)
		if err := l.writerLock.Lock(); err != nil{
			log.Printf("[error] writerLock.Lock in LockReader error: %v\n", err)
			return err
		}
	}
	utils.Debug("[DEBUG] set reader num in lock reader %d", reader)
	if err := SetGlobalReaderNum(l.conn, reader); err != nil{
		log.Printf("[error] SetReaderNum in LockReader error: %v\n", err)
		return err
	}
	if err := l.readerLock.Unlock(); err != nil{
		log.Printf("[error] readerLock.Unlock in LockReader error: %v\n", err)
		return err
	}
	return nil
}

func (l *GlobalRwLock) UnlockReader() error{
	if err := l.readerLock.Lock(); err != nil{
		return err
	}
	reader, err := GetGlobalReaderNum(l.conn)
	if err != nil{
		log.Printf("[error] GetGlobalReaderNum in UnlockReader error: %v\n", err)
		return err
	}
	reader -= 1

	if reader == 0{
		log.Printf("[DEBUG] reader == 0, unlock writer\n")
		if err := l.writerLock.Unlock(); err != nil{
			log.Printf("[error] writerlock.Unlock in UnlockReader error: %v\n", err)
			return err
		}
	}

	utils.Debug("[DEBUG] set reader num in unlock reader %d", reader)
	if err := SetGlobalReaderNum(l.conn, reader); err != nil{
		log.Printf("[error] SetReaderNum in UnlockReader error: %v\n", err)
		return err
	}

	if err := l.readerLock.Unlock(); err != nil{
		log.Printf("[error] readerLock.Unlock in UnlockReader error: %v\n", err)
		return err
	}
	return nil
}

func (l *GlobalRwLock) LockWriter() error{
	utils.Debug("lock writerlock in LockWriter, id: %s", l.writerLock.ID)
	if err := l.writerLock.Lock(); err != nil{
		log.Printf("[error] writerLock.Lock in LockWriter error: %v\n", err)
		return err
	}
	return nil
}

func (l *GlobalRwLock) UnlockWriter() error{
	if err := l.writerLock.Unlock(); err != nil{
		log.Printf("[error] writerLock.Unlock in UnlockWriter error: %v\n", err)
		return err
	}
	return nil
}