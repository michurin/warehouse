#!/bin/sh
set -ex
GO=go
$GO run .
MXXX_STACK=4 $GO run .
MXXX_STACK=40 MXXX_STDERR=1 $GO test . -v
