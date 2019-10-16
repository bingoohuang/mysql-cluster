#!/usr/bin/env bash

useradd -r haproxy -s /sbin/nologin
cp haproxy /usr/sbin/
cp haproxy.service /usr/lib/systemd/system/
cp haproxy.cfg /etc/haproxy.cfg
touch /run/haproxy.pid
chown haproxy:haproxy /run/haproxy.pid
systemctl daemon-reload
systemctl enable haproxy
systemctl start haproxy
