# KDC and SPNEGO enabled Nginx

## How to build and run (first time)

```bash
docker-compose up --build
```

## How to configure Kerberos

* Configure HTTP/nginx-spnego@REALM.SOKOIDE.COM in KDC

```bash
# logon to krb5-server
docker exec -it krb5-server /bin/sh
ktadmin -p admin/admin
addprinc $YOURID
addprinc HTTP/nginx-spnego
ktadd HTTP/nginx-spnego
cp /etc/krb5.keytab /var/lib/krb5kdc # /var/lib/krb5kdc is mapped to /tmp/krb5kdc-data-mac on Mac

# on your docker host
cp /tmp/krb5kdc-data-mac/krb5.keytab ./nginx-spnego/data/etc

# on your mac
sudo cp /tmp/krb5kdc-data-mac/krb5.keytab ./nginx-spnego/data/etc
```

## How to run (second time or later)

```bash
docker-compose up --build

sudo vim /etc/hosts
# add nginx-spnego in the line of your mac IP address
# e.g. 192.168.x.y nginx-spnego

kinit $YOURID # get your credential
curl --negotiate -u: -v http://nginx-spnego:20080/
```
