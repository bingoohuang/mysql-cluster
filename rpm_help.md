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
    * 目标机：执行 ```sudo yum -y install `ls | grep rpm` ```
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

1. Thanks
    * [17 Jul 2018 Install HAProxy 1.8 on CentOS 7](https://pario.no/2018/07/17/install-haproxy-1-8-on-centos-7/)
    * [yum install rpm dependencies from a local directory without a localrepo](https://gist.github.com/ionutz22/ae5d4fae66cd81f27fd0f463ca4a015f)
