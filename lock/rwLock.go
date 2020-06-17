package lock

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

const (
	lockPath          = "/locks"
	readerLockPath    = "/readerLock"
	writerLockPath    = "/writerLock"
	ReaderNumRootPath = "/readers"
)
type RwLock struct{
	readerLock *zk.Lock
	writerLock *zk.Lock
	port string
	conn *zk.Conn
}

func NewRwLock(port string) RwLock {
	readerLockPath := fmt.Sprintf("%s%s", lockPath, readerLockPath)
	writerLockPath := fmt.Sprintf("%s%s", lockPath, writerLockPath)
	hosts := []string{"localhost:2181", "localhost:2182", "localhost:2183"}

	conn, _, err := zk.Connect(hosts, time.Second * 5)
	if err != nil {
		fmt.Println(err)
	}

	return RwLock{
		readerLock: zk.NewLock(conn, readerLockPath, zk.WorldACL(zk.PermAll)),
		writerLock: zk.NewLock(conn, writerLockPath, zk.WorldACL(zk.PermAll)),
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
	return err
}

func (l *RwLock) LockReader() error{
	if err := l.readerLock.Lock(); err != nil{
		return err
	}
	reader, err := GetReaderNum(l.conn, l.port)
	if err != nil{
		return err
	}
	reader += 1
	if reader == 1{
		if err := l.writerLock.Lock(); err != nil{
			return err
		}
	}
	if err := SetReaderNum(l.conn, l.port, reader); err != nil{
		return err
	}
	if err := l.readerLock.Unlock(); err != nil{
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
		return err
	}
	reader -= 1

	if reader == 0{
		if err := l.writerLock.Unlock(); err != nil{
			return err
		}
	}

	if err := SetReaderNum(l.conn, l.port, reader); err != nil{
		return err
	}

	if err := l.readerLock.Unlock(); err != nil{
		return err
	}
	return nil
}

func (l *RwLock) LockWriter() error{
	if err := l.writerLock.Lock(); err != nil{
		return err
	}
	return nil
}

func (l *RwLock) UnlockWriter() error{
	if err := l.writerLock.Unlock(); err != nil{
		return err
	}
	return nil
}