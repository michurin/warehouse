#!/usr/bin/env python3

# LOW LEVEL DRAWING

def drawOpen(w, h):
    w *= 10
    h *= 10
    print('<?xml version="1.0" encoding="UTF-8"?>')
    print('<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">')
    print(f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}">')
    print('<rect width="100%" height="100%" fill="#000000"/>')
    print('<g transform="scale(10, -10) translate(0, -100)">')

def xy(n, c):
    return f'x{n}="{c.real}" y{n}="{c.imag}"'

def line(f, a, b):
    print(f'<line {xy(1, f(a))} {xy(2, f(b))} style="stroke:rgb(255,255,255);stroke-width:.3px;stroke-linecap:round" />')

def drawClose():
    print('</g></svg>')

# FUNCTIONS

def link(f, g):
    return lambda x: f(g(x))

# SIMPLEX

def arch(f): # relatively to bottom-center of square
    f = link(f, lambda x: x+.5j) # move to center of square; just for convenient coordinates
    line(f, -.25, -.125+.125j)
    line(f, .25, .125+.125j)
    line(f, -.125+.125j, .125+.125j)
    line(f, -.125+.125j, -.125+.25j)
    line(f, .125+.125j, .125+.25j)

def archline(fo, depth):
    # show square
    #line(fo, -.5, .5)
    #line(fo, -.5+1j, .5+1j)
    #line(fo, -.5, -.5+1j)
    #line(fo, .5, .5+1j)
    f = fo
    n = 0
    for _ in range(depth):
        arch(f)
        for s in range(n):
           arch(lambda x: f(x+(s+1)*.5))
           arch(lambda x: f(x-(s+1)*.5))
        line(link(f, lambda x: x-.5*n), -.25+.5j, -.25+.25j) # left border line; we do not need right one
        line(link(f, lambda x: x-.5*n), -.25+.25j, 0)
        n = n*2+1
        f = link(f, lambda x: x*.5+.5j)
    return

def spiral(fo):
    f = fo
    deptho = 9
    n = 14
    depth = deptho
    for _ in range(n): # left
        archline(f, int(depth))
        depth-=.5
        f = link(f, lambda x: x*(.5+.5j)-.25+.25j)
    f = fo
    depth = deptho
    for _ in range(n-1): # right
        depth-=.5
        f = link(f, lambda x: x*(.5-.5j)+.25+.25j)
        archline(f, int(depth))
        line(f, 0, -.5j) # right most line

# MAIN

def main():
    drawOpen(100, 100)
    #arch(lambda x: x*80+50+10j)
    #archline(lambda x: x*80+50+10j)
    spiral(lambda x: x*40+50+50j)
    drawClose()

if __name__ == '__main__':
    main()
