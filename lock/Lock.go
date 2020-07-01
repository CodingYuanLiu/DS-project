package lock

import (
	"FinalProject/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
)

type Lock struct{
	lockRootPath string
	ID string
	//lockQueue *queue.Queue
	conn *zk.Conn
}

func NewLock(conn *zk.Conn, path string) *Lock{
	ID := uuid.New().String()
	return 	&Lock{
		lockRootPath: path,
		ID: ID,
		conn: conn,
		//lockQueue: queue.New(0),
	}
}

func (l *Lock) Lock() error{
	exist, _, err := l.conn.Exists(l.lockRootPath)
	if !exist{
		utils.Debug("the lock path does not exist")
		if _, err := l.conn.Create(l.lockRootPath, []byte{}, 0, zk.WorldACL(zk.PermAll)); err != nil{
			log.Printf("create lock path %v err: %v\n", l.lockRootPath, err)
			return err
		}
	} else if err != nil{
		return err
	}
	//Add ":" between l.ID to split the ID
	newLockNodePath := fmt.Sprintf("%s/:%s:", l.lockRootPath, l.ID)
	newLockFullPath, err := l.conn.CreateProtectedEphemeralSequential(newLockNodePath, []byte{}, zk.WorldACL(zk.PermAll))

	if err != nil{
		return errors.New("create new lock node error " + err.Error())
	}
	utils.Debug("set lock %s\n", newLockFullPath)

	for {
		locks, _, getCh, err := l.conn.ChildrenW(l.lockRootPath)
		if err != nil{
			return errors.New("watch pending locks error")
		}
		//utils.Debug("header lock: %v, my lock: %v\n", locks[GetFirstLockNodeIndex(locks)], l.ID)

		firstLockNodeSplitByColon := strings.Split(locks[GetFirstLockNodeIndex(locks)], ":")
		firstLockSequentialNum := firstLockNodeSplitByColon[len(firstLockNodeSplitByColon) - 1]
		thisLockNodeSplitByColon := strings.Split(newLockFullPath, ":")
		thisLockSequentialNum := thisLockNodeSplitByColon[len(thisLockNodeSplitByColon) - 1]
		//The last element of the splitted node name is the lock sequential number
		if firstLockSequentialNum == thisLockSequentialNum{
			utils.Debug("%s acquire lock\n", newLockFullPath)
			return nil
		}

		select {
		case chEvent := <-getCh:
			if chEvent.Type == zk.EventNodeChildrenChanged {
				//utils.Debug("[DEBUG] some one lock or unlock on the node")
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
	if len(locks) == 0{
		return errors.New("unlock error: there's no locked node in " + l.lockRootPath)
	}
	releaseLockPath := fmt.Sprintf("%s/%s", l.lockRootPath, locks[GetFirstLockNodeIndex(locks)])

	if err := l.conn.Delete(releaseLockPath, s.Version); err != nil && err != zk.ErrNoNode{
		log.Printf("unlock delete lock node error\n")
		return err
	}
	utils.Debug("%s: release lock\n",releaseLockPath)
	return nil
}

func GetFirstLockNodeIndex(nodes []string) (index int){
	index = 0
	minSeq := strings.Split(nodes[0],":")[2]
	for i := index + 1; i < len(nodes); i++{
		Seq := strings.Split(nodes[i],":")[2]
		if  Seq < minSeq{
			index = i
			minSeq = Seq
		}
	}
	return
}