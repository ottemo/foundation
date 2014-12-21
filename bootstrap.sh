apt-get -y install bzr
apt-get -y install git
apt-get -y install golang-go

mkdir -pv /opt/go/src/github.com/ottemo
mkdir -pv /opt/go/bin
mkdir -pv /opt/go/pkg
mkdir -pv /opt/media

export GOPATH=/opt/go

git clone https://ottemo-dev:freshbox111222333@github.com/ottemo/foundation.git /opt/go/src/github.com/ottemo/foundation

cd $GOPATH/src/github.com/ottemo/foundation
echo "db.sqlite3.uri=ottemo.db" >> ottemo.ini
echo "media.fsmedia.folder=/opt/media" >> ottemo.ini

cd $GOPATH/src/github.com/ottemo/foundation && go get -t 
cd $GOPATH/src/github.com/ottemo/foundation && go build && go install
