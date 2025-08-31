#!/bin/sh

# You may want to add to PATH something like this
# /usr/local/texlive/2020/bin/x86_64-darwin
#
# Preparing Arch Linux:
# pacman -Suy texlive-latex
# pacman -Suy texlive-latexrecommended
# pacman -Suy texlive-latexextra
# pacman -Suy texlive-fontsextra

pdflatex michurin.tex
