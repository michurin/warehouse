#!/usr/bin/env python3

import math

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

def line(f, a, b, fo=None):
    a = f(a)
    b = f(b)
    if fo is not None:
        a = fix(fo, a)
        b = fix(fo, b)
    print(f'<line {xy(1, a)} {xy(2, b)} style="stroke:rgb(255,255,255);stroke-width:.3px;stroke-linecap:round" />')

def drawClose():
    print('</g></svg>')

# GEOMETRY TRICKS

delta = 1/(2*math.tan(math.pi/8))-1
radius = math.hypot(1+delta, .5)
factor2 = delta*delta - radius*radius

def fix(g, a):
    o = g(0)
    dr = g(1)-o
    a = a-o # hack
    a1 = a/dr
    e = abs(a1.imag/a1.real)
    qa = (1+e*e)
    qb = 2*e*delta
    qc = factor2
    qd = qb*qb - 4*qa*qc
    qds = math.sqrt(qd)
    h = (-qb+qds)/2/qa
    f = h*e
    return a * f + o

# FUNCTIONS

def link(f, g):
    return lambda x: f(g(x))

# SIMPLEX

def arch(f, fx): # relatively to bottom-center of square
    f = link(f, lambda x: x+.5j) # move to center of square; just for convenient coordinates
    line(f, -.25, -.125+.125j, fx)
    line(f, .25, .125+.125j, fx)
    line(f, -.125+.125j, .125+.125j, fx)
    line(f, -.125+.125j, -.125+.25j, fx)
    line(f, .125+.125j, .125+.25j, fx)

def archline(fo, depth, fx):
    # show square
    #line(fo, -.5, .5)
    #line(fo, -.5+1j, .5+1j)
    #line(fo, -.5, -.5+1j)
    #line(fo, .5, .5+1j)
    f = fo
    n = 0
    for _ in range(depth):
        arch(f, fx)
        for s in range(n):
           arch(lambda x: f(x+(s+1)*.5), fx)
           arch(lambda x: f(x-(s+1)*.5), fx)
        line(link(f, lambda x: x-.5*n), -.25+.5j, -.25+.25j, None) # left border line; we do not need right one
        line(link(f, lambda x: x-.5*n), -.25+.25j, 0, None)
        n = n*2+1
        f = link(f, lambda x: x*.5+.5j)
    return

def spiral(fo, n, fxo):
    f = fo
    fx = fxo
    deptho = 9
    depth = deptho
    for _ in range(n): # left
        archline(f, int(depth), fx)
        depth-=.5
        f = link(f, lambda x: x*(.5+.5j)-.25+.25j)
        fx = link(fx, lambda x: x*(.5+.5j)-.25+.25j)
    f = fo
    fx = fxo
    depth = deptho
    for _ in range(n-1): # right
        depth-=.5
        f = link(f, lambda x: x*(.5-.5j)+.25+.25j)
        fx = link(fx, lambda x: x*(.5-.5j)+.25+.25j)
        archline(f, int(depth), fx)
        line(f, 0, -.5j, None)

# MAIN

def main():
    drawOpen(100, 100)
    position = lambda x: x*40+50+30j
    #arch(position, None)
    #arch(position, position)
    #archline(position, 4, position)
    spiral(position, 14, position)
    drawClose()

if __name__ == '__main__':
    main()
