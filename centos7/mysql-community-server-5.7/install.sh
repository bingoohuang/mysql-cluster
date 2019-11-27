#!/usr/bin/env bash

yum install -y libaio*
yes | cp my.cnf /etc/my.cnf
rpm -ivh *.rpm --nodeps --force
chown mysql:mysql /var/lib/mysql
systemctl daemon-reload
systemctl enable mysqld
# systemctl is-enabled mysqld
systemctl start mysqld
grep 'temporary password' /var/lib/mysql/mysqld.log
password=$(grep -oP 'temporary password(.*): \K(\S+)' /var/lib/mysql/mysqld.log)
newpwd=$(grep -oP '\bPassword\s*=\s*"\K[^"]+' ../mysqlclusterinit/config.toml)
mysqladmin --user=root --password="$password" password "$newpwd"
