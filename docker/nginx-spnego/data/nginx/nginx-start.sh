#!/bin/sh

cp /data/nginx/conf/nginx.conf /usr/local/nginx/conf
cp /data/etc/krb5.conf /etc/krb5.conf
cp /data/etc/krb5.keytab /etc/krb5.keytab

echo "starting nginx..."
echo "log dir: ./docker/nginx-spnego/data/nginx/logs"
exec /usr/local/nginx/sbin/nginx
