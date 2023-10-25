#!/bin/sh
set -x
user=`whoami`
set -e
go mod tidy
go build -o whoami
sudo ./whoami
