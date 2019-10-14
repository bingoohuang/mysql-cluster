# MySQL和HAProxy的离线安装

## MySQL

1. 下载
    `wget https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm`
1. 安装MySQL源
    `sudo yum localinstall mysql57-community-release-el7-11.noarch.rpm`
1. 检查MySQL源是否安装成功

    ```bash
    [vagrant@bogon ~]$ sudo yum repolist enabled | grep "mysql.*-community.*"
    mysql-connectors-community/x86_64 MySQL Connectors Community 118
    mysql-tools-community/x86_64 MySQL Tools Community 95
    mysql57-community/x86_64 MySQL 5.7 Community Server 364
    ```

1. 下载全部依赖包到本地目录(vagrant centos7)
    * vagrant centos7上：安装插件 `sudo yum install yum-plugin-downloadonly`
    * vagrant centos7上：下载依赖包 `sudo yum install -y --downloadonly --downloaddir=/vagrant/mysql57  mysql-community-server`
        > [root@BJCA-device ~]# yum -h
        > -y, --assumeyes       回答全部问题为是
        > --downloadonly        仅下载而不更新
        > --downloaddir=DLDIR   指定一个其他文件夹用于保存软件包
    * 本机：上传目录 `sshpass -p mima scp -P1122 -o StrictHostKeyChecking=no ./*.rpm root@192.168.1.23:./mysql/`
    * 目标机：执行 `sudo rpm -ivh *.rpm --nodeps --force`

        ```bash
        [vagrant@bogon mysql57]$ sudo rpm -ivh *.rpm --nodeps --force
        警告：mysql-community-client-5.7.27-1.el7.x86_64.rpm: 头V3 DSA/SHA1 Signature, 密钥 ID 5072e1f5: NOKEY
        准备中...                          ################################# [100%]
        正在升级/安装...
        1:mysql-community-common-5.7.27-1.e################################# [ 17%]
        2:mysql-community-libs-5.7.27-1.el7################################# [ 33%]
        3:mysql-community-client-5.7.27-1.e################################# [ 50%]
        4:mysql-community-libs-compat-5.7.2################################# [ 67%]
        5:postfix-2:2.10.1-7.el7           ################################# [ 83%]
        6:mysql-community-server-5.7.27-1.e################################# [100%]
        ```

    * 目标机：开启启动 `systemctl enable mysqld`
    * 目标机：启动服务 `systemctl start mysqld`
    * 目标机：查看状态 `systemctl status mysqld`

1. 修改 root 本地账户密码

    安装完成之后，生成的默认密码在 /var/log/mysqld.log 文件中。

    使用 `grep 'temporary password' /var/log/mysqld.log` 命令找到日志中的密码。

    ```sql
    ALTER USER 'root'@'localhost' IDENTIFIED BY 'A1765527-61a0';
    ```

    > 注意：mysql 5.7 默认安装了密码安全检查插件（validate_password），默认密码检查策略要求密码必须包含：大小写字母、数字和特殊符号，并且长度不能少于8位。
    > 否则会提示 ERROR 1819 (HY000): Your password does not satisfy the current policy requirements 错误。

Thanks:

