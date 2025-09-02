#!/bin/sh -xe

openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj "/CN=localhost/C=GB/O=MornInc/OU=MornUnit"
