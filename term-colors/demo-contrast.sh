#!/bin/sh

for p in 3 9
do
    for a in 0 1 2 3 4 5 6 7
    do
        for q in 4 10
        do
            for b in 0 1 2 3 4 5 6 7
            do
                printf " \033[${p}${a};${q}${b}m ${p}${a}/${q}${b} \033[m"
            done
        done
        printf "\n"
    done
done
