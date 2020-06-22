package lock

import (
	"FinalProject/utils"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)


type RwLock struct{
	readerLock *Lock
	writerLock *Lock
	port string
	conn *zk.Conn
}


func NewRwLock(port string) RwLock {
	readerLockPath := fmt.Sprintf("%s%s/%s", LockPath, ReaderLockPath, port)
	writerLockPath := fmt.Sprintf("%s%s/%s", LockPath, WriterLockPath, port)
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	conn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}

	return RwLock{
		readerLock: NewLock(conn, readerLockPath),
		writerLock: NewLock(conn, writerLockPath),
		port: port,
		conn: conn,
	}
}

func GetReaderNum(conn *zk.Conn, port string) (int, error){
	path := fmt.Sprintf("%s/%s", ReaderNumRootPath, port)
	v, _, err := conn.Get(path)
	if err != nil{
		return -1, err
	}
	return strconv.Atoi(string(v[:]))//Byte[] to string to int
}

func SetReaderNum(conn *zk.Conn, port string, reader int) error{
	path := fmt.Sprintf("%s/%s", ReaderNumRootPath, port)
	exist, s, err := conn.Exists(path)
	if !exist{
		return errors.New("the node does not exist")
	} else if err != nil{
		return err
	}
	_, err = conn.Set(path, []byte(strconv.Itoa(reader)), s.Version)
	if err != nil{
		utils.Error("set reader num %d error: %v, version: %d\n", reader, err, s.Version)
	}
	return err
}

func (l *RwLock) LockReader() error{
	if err := l.readerLock.Lock(); err != nil{
		return err
	}
	reader, err := GetReaderNum(l.conn, l.port)
	if err != nil{
		utils.Error("GetReaderNum in LockReader error: %v\n", err)
		return err
	}
	reader += 1
	//log.Printf("update reader: %d\n", reader)
	if reader == 1{
		utils.Debug("reader == 1, lock writer with ID: %s\n", l.writerLock.ID)
		if err := l.writerLock.Lock(); err != nil{
			utils.Error("writerLock.Lock in LockReader error: %v\n", err)
			return err
		}
	}
	utils.Debug("set reader num in lock reader %d", reader)
	if err := SetReaderNum(l.conn, l.port, reader); err != nil{
		utils.Error("SetReaderNum in LockReader error: %v\n", err)
		return err
	}
	if err := l.readerLock.Unlock(); err != nil{
		utils.Error("readerLock.Unlock in LockReader error: %v\n", err)
		return err
	}
	return nil
}

func (l *RwLock) UnlockReader() error{
	if err := l.readerLock.Lock(); err != nil{
		return err
	}
	reader, err := GetReaderNum(l.conn, l.port)
	if err != nil{
		utils.Error("[error] GetReaderNum in UnlockReader error: %v\n", err)
		return err
	}
	reader -= 1

	if reader == 0{
		utils.Debug("reader == 0, unlock writer\n")
		if err := l.writerLock.Unlock(); err != nil{
			utils.Error("writerlock.Unlock in UnlockReader error: %v\n", err)
			return err
		}
	}

	utils.Debug("set reader num in unlock reader %d", reader)
	if err := SetReaderNum(l.conn, l.port, reader); err != nil{
		utils.Error("SetReaderNum in UnlockReader error: %v\n", err)
		return err
	}

	if err := l.readerLock.Unlock(); err != nil{
		utils.Error("readerLock.Unlock in UnlockReader error: %v\n", err)
		return err
	}
	return nil
}

func (l *RwLock) LockWriter() error{
	utils.Debug("lock writerlock in LockWriter, id: %s", l.writerLock.ID)
	if err := l.writerLock.Lock(); err != nil{
		utils.Error("writerLock.Lock in LockWriter error: %v\n", err)
		return err
	}
	return nil
}

func (l *RwLock) UnlockWriter() error{
	if err := l.writerLock.Unlock(); err != nil{
		utils.Error("writerLock.Unlock in UnlockWriter error: %v\n", err)
		return err
	}
	return nil
}