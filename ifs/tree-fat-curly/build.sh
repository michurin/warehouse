#!/bin/sh -ex

DPI=100

for x in tree-idea tree
do
  ./$x.py >$x.svg
  magick -density $DPI $x.svg -alpha off $x.png
done
