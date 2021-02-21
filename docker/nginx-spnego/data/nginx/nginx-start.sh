#!/bin/sh

cp /data/nginx/conf/nginx.conf /usr/local/nginx/conf
# krb5.conf and keytab are copied inside Dockerfile
# when you update it, please rebuild the image
# cp /data/etc/krb5.conf /etc/krb5.conf
# cp /data/etc/krb5.keytab /etc/krb5.keytab

echo "starting nginx..."
echo "log dir: ./docker/nginx-spnego/data/nginx/logs"
exec /usr/local/nginx/sbin/nginx
