#!/bin/bash

cp /data/nginx/conf/nginx.conf /usr/local/nginx/conf
cp /data/etc/krb5.conf /etc/krb5.conf
cp /data/etc/krb5.keytab /etc/krb5.keytab

exec /usr/sbin/nginx
