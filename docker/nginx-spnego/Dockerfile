FROM alpine:latest

MAINTAINER "So Koide" <sokoide@gmail.com>

WORKDIR /tmp

# Install prerequisites for Nginx compile
RUN apk add --no-cache tzdata krb5
ENV TZ=Asia/Tokyo \
	GOPATH=/go \
	PATH=$GOPATH/bin:$PATH

# DATA VOLUME
RUN mkdir -p /go /data/nginx/www /data/nginx/logs /data/nginx/conf /data/etc
VOLUME ["/data"]
VOLUME ["/conf"]
WORKDIR /tmp/nginx
ADD modules3rd.ini /tmp/nginx
ADD pcre-8.32.tar.bz2 /tmp/nginx

# RUN apt install -y \
RUN apk add --no-cache --virtual=.tmp-pkgs \
    wget tar openssl-dev gcc g++ make krb5-dev git go \
	libc-dev pcre pcre-dev zlib-dev linux-headers \
	curl gnupg libxslt-dev gd-dev geoip-dev perl-dev && \
	go get -u github.com/cubicdaiya/nginx-build && $GOPATH/bin/nginx-build -d work -v 1.21.4 -with-pcre=/tmp/nginx/pcre-8.32 -zlib -openssl -m modules3rd.ini -clear && \
    cd /tmp/nginx/work/nginx/*/nginx-* && make install && cd - && \
    apk del --purge .tmp-pkgs && \
    rm -rf /tmp/*

RUN addgroup -S nginx && adduser -S nginx -G nginx

# Kerberos
# COPY ./data/etc/krb5.conf ./data/etc/krb5.keytab /etc/
COPY ./data/etc/krb5.conf /etc/
ADD ./data/nginx/nginx-start.sh /opt/nginx-start.sh

RUN chmod u=rwx /opt/nginx-start.sh && \
	chown nginx:nginx -R /opt/nginx-start.sh /data/nginx /usr/local/nginx/conf && \
	chown nginx:nginx -R /data/nginx/logs /usr/local/nginx /etc/krb5.conf /usr/local/nginx/sbin/nginx

# PORTS
EXPOSE 10080
EXPOSE 10443

USER nginx
ENTRYPOINT ["/opt/nginx-start.sh"]
