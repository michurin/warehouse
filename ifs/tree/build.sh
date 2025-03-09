#!/bin/sh -ex

./full.py >svg.svg
magick -density 400 svg.svg full.png
./rounded.py >svg.svg
magick -density 400 svg.svg rounded.png
rm svg.svg

# convert to black and white:
# magick rounded.png -alpha off -auto-threshold otsu rounded-bw.png
