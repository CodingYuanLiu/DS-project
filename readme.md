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
    * 因此，我们让每个client都拥有一把读锁和写锁，指向zookeeper，这样就可以进行分布式加锁放锁了。但是，读写锁需要维护一个reader的数量，用来判断是否需要lock/unlock写锁。因此这个全局的reader数量把它放到一个zookeeper一个全局节点里面，用于进程间共享。这个节点的初始化由data/master server来做(master server初始化根节点`/readers`, data server初始化自己的port, 这设计的不是很好，后面可以改一下）。
    * 不知道这个测试应该怎么写比较好。瞎jb写了一个测试，似乎是可以跑了。
    
## Day4
 * 实现自己的锁。
   * Challenge: 我们需要根据节点创建的顺序来判断谁应该拿到锁。但是直接使用ChildrenW或者Children读出来的node顺序是乱的。比如我先插了node1和node2，有可能Children读出来的顺序就是node2先于node1.
 * 尝试实现data 节点的扩容。扩容要考虑
   1. 怎么加锁：使用一把全局的读写锁。每个client在向master请求RPC之前，无论是READ还是PUT，DELETE，都申请一把全局的读锁。而Master在注册的时候，申请一把全局的写锁。这样，只有没有client请求的时候，master会注册新的节点。
   2. 怎么reshard：
   
   
## Notes
* 注意一些不会自己清除的状态:
    1. 先shutdown master的话，DataNode不会删除
    2. readers/$port里面存着readers的数量。错误的情况下可能运行完之后不为0
    3. locks/$port里面存着用来加锁的节点。不过这个正常情况下会自行消除。