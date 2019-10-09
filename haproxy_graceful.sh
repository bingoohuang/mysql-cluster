#!/bin/sh

# https://manageacloud.com/configuration/haproxy_graceful_restart

# hold/pause new requests
iptables -I INPUT -p tcp --dport 13306,23306 --syn -j DROP
sleep 1

# gracefully restart haproxy
/usr/sbin/haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -sf $(cat /var/run/haproxy.pid)

# allow new requests to come in again
iptables -D INPUT -p tcp --dport 13306,23306 --syn -j DROP
