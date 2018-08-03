#!/bin/bash

COMMIT=$1
TOKEN=$2
FULLNAME=$3/$4
REPORTPATH=$5
DIR=`pwd`

# fetch git repo
git clone https://${TOKEN}:${TOKEN}@github.com/${FULLNAME}.git /tmp/${COMMIT}
cd /tmp/${COMMIT}
git checkout ${COMMIT}

# run flake8
flake8 . &> report_output.txt.txt

# save persistent
cd $DIR
cp /tmp/${COMMIT}/report_output.txt.txt ${REPORTPATH}/${COMMIT}

# clean up
rm -rf /tmp/${COMMIT}
