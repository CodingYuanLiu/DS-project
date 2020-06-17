# DS-lab: final project
## Day1
 * 1 master + 2 data node + sequential requests
 * 没有考虑多client加锁问题
 * 没有考虑reshard问题。
   * 现在加节点是直接硬加，改变key对应的data node，但是没有迁移更改了node 的 key-value pair
 * 没有考虑注册data node的问题, data node是在master启动的时候静态注册上去的。
    * 正常情况应该是data node启动以后通知zookeeper去master注册。
 * 运行方式：先运行master节点和两个data节点（必须两个都运行，因为master静态注册了这俩data节点），然后跑client的测试。
 
## Day2
 * zookeeper实现简单动态nameservice注册，但不支持reshard，也不加锁，可能会导致并发错误。（后面补充）
 * 实现了心跳检测，暂定每2s心跳检测一次。如果一个data node 2s都未回应则认为他挂掉了，删除其zk里面的节点以及注册在hashring里面的节点
    * 检测方法：dataNodeManager修改data node在zookeeper里面节点的值为"Is alive?"，data node监听到这个修改之后将这个值改成"Alive". dataNodeManager在下一次修改前，读到这个value是"Alive"就知道检测成功了。
    * 没有考虑data node挂掉的时候是接受任务会导致任务被分配到挂掉的节点的问题。  
 * NameService的监听，心跳检测的发送和接受，使用了go routine。
 
## Day3 
 * 实现读写锁满足多client运行 
    * zookeeper的锁是以节点为准，而不是以lock object（指针）为准的，i.e., 即使是两个不同的lock object，如果指向同一个节点也可以lock。
    每次newLock生成的lock只能lock一次，否则将抛出死锁错误。所以我们不能使用共享的reader lock和writer lock
    * 对于reader lock，它的lock 和 unlock在一个函数里面，因此直接新生成一把锁去lock和unlock readerlock节点就行。
    * 对于writer lock，我们在锁它之后