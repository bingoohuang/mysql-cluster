# MySQL 安装

准备MySQL安装包：

1. mysql-community-client-5.7.27-1.el7.x86_64.rpm
1. mysql-community-common-5.7.27-1.el7.x86_64.rpm
1. mysql-community-libs-5.7.27-1.el7.x86_64.rpm
1. mysql-community-server-5.7.27-1.el7.x86_64.rpm

执行安装脚本 `sudo ./install.sh`

## 参考

1. 修改 root 本地账户密码
    安装完成之后，生成的默认密码在 /var/log/mysqld.log 文件中
    使用 `grep 'temporary password' /var/log/mysqld.log` 命令找到日志中的密码。

    ```bash
    [root@BJCA-device ~]# uuidgen
    d3472e8c-5885-4124-b7f2-6df505733f9d
    ```

    ```sql
    ALTER USER 'root'@'localhost' IDENTIFIED BY 'b7f2-6df505733f9D';
    ```

1. 使用 `mysql -uroot -p` 输入密码后，检查是否可以正常连接MySQL。
1. 参考 [Removing the MySQL root password](https://medium.com/@benmorel/remove-the-mysql-root-password-ba3fcbe29870)

    * On MySQL 5.7:

        ```bash
        password=$(grep -oP 'temporary password(.*): \K(\S+)' /var/log/mysqld.log)
        mysqladmin --user=root --password="$password" password aaBB@@cc1122
        mysql --user=root --password=aaBB@@cc1122 -e "UNINSTALL PLUGIN validate_password;"
        mysqladmin --user=root --password="aaBB@@cc1122" password ""
        ```

    * On MySQL 8.0:

        ```bash
        password=$(grep -oP 'temporary password(.*): \K(\S+)' /var/log/mysqld.log)
        mysqladmin --user=root --password="$password" password aaBB@@cc1122
        mysql --user=root --password=aaBB@@cc1122 -e "UNINSTALL COMPONENT 'file://component_validate_password';"
        mysqladmin --user=root --password="aaBB@@cc1122" password ""
        ```