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
 * zookeeper实现简单lock service实现读写锁。
 