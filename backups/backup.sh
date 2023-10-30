#!/bin/bash
set -x
set -e

wd=`pwd`

service osintami-gateway stop
service osintami-etlr stop
service osintami-nods stop
service osintami-whoami stop

cd /home/osintami
rm -f /logs/*.log

now=`date +"%Y-%m-%d"`
cd /tmp
file=osintami_${now}.dump
sudo -u postgres pg_dump -Fc -f $file -d osintami
mv ${file} ${wd}/.
chown root:root ${file}
file=osintami_${now}.dump

cd /home/osintami
tar -cvf ${wd}/osintami_${now}.tar gateway etlr nods whoami data logs
service osintami-gateway start
service osintami-etlr start
service osintami-nods start
service osintami-whoami start

cd ${wd}
gzip osintami_${now}.tar
