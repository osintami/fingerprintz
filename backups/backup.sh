#!/bin/bash
set -x
set -e

wd=`pwd`

sudo service osintami-gateway stop
sudo service osintami-etlr stop
sudo service osintami-nods stop
sudo service osintami-whoami stop

cd /home/osintami
sudo rm -f /logs/*.log

now=`date +"%Y-%m-%d"`
cd /tmp
file=osintami_${now}.dump
sudo -u postgres pg_dump -Fc -f $file -d osintami
sudo mv ${file} ${wd}/.
file=osintami_${now}.dump

cd /home/osintami
tar -cvf ${wd}/osintami_${now}.tar gateway etlr nods whoami data logs
sudo service osintami-gateway start
sudo service osintami-etlr start
sudo service osintami-nods start
sudo service osintami-whoami start

cd ${wd}
gzip osintami_${now}.tar