* [CentOS 7 下 MySQL 5.7 的安装与配置](https://www.jianshu.com/p/1dab9a4d0d5f)
* [Installing a package forcefully without dependencies](https://www.golinuxhub.com/2014/01/how-to-installuninstallupgrade-rpm.html)

## HAProxy

1. 使用阿里云的源，将下面脚本写入文件 init_aliyun_repo.sh，然后执行 `sudu sh init_aliyun_repo.sh`

    ```bash
    rm -rf /etc/yum.repos.d/*.repo
    wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-7.repo
    wget -O /etc/yum.repos.d/epel.repo https://mirrors.aliyun.com/repo/epel-7.repo
    sed -i '/aliyuncs/d' /etc/yum.repos.d/CentOS-Base.repo
    sed -i 's/http/https/g' /etc/yum.repos.d/CentOS-Base.repo
    sed -i 's/$releasever/7/g' /etc/yum.repos.d/CentOS-Base.repo
    sed -i '/aliyuncs/d' /etc/yum.repos.d/epel.repo
    sed -i 's/http/https/g' /etc/yum.repos.d/epel.repo
    ```

1. 从公网下载到本机
    * `sudo yum install yum-plugin-downloadonly`
    * `yum install centos-release-scl`
    * `yum install --downloadonly --downloaddir=/vagrant/haproxy18 rh-haproxy18-haproxy rh-haproxy18-haproxy-syspaths`

1. 从本机上传到目标机器
    * `sshpass -p mima scp -P1122 -o StrictHostKeyChecking=no ./*.rpm root@192.168.1.23:./haproxy/`

1. 在目标机器上安装
    * ```sudo yum -y install `ls | grep rpm` ```

1. 在目标机器上查看安装

    ```bash
    [root@BJCA-device ~]# more /usr/lib/systemd/system/rh-haproxy18-haproxy.service
    [Unit]
    Description=HAProxy Load Balancer
    After=network.target

    [Service]
    Environment="CONFIG=/etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg" "PIDFILE=/run/rh-haproxy18-haproxy.pid"
    EnvironmentFile=/etc/sysconfig/rh-haproxy18-haproxy
    ExecStartPre=/opt/rh/rh-haproxy18/root/usr/sbin/haproxy -f $CONFIG -c -q
    ExecStart=/opt/rh/rh-haproxy18/root/usr/sbin/haproxy -Ws -f $CONFIG -p $PIDFILE $OPTIONS
    ExecReload=/opt/rh/rh-haproxy18/root/usr/sbin/haproxy -f $CONFIG -c -q
    ExecReload=/bin/kill -USR2 $MAINPID
    KillMode=mixed
    Type=notify

    [Install]
    WantedBy=multi-user.target
    [root@BJCA-device ~]# ls -l /etc/haproxy/
    总用量 0
    lrwxrwxrwx 1 root root 44 10月 12 15:15 haproxy.cfg -> /etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg
    ```

1. 设置开机启动
    * 查看 `systemctl is-enabled rh-haproxy18-haproxy`
    * 设置 `systemctl enable rh-haproxy18-haproxy`

1. 检查状态

    ```bash
    [root@BJCA-device ~]# systemctl status rh-haproxy18-haproxy
    ● rh-haproxy18-haproxy.service - HAProxy Load Balancer
    Loaded: loaded (/usr/lib/systemd/system/rh-haproxy18-haproxy.service; enabled; vendor preset: disabled)
    Active: active (running) since 六 2019-10-12 15:41:34 CST; 7min ago
    Main PID: 31707 (haproxy)
    CGroup: /system.slice/rh-haproxy18-haproxy.service
            ├─31707 /opt/rh/rh-haproxy18/root/usr/sbin/haproxy -Ws -f /etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg -p /run/rh-haproxy18-haproxy.pid
            └─31708 /opt/rh/rh-haproxy18/root/usr/sbin/haproxy -Ws -f /etc/opt/rh/rh-haproxy18/haproxy/haproxy.cfg -p /run/rh-haproxy18-haproxy.pid

    10月 12 15:41:34 BJCA-device systemd[1]: Starting HAProxy Load Balancer...
    10月 12 15:41:34 BJCA-device haproxy[31707]: [WARNING] 284/154134 (31707) : config : log format ignored for proxy 'mysql-rw' since it has no log address.
    10月 12 15:41:34 BJCA-device systemd[1]: Started HAProxy Load Balancer.
    ```

1. `Starting proxy mysql-rw: cannot bind socket 127.0.0.1:13306`
    * `setsebool -P haproxy_connect_any=1` or Disable SELinux.
    * [Cannot bind socket](https://discourse.haproxy.org/t/solved-cannot-bind-socket/3180)
    * [How to Disable SELinux on CentOS 7](https://linuxize.com/post/how-to-disable-selinux-on-centos-7/)

1. Thanks
    * [17 Jul 2018 Install HAProxy 1.8 on CentOS 7](https://pario.no/2018/07/17/install-haproxy-1-8-on-centos-7/)
    * [yum install rpm dependencies from a local directory without a localrepo](https://gist.github.com/ionutz22/ae5d4fae66cd81f27fd0f463ca4a015f)

## HAProxy 编译安装

下载源代码[haproxy-1.8.21.tar.gz](https://www.haproxy.org/download/1.8/src/haproxy-1.8.21.tar.gz)，查看[其它可下载源代码](https://www.haproxy.org/#down)

在CentOS 7上(vagrant)

```bash
$ sudo yum -y install systemd-devel
...
已安装:
  systemd-devel.x86_64 0:219-67.el7_7.1
...
$ tar vxf /vagrant/haproxy-1.8.21.tar.gz
haproxy-1.8.21/
...
$ cd haproxy-1.8.21/
...
$ make ARCH=x86_64 TARGET=linux2628 USE_PCRE=1 USE_OPENSSL=1 USE_ZLIB=1 USE_SYSTEMD=1 USE_CPU_AFFINITY=1
gcc -m64 -march=x86-64 -g -o haproxy src/ev_poll.o src/ev_epoll.o src/ssl_sock.o ebtree/ebtree.o ebtree/eb32sctree.o ebtree/eb32tree.o ebtree/eb64tree.o ebtree/ebmbtree.o ebtree/ebsttree.o ebtree/ebimtree.o ebtree/ebistree.o src/proto_http.o src/cfgparse.o src/server.o src/stream.o src/flt_spoe.o src/stick_table.o src/stats.o src/mux_h2.o src/checks.o src/haproxy.o src/log.o src/dns.o src/peers.o src/standard.o src/sample.o src/cli.o src/stream_interface.o src/proto_tcp.o src/backend.o src/proxy.o src/tcp_rules.o src/listener.o src/flt_http_comp.o src/pattern.o src/cache.o src/filters.o src/vars.o src/acl.o src/payload.o src/connection.o src/raw_sock.o src/proto_uxst.o src/flt_trace.o src/session.o src/ev_select.o src/channel.o src/task.o src/queue.o src/applet.o src/map.o src/frontend.o src/freq_ctr.o src/lb_fwlc.o src/mux_pt.o src/auth.o src/fd.o src/hpack-dec.o src/memory.o src/lb_fwrr.o src/lb_chash.o src/lb_fas.o src/hathreads.o src/chunk.o src/lb_map.o src/xxhash.o src/regex.o src/shctx.o src/buffer.o src/action.o src/h1.o src/compression.o src/pipe.o src/namespace.o src/sha1.o src/hpack-tbl.o src/hpack-enc.o src/uri_auth.o src/time.o src/proto_udp.o src/arg.o src/signal.o src/protocol.o src/lru.o src/hdr_idx.o src/hpack-huff.o src/mailers.o src/h2.o src/base64.o src/hash.o   -lcrypt  -lz -ldl -lpthread  -lssl -lcrypto -ldl -lsystemd -L/usr/lib -lpcreposix -lpcre
$ mv haproxy /vagrant/
...
```

在本机：

```bash
➜ ls -lh haproxy
-rwxrwxr-x  1 bingoobjca  staff   7.8M 10 14 17:41 haproxy
➜ upx haproxy
                       Ultimate Packer for eXecutables
                          Copyright (C) 1996 - 2018
UPX 3.95        Markus Oberhumer, Laszlo Molnar & John Reiser   Aug 26th 2018

        File size         Ratio      Format      Name
   --------------------   ------   -----------   -----------
   8208006 ->   3236040   39.43%   linux/amd64   haproxy

Packed 1 file.
➜ ls -lh haproxy
-rwxrwxr-x  1 bingoobjca  staff   3.1M 10 14 17:41 haproxy
➜ sshpass -p mypwd scp -P8022 -o StrictHostKeyChecking=no ./haproxy root@192.168.1.22:.
...
```

在目标机器上：

```bash
# cp haproxy /usr/sbin/
# cat <<EOF | tee /usr/lib/systemd/system/haproxy.service
[Unit]
Description=HAProxy Load Balancer
After=syslog.target network.target

[Service]
Environment="CONFIG=/etc/haproxy.cfg" "PIDFILE=/run/haproxy.pid"
ExecStartPre=/usr/sbin/haproxy -f \$CONFIG -c -q
ExecStart=/usr/sbin/haproxy -Ws -f \$CONFIG -p \$PIDFILE
ExecReload=/usr/sbin/haproxy -f \$CONFIG -c -q
ExecReload=/bin/kill -USR2 \$MAINPID
KillMode=mixed
Type=notify

[Install]
WantedBy=multi-user.target
EOF

# cat <<EOF | tee /etc/haproxy.cfg
global
    maxconn 100000
    #stats socket /var/lib/haproxy/haproxy.sock mode 600 level admin
    daemon
    pidfile /run/haproxy.pid
    log 127.0.0.1 local3 info

defaults
    log global
    option http-keep-alive
    option  forwardfor
    maxconn 100000
    mode http
    timeout connect 30000ms
    timeout client  30000ms
    timeout server  30000ms

listen stats
    mode http
    bind 0.0.0.0:8000
    stats enable
    stats uri     /haproxy-status
    stats auth    admin:123456

listen  web_port
    bind 0.0.0.0:80
    mode http
    server web1  127.0.0.1:8080  check inter 3000 fall 2 rise 5
EOF
```

在目标机器上：

* 启动并验证haproxy `systemctl daemon-reload`, `systemctl start haproxy`
* 访问haproxy状态页 `http://92.168.1.22:8000/haproxy-status`

```bash
[root@BJCA-device ~]# systemctl status haproxy
● rh-haproxy18-haproxy.service - HAProxy Load Balancer
   Loaded: loaded (/usr/lib/systemd/system/rh-haproxy18-haproxy.service; disabled; vendor preset: disabled)
   Active: active (running) since 一 2019-10-14 18:00:23 CST; 1min 49s ago
  Process: 22798 ExecStartPre=/usr/sbin/haproxy -f $CONFIG -c -q (code=exited, status=0/SUCCESS)
 Main PID: 22801 (haproxy)
   CGroup: /system.slice/rh-haproxy18-haproxy.service
           ├─22801 /usr/sbin/haproxy -Ws -f /etc/haproxy.cfg -p /run/haproxy.pid
           └─22802 /usr/sbin/haproxy -Ws -f /etc/haproxy.cfg -p /run/haproxy.pid

10月 14 18:00:23 BJCA-device systemd[1]: Starting HAProxy Load Balancer...
10月 14 18:00:23 BJCA-device systemd[1]: Started HAProxy Load Balancer.
[root@BJCA-device ~]# systemctl restart haproxy
[root@BJCA-device ~]# systemctl status haproxy
● rh-haproxy18-haproxy.service - HAProxy Load Balancer
   Loaded: loaded (/usr/lib/systemd/system/rh-haproxy18-haproxy.service; disabled; vendor preset: disabled)
   Active: active (running) since 一 2019-10-14 18:02:20 CST; 1s ago
  Process: 22943 ExecStartPre=/usr/sbin/haproxy -f $CONFIG -c -q (code=exited, status=0/SUCCESS)
 Main PID: 22946 (haproxy)
   CGroup: /system.slice/rh-haproxy18-haproxy.service
           ├─22946 /usr/sbin/haproxy -Ws -f /etc/haproxy.cfg -p /run/haproxy.pid
           └─22947 /usr/sbin/haproxy -Ws -f /etc/haproxy.cfg -p /run/haproxy.pid

10月 14 18:02:20 BJCA-device systemd[1]: Starting HAProxy Load Balancer...
10月 14 18:02:20 BJCA-device systemd[1]: Started HAProxy Load Balancer.
```

参考

* [编译安装haproxy-1.8](https://www.s4lm0x.com/archives/116.html)
* [Configuration Manual version 1.8](http://cbonte.github.io/haproxy-dconv/1.8/configuration.html)
