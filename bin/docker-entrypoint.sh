#!/bin/sh

echo "generate foundation config and start it"
echo "you have to define at least MONGOHOST and MONGODB"
echo "MONGOUSER, MONGOPASS, SECURE_COOKIE, SSL_CERT, XDOMAIN, REDIS are optional parametes"

if ! [ -n "$MONGOHOST" ] ; then
  echo "you have to define MONGOHOST environment parameter"
  exit 2
fi
if ! [ -n "$MONGODB" ] ; then
  echo "you have to define MONGODB environment parameter"
  exit 2
fi

uri="mongodb://"
if [ -n "$MONGOUSER" ] ; then
  if [ -n "$MONGOPASS" ] ; then
    uri="$uri$MONGOUSER:$MONGOPASS@"
  else
    uri="$uri$MONGOUSER@"
  fi
fi
uri="$uri$MONGOHOST/$MONGODB"

echo "media.fsmedia.folder=/home/ottemo/media" > ottemo.ini
echo "mongodb.db=$MONGODB" >> ottemo.ini
echo "mongodb.uri=$uri" >> ottemo.ini

if [ -n "$SECURE_COOKIE" ] ; then
  echo "secure_cookie=$SECURE_COOKIE" >> ottemo.ini
else
  echo "secure_cookie=false" >> ottemo.ini
fi
# ssl cert can be placed to nfs share or mounted as secret into container on kubernetes
if [ -n "$SSL_CERT" ] ; then
  echo "ssl.cert=$SSL_CERT" >> ottemo.ini
fi
if [ -n "$XDOMAIN" ] ; then
  echo "xdomain=$XDOMAIN" >> ottemo.ini
fi
if [ -n "$REDISHOST" ] ; then
  echo "redis.servers=$REDISHOST" >> ottemo.ini
fi

echo "use follow ottemo.ini config:"
cat ottemo.ini

./foundation
