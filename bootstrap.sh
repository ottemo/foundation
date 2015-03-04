apt-get -y install bzr
apt-get -y install git
apt-get -y install golang-go


rm -rf /opt/ottemo
mkdir -pv /opt/ottemo/go/src/github.com/ottemo
mkdir -pv /opt/ottemo/go/bin
mkdir -pv /opt/ottemo/go/pkg
mkdir -pv /opt/ottemo/media

export GOPATH=/opt/ottemo/go

git clone https://ottemo-dev:freshbox111222333@github.com/ottemo/foundation.git /opt/ottemo/go/src/github.com/ottemo/foundation

cd $GOPATH/bin
echo "media.fsmedia.folder=/opt/ottemo/media" >> ottemo.ini
echo "mongodb.db=ottemo-demo" >> ottemo.ini
echo "mongodb.uri=mongodb://ottemo:ottemo2014@candidate.42.mongolayer.com:10243,candidate.43.mongolayer.com:10327/ottemo-demo" >> ottemo.ini

cd $GOPATH/src/github.com/ottemo/foundation && go get -t
cd $GOPATH/src/github.com/ottemo/foundation && go get gopkg.in/mgo.v2
cd $GOPATH/src/github.com/ottemo/foundation && go get gopkg.in/mgo.v2/bson
cd $GOPATH/src/github.com/ottemo/foundation && go build -tags mongo
