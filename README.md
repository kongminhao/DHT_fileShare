# 基于DHT网络的局域网文件分享 

1. 提供windows,Linux, freebsd, Macos系统上的可执行文件(包含64位与32位的）
2. 一个daemon进程, 负责对路由表进行维护, 以及对dht网络中的信息进行自动化的收集。
3. 一个client进程，负责与daemon进程进行沟通，以完成下载功能，发送文件信息到DHT网络中

使用了一个简化的DHT， 只存储了前者节点和后者节点的路由信息。使用了一个16位的无符号整形作为nodeid， 保证了小型局域网之内无重复。

## daemon 进程
### 节点初始化 

从操作系统中读取出本地的ip地址，然后以时间作为随机数的种子, 随机一个Nodeid， 然后监听8765端口（udp), 同时启动一个rpc调用端口，监听2333（tcp）端口

### 路由表构造 

为了保证即使不知道DHT节点的情况，也可以加入局域网中的DHT网络，我在这里使用了一个广播地址，节点在进入dht网络之前先执行一次ping_all() 操作， 保证每个节点都能收到我的ping请求,然后对节点的返回信息做处理，新开一个协程对节点的路由表做更新操作。
我使用简单的nodeid之差的绝对值来衡量节点间的距离

### ping 

对DHT网络的路由表的前后节点进行探测。

### announce_peer 

对dht网络中的announce_peer, 宣告我这里有一个端口，有infohash对应的数据

### get_peers

对dht网络，获取其中的拥有infohash对应数据的ip及端口 

### find_node

寻找节点（未实现，因为在这个需求中不需要用到它） 

### 自动获取DHT网络中的infohash值 

通过每10s向DHT网络中发送broadcastinfo消息， 由路由表进行转发， 获取DHT网络中的infohash值 

当然，对应上面所有处理的handle也写了

### 错误处理 
如果路由表有明显错误，比如，ping超时无返回值， 或者broadcastinfo从一个不在自己路由表中的节点转发过来, 此时认为路由表出错, 重建路由表
## client进程 

连接到daemon进程的rpc调用端口，2333（tcp） 
记得一定要先开daemon进程，否则client进程连不上。

### /q

退出client

### /help 

显示帮助信息

### /upload filepath

将一个文件上传到DHT网络中 

### /showinfo 

/showinfo ， 展示当前dht网络中拥有的资源 

### /download infohash filename 

将infohash 对应的文件下载到本地，命名为filename




# 用户手册 

首先打开daemon进程,

![](http://caiji.scuseek.com/396ddecd017a020e9ae0d0a9ba6d5e78.png)

可以看到这两个文件的位数是不同的. 

然后我们执行它们, 这里我还开了一台windows主机做测试， 所以我们总共拥有3台主机。

我们在32位的linux上执行 client ， 并且上传一个文件

![](http://caiji.scuseek.com/7e46ddceecedf7fbe58cd9710371ac11.png)

![](http://caiji.scuseek.com/ab486dd4d786c1c27d7d5503fb8192dc.png)

然后我们获取了这个dht网络中拥有的infohash值和文件， 然后我们选择下载image.tar.gz 

![](http://caiji.scuseek.com/804d078898bb25b0d728bd95a638fcb6.png)

ok, download成功。




 
