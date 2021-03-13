#!/bin/sh

cp /data/nginx/conf/nginx.conf /usr/local/nginx/conf
echo "starting nginx..."
echo "log dir: ./docker/nginx-spnego/data/nginx/logs"
rm -rf /data/nginx/logs/*
exec /usr/local/nginx/sbin/nginx
