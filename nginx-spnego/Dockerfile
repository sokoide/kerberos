#FROM ubuntu:latest
FROM alpine

MAINTAINER "So Koide" <sokoide@gmail.com>

WORKDIR /tmp

# Install prerequisites for Nginx compile
RUN apk add --no-cache tzdata krb5
ENV TZ=Asia/Tokyo

# RUN apt install -y \
RUN apk add --no-cache \
    wget tar openssl-dev gcc g++ make krb5-dev git

# Download nginx
RUN wget http://nginx.org/download/nginx-1.17.6.tar.gz -O nginx.tar.gz && \
    mkdir /tmp/nginx && \
    tar -xzvf nginx.tar.gz -C /tmp/nginx --strip-components=1 && \
    git clone https://github.com/stnoonan/spnego-http-auth-nginx-module.git /tmp/nginx/spnego-http-auth-nginx-module

RUN apk add --no-cache \
    libc-dev \
    pcre-dev \
    zlib-dev \
    linux-headers \
    curl \
    gnupg \
    libxslt-dev \
    gd-dev \
    geoip-dev \
    perl-dev

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

# DATA VOLUME
RUN mkdir -p /data/nginx/www /data/nginx/logs /data/nginx/conf /data/etc
VOLUME ["/data"]

# Cleanup after Nginx build
RUN apk del --purge \
        wget \
        gcc \
        g++ \
        make \
        linux-headers \
        git && \
    rm -rf /tmp/*

RUN addgroup -S nginx && adduser -S nginx -G nginx

# Kerberos
ADD ./data/etc/krb5.conf /etc/krb5.conf
ADD ./data/etc/krb5.keytab /etc/krb5.keytab
ADD ./data/nginx/nginx-start.sh /opt/nginx-start.sh

RUN chmod u=rwx /opt/nginx-start.sh && \
    chown nginx:nginx -R /opt/nginx-start.sh /data/nginx /usr/local/nginx/conf \
      /data/nginx/logs /usr/local/nginx /etc /usr/sbin/nginx && \
    rm /usr/local/nginx/conf/nginx.conf /etc/krb5.conf /etc/krb5.keytab

# PORTS
EXPOSE 10080
EXPOSE 10443

USER nginx
ENTRYPOINT ["/opt/nginx-start.sh"]
