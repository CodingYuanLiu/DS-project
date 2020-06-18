package lock

import (
	"errors"
	"fmt"
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
)

type Lock struct{
	lockRootPath string
	ID string
	lockQueue *queue.Queue
	conn *zk.Conn
}

func NewLock(conn *zk.Conn, path string) *Lock{
	ID := uuid.New().String()
	return 	&Lock{
		lockRootPath: path,
		ID: ID,
		conn: conn,
		lockQueue: queue.New(0),
	}
}

func (l *Lock) Lock() error{
	exist, _, err := l.conn.Exists(l.lockRootPath)
	if !exist{
		return errors.New("the lock path does not exist")
	} else if err != nil{
		return err
	}
	//Add ":" between l.ID to split the ID
	newLockNodePath := fmt.Sprintf("%s/:%s:", l.lockRootPath, l.ID)
	_, err = l.conn.CreateProtectedEphemeralSequential(newLockNodePath, []byte{}, zk.WorldACL(zk.PermAll))
	if err != nil{
		return errors.New("create new lock node error " + err.Error())
	}
	for {
		locks, _, getCh, err := l.conn.ChildrenW(l.lockRootPath)
		if err != nil{
			return errors.New("watch pending locks error")
		}
		//log.Printf("[DEBUG] header lock: %v, my lock: %v\n", locks[0], l.ID)

		if len(locks) == 1 && strings.Split(locks[0],":")[1] == l.ID{
			//log.Printf("[DEBUG] %s acquire lock\n", newLockNodePath)
			return nil
		}

		select {
		case chEvent := <-getCh:
			if chEvent.Type == zk.EventNodeChildrenChanged {
				//log.Printf("[DEBUG] some one lock or unlock on the node")
			}
		}
	}
}

func (l *Lock) Unlock() error{
	//log.Printf("start to unlock...")
	exist, _, err := l.conn.Exists(l.lockRootPath)
	if !exist{
		return errors.New("the lock path does not exist")
	} else if err != nil{
		return err
	}
	locks, s, err := l.conn.Children(l.lockRootPath)
	if err != nil{
		log.Printf("unlock get children error\n")
		return err
	}
	releaseLockPath := fmt.Sprintf("%s/%s", l.lockRootPath, locks[0])
	if err := l.conn.Delete(releaseLockPath, s.Version); err != nil{
		log.Printf("unlock delete lock node error\n")
		return err
	}
	//log.Printf("[DEBUG] %s: release lock\n",releaseLockPath)
	return nil
}