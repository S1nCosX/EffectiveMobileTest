#!/bin/sh

echo $CONSUL_ADDR:$CONSUL_PORT;

consul-template \
    -consul-addr=$CONSUL_ADDR:$CONSUL_PORT \
    -template="/etc/haproxy/haproxy.cfg.tmpl:/etc/haproxy/haproxy.cfg"\
    -once

haproxy -W -db -f /etc/haproxy/haproxy.cfg &
HAPROXY_PID=$!

consul-template \
    -consul-addr=$CONSUL_ADDR:$CONSUL_PORT \
    -template="/etc/haproxy/haproxy.cfg.tmpl:/etc/haproxy/haproxy.cfg:kill -SIGUSR2 $HAPROXY_PID"\
    -wait=$HEALTH_CHECK_TIMEOUT

wait "$HAPROXY_PID"
