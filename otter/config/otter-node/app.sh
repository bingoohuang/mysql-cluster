#!/bin/bash
set -e

# waitterm
#   wait TERM/INT signal.
#   see: http://veithen.github.io/2014/11/16/sigterm-propagation.html
waitterm() {
        local PID
        # any process to block
        tail -f /dev/null &
        PID="$!"
        # setup trap, could do nothing, or just kill the blocker
        trap "kill -TERM ${PID}" TERM INT
        # wait for signal, ignore wait exit code
        wait "${PID}" || true
        # clear trap
        trap - TERM INT
        # wait blocker, ignore blocker exit code
        wait "${PID}" 2>/dev/null || true
}

cp /otter/file /bin/
cp /otter/aria2c /bin/
rm -fr /app/bin/*.pid
rm -fr /app/log/*
bash /app/bin/startup.sh
echo $NID > /app/conf/nid;
cp /otter/otter.properties /app/conf/otter.properties
/otter/wait-for mr:2901 -- echo OK;
tail -f /dev/null &
# wait TERM signal
waitterm
