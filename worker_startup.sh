#!/bin/bash
CURDIR=`pwd`
cd /tmp
wget https://bootstrap.pypa.io/get-pip.py
python get-pip.py
pip install flake8
cd ${CURDIR}
go run worker.go db.go config.go payload.go redis.go