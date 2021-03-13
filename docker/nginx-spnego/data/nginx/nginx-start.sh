#!/bin/sh

echo "copying nginx.conf..."
cp /data/nginx/conf/nginx.conf /usr/local/nginx/conf
echo "removing existing logs..."
rm -rf /data/nginx/logs/*
echo "log dir: ./docker/nginx-spnego/data/nginx/logs"
echo "starting nginx..."
exec /usr/local/nginx/sbin/nginx
