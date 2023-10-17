#!/bin/sh
# service
svc=nods
# boilerplate commands
set -x
set -e
wd=/home/osintami/${svc}
sudo mkdir -p ${wd}
sudo cp osintami-${svc}.service /lib/systemd/system/.
sudo systemctl daemon-reload
sudo service osintami-${svc} stop
git stash
git stash clear
git pull
go mod tidy
go get -u
go build -o ${svc}
sudo -u osintami cp ${svc} ${wd}/.
# copy service specific files
sudo -u osintami cp config.json ${wd}/.
sudo -u osintami cp schema/*.json ${wd}/../data/.
# restart service
sudo service osintami-${svc} start
