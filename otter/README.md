# otter

1. `docker-compose -f zk.yml up -d` startup and then open `http://localhost:2900/`
2. `docker-compose -f zk.yml ps`
3. `docker-compose -f zk.yml rm -fsv`
4. `echo stat | nc localhost 2181`
5. `echo "ruok" | nc localhost 2181 ; echo`  verify if zookeeper is running, expected response: imok

![image](https://user-images.githubusercontent.com/1940588/78862319-ac5aad80-7a69-11ea-911c-134e7e1c9b02.png)

```bash
# echo "ruok" | nc localhost 2181 ; echo
imok
```

## Zookeeper ports usage

from [stackoverflow](https://stackoverflow.com/a/18186224):

```properties
clientPort=2181
server.1=zookeeper1:2888:3888
server.2=zookeeper2:2888:3888
server.3=zookeeper3:2888:3888
```

Out of these one server will be the master and rest all will be slaves.If any server goes OFF then zookeeper elects leader automatically .

Servers listen on three ports:

1. `2181` for client connections;
2. `2888` for follower connections, if they are the leader;
3. `3888` for other server connections during the leader election phase.

otter Dockerfile:

move dockerfile to the target dir of the otter 
move otter-node/aria2c to the target 
```
 docker build -t footstone-otter-manager:v0.0.1 -f ManagerDockerfile .
 docker build -t footstone-otter-node:v0.0.1 -f NodeDockerfile .
```
