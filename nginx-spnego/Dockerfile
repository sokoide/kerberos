FROM ubuntu:latest

MAINTAINER "So Koide" <sokoide@gmail.com>

WORKDIR /tmp

# Install prerequisites for Nginx compile
RUN DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y tzdata krb5-user
ENV TZ=Asia/Tokyo

RUN apt install -y \
    wget tar libssl-dev gcc g++ make zlib1g-dev libpcre3-dev libgd3 libgd-dev libkrb5-dev git libcap2-bin php-fpm

# Download nginx
RUN wget http://nginx.org/download/nginx-1.17.6.tar.gz -O nginx.tar.gz && \
    mkdir /tmp/nginx && \
    tar -xzvf nginx.tar.gz -C /tmp/nginx --strip-components=1 && \
    git clone https://github.com/stnoonan/spnego-http-auth-nginx-module.git /tmp/nginx/spnego-http-auth-nginx-module

# Build Nginx
WORKDIR /tmp/nginx
RUN ./configure \
    --user=nginx \
    --group=nginx \
    --prefix=/usr/local/nginx \
    --sbin-path=/usr/sbin/nginx \
    --conf-path=/usr/local/nginx/conf/nginx.conf \
    --pid-path=/run/nginx.pid \
    --lock-path=/run/lock/subsys/nginx \
    --error-log-path=/data/nginx/logs/error.log \
    --http-log-path=/data/nginx/logs/access.log \
    --with-http_gzip_static_module \
    --with-http_stub_status_module \
    --with-http_ssl_module \
    --with-pcre \
    --with-http_image_filter_module \
    --with-file-aio \
    --with-http_dav_module \
    --with-http_flv_module \
    --with-http_mp4_module \
    --with-http_gunzip_module \
    --add-module=spnego-http-auth-nginx-module && \
    make -j4 && \
    make install

# DATA VOLUMES
RUN mkdir -p /data/nginx/www/
RUN mkdir -p /data/nginx/logs/
RUN mkdir -p /data/nginx/conf/
RUN mkdir -p /data/etc

VOLUME ["/data"]

# Cleanup after Nginx build
RUN apt remove -y \
        wget \
        gcc \
        g++ \
        make \
        git && \
    apt autoremove -y && \
    rm -rf /tmp/*

# Configure filesystem to support running Nginx
RUN useradd -ms /bin/bash nginx && \
    setcap cap_net_bind_service=ep /usr/sbin/nginx

# Kerberos
ADD ./data/etc/krb5.conf /etc/krb5.conf
ADD ./data/etc/krb5.keytab /etc/krb5.keytab

ADD ./data/nginx/nginx-start.sh /opt/bin/nginx-start.sh
RUN chmod u=rwx /opt/bin/nginx-start.sh && \
    chown nginx:nginx /opt/bin/nginx-start.sh /data/nginx /usr/local/nginx/conf/nginx.conf /data/nginx/logs /usr/local/nginx /etc/krb5.*

# For troubleshooting
RUN apt install -y \
    inetutils-ping curl

# PORTS
EXPOSE 80
EXPOSE 443

USER nginx
ENTRYPOINT ["/opt/bin/nginx-start.sh"]
