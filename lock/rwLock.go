package lock

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

const (
	lockPath = "/locks"
	readerLockPath = "/readerLock"
	writerLockPath = "/writerLock"
)
type RwLock struct{
	readerLock *zk.Lock
	writerLock *zk.Lock
	reader int
}

func NewRwLock() RwLock {
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
		reader: 0,
	}
}



func (l *RwLock) LockReader() error{
	if err := l.readerLock.Lock(); err != nil{
		return err
	}

	l.reader += 1
	if l.reader == 1{
		if err := l.writerLock.Lock(); err != nil{
			return err
		}
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
	l.reader -= 1
	if l.reader == 0{
		if err := l.writerLock.Unlock(); err != nil{
			return err
		}
	}
	if err := l.readerLock.Unlock(); err != nil{
		return err
	}
	return nil
}

func (l *RwLock) LockWriter() error{
	return l.writerLock.Lock()

}

func (l *RwLock) UnlockWriter() error{
	return l.writerLock.Unlock()
}