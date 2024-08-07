#!/bin/bash

# just check correctness of call
if test "$0" = "${0%_debug.sh}"
then
    echo "To use wrapper you are to create symlink to your original."
    echo "To wrap demo_bot.sh, lease do"
    echo "ln -s wrapper.sh demo_bot_debug.sh"
    exit 1
fi

# put your command here
cmd="${0%_debug.sh}.sh"

# tune naming for your taste
t="${cmd##*/}"
t="${t%%.*}"
base="${0%/*}/logs/$(date +%s)_${t}_${$}_"
ext='.log'

n=0
for a in "$@"
do
    echo "$a" >"${base}arg_${n}${ext}"
    n="$(($n+1))"
done

env | sort >"${base}env${ext}"

"$cmd" "$@" 2>"${base}err${ext}" | tee "${base}out${ext}"

echo "$?" >"${base}status${ext}"
