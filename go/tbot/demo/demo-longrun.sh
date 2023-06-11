#!/bin/sh

# TODO sections, comments and use cases
# TODO mention $PATH, $CWD, etc.

# Very simple, just output to stdout
if [ -z "$1" ] # no args
then
    echo '%!PRE'
    echo 'Long-running environment:'
    env | sort
    exit
fi

# interaction

CTRL="http://localhost$tg_x_ctrl_addr"

# show how to use tg_x_to
if [ "$1" = "countdown" ]
then
    for i in 5 4 3 2 1
    do
        curl -qs "$CTRL/x?to=$tg_x_to" -d "$i..." >&2
        sleep 1
    done
    echo 'done.'
    exit
fi
