# KDC and SPNEGO enabled Nginx

## Prereqs
This has been tested in the following environments

* kdc: Ubuntu 20.04 Intel + Docker stable, kerberos client: Ubuntu 20.04 Intel
* kdc: Ubuntu 20.04 Intel + Docker stable, kerberos client: Macos 11.2 Apple Silicon
* kdc: Macos 11.2 Apple Silicon + Docker preview 3.1.0(60984), client: Macos 11.2 Apple Silicon

* Ubuntu

```bash
sudo apt install krb5-user
```

* Macos

```bash
brew install krb5
```
* Ubuntu/Macos: Create/configure /etc/krb5.conf like docker/nginx-spnego/dagta/etc/krb5.conf. Replace 'localhost' with your Docker host name (KDC)
* MacOs uses Heimdal Kerberos which needs 'tcp/' in /etc/krb5.conf as below

```bash
[realms]
        REALM.SOKOIDE.COM = {
                kdc = tcp/scottmm.local:10088
                admin_server = tcp/scottmm.local:10749
                                kpasswd_server = tcp/scottmm.local:10464
        }
```


## How to build and run

* Both KDC and Nginx containers use alpine:latest base image. If you build it on M1 mac, it'll use alpine arm64 image. Otherwise x86_64.

```bash
rm -rf tmp/krb5kdc-data/*
docker network create shared
docker-compose up --build
```

### What is configured

* By default, the following 3 users and 1 HTTP SPN are configured
* Password for the first 3 are `admin`. The last one is random (-randkey)
   * admin
   * scott
   * sandy
   * HTTP/nginx-spnego
* See docker/kdc/docker-entrypoint.sh for details
* To add more principals manually, you can do this

```bash
# logon to krb5-server
docker exec -it krb5-server /bin/sh
kadmin -p admin/admin
# default password is defined as 'admin'. see docker-compose.yml -> KRB5_PASS

# add your id (e.g. sokoide@REALM.SOKOIDE.COM)
addprinc $YOURID
# if you cant to add service id (e.g. HTTP/nginx-spnego@REALM.SOKOIDE.COM) and export the keytab
addprinc HTTP/nginx-spnego
ktadd HTTP/nginx-spnego
exit

cp /etc/krb5.keytab /var/lib/krb5kdc # /var/lib/krb5kdc is mapped to ./tmp/krb5kdc-data on Mac
exit
```

## How to test

* CLI w/ curl

```bash
sudo vim /etc/hosts
# add nginx-spnego in the line of your mac IP address
# e.g. 192.168.x.y nginx-spnego

kinit $YOURID # get your credential
curl --negotiate -u: -v http://nginx-spnego:20080/

# confirm curl fails as expected w/o --negotiate and -u:
```

* w/ Go
```bash
cd go

# change your KDC server name from 'scottmm.local' to your host in main.go

# 1. Keytab
# generate your keytab. type your id's password you used above (addpring $YOURID) when prompted
# Mac
ktutil --keytab=scott.keytab add -password -p scott -V 1 -e aes256-cts-hmac-sha1-96
# Mac verify
ktutil --keytab=scott.keytab list --keys

# Linux
ktutil
addent -password -p sokoide -v 1 -f
wkt sokoide.keytab
# Linux verify
list -e
exit # run SPNEGO go run main.go -kt -ktpath ./scott.keytab
# 2. Ccache
kinit -c hoge.ccache scott # default password is 'admin'

# run SPNEGO
go run main.go -cc -ccpath ./hoge.ccache

```

# To geenrate ccache on MacOS
kinit -c hoge.ccache scott


