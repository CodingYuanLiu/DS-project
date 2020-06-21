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
   
## Day5 
* 尝试实现扩容：1.知道reshard的地址。2.做reshard。3.加锁 4.写测试
* 注意:首先需要修改现有RPC的结构。DataMaster和MasterData要区分开。

## Day6
* 尝试实现容错
* TODO写在注释里面了。
  * backup需要注册，并创建/BackupNode/$dataPort/$backupPort 的znode节点，用来做心跳检测
  * master需要同时监听znode来注册backupnode，然后建立和backup node的心跳检测。
  * backup注册完成之后要和data server保持数据同步，因此注册的时候要首先做一次数据迁移。
  * 注册过程记得加锁，加锁，加锁。考虑好如何加锁
* 容错：
  * data挂了以后，心跳检测失败，如果有backup存在，则通知backup变成data server
  * 这个过程需要对port加写锁。如果有req先抢到锁，他们可能会失败返回error（或者更糟糕的是直接block住导致死锁）
  * 因此明天必须要先测试一下是会block还是会返回error
  
* 实现了容错里面的backup node的服务注册及心跳检测功能
  * BUG：如果一个data server挂了，那么在注册它的时候启动的watch backup node的goroutine不会挂，导致会有N个watcher在监听 backup 节点。这里面N-1个都是孤儿线程。
    * 如果data 死了以后能够重启，似乎不会创建新的watcher所以不会有这个问题?
    * 明天先检查一下有没有别的孤儿线程bug
* 目前使用一个backupServer数据结构监听backup节点用来做sync的port。目前想的是，在InitializeBackupServer里面注册好backup server以后，
  阻塞在一个channel里，如果通过zk收到了promote的message，那么就去要其他backupnode的port list，然后升级成data server
  * 这里面也要注意有没有孤儿线程的问题 
  
## Day7
* 实现容错里面的sync 和 重启功能
* master通过rpc告诉backup server需要promote，并且在参数里面告诉backup server他现在还有哪些backup节点
  * backup server直接在rpc处理里面，创建新的dataServer节点（database，注册dataServer rpc，并且通过goroutine去监听rpc的端口以及在原先data node的节点上面做心跳检测）
  * 同时，backup server把之前backup节点删掉，这样master在心跳检测的时候就知道backup server已经在promote了，就停掉对该backup server的心跳检测。同样，backup server的心跳检测response goroutine也能察觉到节点被删掉然后停掉。
    * （这里通过删节点传信息鲁棒性可能不太好？）
  * backup server保留了之前做backup server的时候listen的rpc节点端口（主要是不知道怎么关），但是因为backupNodemanager里面选择了这个backup server 来promote之后，就把这个port给删掉了，因此没有人知道这个backup server的port，从而这个port再也不会被访问到了。
* 通过了初步的测试：和之前相同的scalability，以及data server挂了以后backup server晋升为data server，仍然能read到之前的请求。

## Day8 收尾
* 创建一些不会自己荆楚或者需要手动创建的zookeeper状态
* 把log完善一下。有些log被写道utils.Debug里面去了。
* 补充一些额外的测试
  * 需要补充和优化的测试：并发测试，scalability测试，availability测试（这个不能像scalability一样实时）
* 写文档
* 可能需要自己抄一下zookeeper的dockerfile自己起一下zookeeper

## Notes
* 注意一些不会自己清除的状态或者需要手动创建的状态:
    1. 先shutdown master的话，DataNode不会删除
    2. /readers/$port里面存着readers的数量。错误的情况下可能运行完之后不为0
    3. /locks/$port里面存着用来加锁的节点。不过这个正常情况下会自行消除。
    4. /globalLock
    5. /readers/global
    
    
    
HeartBeatFailure => 删datanodeset，删data节点和backup节点 => datawatch注意到，但是旧的和新的一样所以没反应。