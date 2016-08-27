#!/bin/bash
##############################################################
# RPI BUILD SCRIPT
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
LDFLAGS="-w -s"
cd "${CURRENT_PATH}/.."

##############################################################
# Sanity checks

if [ ! -d ${CURRENT_PATH} ] ; then
  echo "Not found: ${CURRENT_PATH}" >&2
  exit -1
fi
if [ "${GO}" == "" ] || [ ! -x ${GO} ] ; then
  echo "go not installed or executable" >&2
  exit -1
fi

##############################################################
# go get dependencies

##############################################################
# install

echo "go install helloworld"
go install -ldflags "${LDFLAGS}" cmd/helloworld.go

echo "go install vcgencmd"
go install -ldflags "${LDFLAGS}" cmd/vcgencmd.go

echo "go install gpio"
go install -ldflags "${LDFLAGS}" cmd/gpio.go

echo "go install i2c"
go install -ldflags "${LDFLAGS}" cmd/i2c.go

