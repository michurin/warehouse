#!/bin/sh

echo OKOKOK
env
echo 0--------------------------
echo "1=$1"
echo 0--------------------------
echo 'ifconfig-push 10.8.0.8 10.8.0.9' >"$1"
