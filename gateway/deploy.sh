#!/bin/sh
# service
svc=gateway
# boilerplate commands
set -x
set -e
wd=/home/osintami/${svc}
sudo mkdir -p ${wd}
sudo cp osintami-${svc}.service /lib/systemd/system/.
sudo systemctl daemon-reload
sudo service osintami-${svc} stop
sudo rm -f /home/osintami/logs/${svc}.log
git stash
git stash clear
git pull
go mod tidy
go get -u
go build -o ${svc}
sudo -u osintami cp ${svc} ${wd}/.
# copy service specific files
sudo setcap CAP_NET_BIND_SERVICE=+eip ${wd}/${svc}
sudo -u osintami cp config.json ${wd}/.
sudo -u osintami cp welcome.template ${wd}/.
# restart service
sudo service osintami-${svc} start
