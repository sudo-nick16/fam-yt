#!/bin/bash
#
HOST=$1
if [ -z "$HOST" ]; then
	echo "Usage: $0 <hostname>"
	HOST="localhost"
fi
if [ -z "$GOROOT" ]; then
	echo "GOROOT is not set"
	GOROOT="/usr/local/go"
fi
go run $GOROOT/src/crypto/tls/generate_cert.go --host $HOST
