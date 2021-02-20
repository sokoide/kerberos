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


## How to build and run (first time)

```bash
docker network create shared
docker-compose up --build
```

## How to configure Kerberos

* Configure HTTP/nginx-spnego@REALM.SOKOIDE.COM in KDC

```bash
# logon to krb5-server
docker exec -it krb5-server /bin/sh
kadmin -p admin/admin
# default password is defined as 'admin'. see docker-compose.yml -> KRB5_PASS

# add your id (e.g. sokoide@REALM.SOKOIDE.COM)
addprinc $YOURID
# add service id (e.g. HTTP/nginx-spnego@REALM.SOKOIDE.COM)
addprinc HTTP/nginx-spnego
ktadd HTTP/nginx-spnego
exit

# copy docker container's /etc/krb5.keytab to the docker host
cp /etc/krb5.keytab /var/lib/krb5kdc # /var/lib/krb5kdc is mapped to ./tmp/krb5kdc-data on Mac
exit

# on your docker host (mac, linux or wsl)
sudo cp ./tmp/krb5kdc-data/krb5.keytab ./docker/nginx-spnego/data/etc
```

## How to run (second time or later)

```bash
docker-compose up
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
# generate your keytab. type your id's password you used above (addpring $YOURID) when prompted
# Mac
ktutil --keytab=sokoide.keytab add -password -p sokoide -V 1 -e aes256-cts-hmac-sha1-96
# Mac verify
ktutil --keytab=sokoide.keytab list --keys

# Linux
ktutil
addent -password -p sokoide -v 1 -f
wkt sokoide.keytab
# Linux verify
list -e
exit

# change your KDC server name from 'scottmm.local' or 'timemachine' to your host in main.go

# run SPNEGO
go run main.go
```
