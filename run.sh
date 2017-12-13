#!/bin/bash

COMMIT=$1
TOKEN=$2
FULLNAME=$3/$4
DIR=`pwd`

echo "--> git clone https://github.com/${FULLNAME}.git"
git clone https://${TOKEN}:${TOKEN}@github.com/${FULLNAME}.git /tmp/${COMMIT}
cd /tmp/${COMMIT}
git checkout ${COMMIT}
echo "--> run flake8 ."

# dir
# pwd
# flake8 .

flake8 . &> report_${COMMIT}

cd ${DIR}
cp /tmp/${COMMIT}/report_${COMMIT} /data/reports/${COMMIT}

rm -rf /tmp/${COMMIT}