#!/bin/sh -ex

pandoc --standalone --embed-resources --metadata title="$(date '+ENGLISH %Y/%m/%d %H%P')" --number-sections -c style.css --toc ENGLISH.md >index.html
pandoc --standalone --embed-resources --metadata title="$(date '+ENGLISH %Y/%m/%d %H%P')" --number-sections -c style.css ENGLISH.md >index-wo-toc.html

echo python -m http.server 8112
