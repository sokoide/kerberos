FROM alpine:latest

MAINTAINER "So Koide" <sokoide@gmail.com>

WORKDIR /tmp

# Install prerequisites for Nginx compile
RUN apk add --no-cache tzdata krb5
ENV TZ=Asia/Tokyo \
	GOPATH=/go \
	PATH=$GOPATH/bin:$PATH

# RUN apt install -y \
RUN apk add --no-cache \
    wget tar openssl-dev gcc g++ make krb5-dev git go \
	libc-dev pcre-dev zlib-dev linux-headers \
	curl gnupg libxslt-dev gd-dev geoip-dev perl-dev

# DATA VOLUME
RUN mkdir -p /go /data/nginx/www /data/nginx/logs /data/nginx/conf /data/etc
VOLUME ["/data"]

# Build Nginx
WORKDIR /tmp/nginx
ADD modules3rd.ini /tmp/nginx
RUN go get -u github.com/cubicdaiya/nginx-build;cd /tmp/nginx; $GOPATH/bin/nginx-build -d work -v 1.19.7 -pcre -zlib -openssl -m modules3rd.ini -clear
RUN cd /tmp/nginx/work/nginx/*/nginx-*;make install;cd -

# Cleanup after Nginx build
RUN apk del --purge \
        wget \
        gcc \
        g++ \
        make \
        linux-headers \
        git \
		go && \
    rm -rf /tmp/*

RUN addgroup -S nginx && adduser -S nginx -G nginx

# Kerberos
# ADD ./data/etc/krb5.conf /etc/krb5.conf
# ADD ./data/etc/krb5.keytab /etc/krb5.keytab
# ADD ./data/nginx/nginx-start.sh /opt/nginx-start.sh
COPY ./data/etc/krb5.conf ./data/etc/krb5.keytab /etc/
ADD ./data/nginx/nginx-start.sh /opt/nginx-start.sh

RUN chmod u=rwx /opt/nginx-start.sh && \
	chown nginx:nginx -R /opt/nginx-start.sh /data/nginx /usr/local/nginx/conf && \
	chown nginx:nginx -R /data/nginx/logs /usr/local/nginx /etc/krb5.conf /etc/krb5.keytab /usr/local/nginx/sbin/nginx

# PORTS
EXPOSE 10080
EXPOSE 10443

USER nginx
ENTRYPOINT ["/opt/nginx-start.sh"]
