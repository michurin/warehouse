#!/bin/sh

i=1
for t in Bold Faint Italic Underline BlinkSlow BlinkRapid ReverseVideo Concealed CrossedOut
do
    printf "\033[${i}m%-12s\033[0m %s %s\n" "$t" "E[${i}m" "$t"
    i=$((i+1))
done

a=0
for x in Black Red Green Yellow Blue Magenta Cyan White
do
    b=0
    for y in '' Hi
    do
        c=30
        for z in Fg Bg
        do
            q=$((a+b+c))
            printf "\033[${q}m%-12s\033[0m %5s  " "$z$y$x" "E[${q}m"
            c=$((c+10))
        done
        b=$((b+60))
    done
    a=$((a+1))
    printf "\n"
done
