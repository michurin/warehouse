#!/usr/bin/env python3

import math

sq2 = math.sqrt(2)

# LOW LEVEL DRAWING

mincolor = 1000
maxcolor = -1000

def drawOpen(w, h):
    w *= 10
    h *= 10
    print('<?xml version="1.0" encoding="UTF-8"?>')
    print('<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">')
    print(f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}">')
    print('<rect width="100%" height="100%" fill="#000000"/>')
    print(f'<g transform="scale(10, -10) translate(0, -65)">') # oh, dirty hack

def xy(n, c):
    return f'x{n}="{c.real}" y{n}="{c.imag}"'

def line(f, a, b, color=None, width=None):
    global maxcolor, mincolor
    if color is None:
        color = '#ffffff'
    if width is None:
        raise NotImplemented
    if width < .02:
        width = .02
    cf = -math.log2(width)
    if cf > maxcolor:
        maxcolor = cf
    if cf < mincolor:
        mincolor = cf
    if cf < 0:
        cf = 0
    c = 55 + 200*(cf/6)
    if c > 255:
        c = 255
    color = f'rgb({c}, {c}, {c})'
    print(f'<line {xy(1, f(a))} {xy(2, f(b))} style="stroke:{color};stroke-width:{width}px;stroke-linecap:round" />')

def drawClose():
    print('</g></svg>')
    print(f'<!-- {mincolor} {maxcolor} -->')

green = 'rgba(63, 255, 127, .25)'

# FUNCTIONS

def link(f, g):
    return lambda x: f(g(x))

# SIMPLEX

def root(f, w):
    line(f, -.2, .2, width=w)
    line(f, -.2+.2j, .2+.2j, width=w)
    line(f, -.2, -.2+.2j, width=w)
    line(f, .2, .2+.2j, width=w)

def arch(f, n, w): # relatively to bottom-center of square
    w2 = w / sq2
    line(f, -.2+.7j, -.1+.6j, width=w)
    line(f, -.2+.7j, -.15+.75j, width=w)
    line(f, -.05+.65j, -.15+.75j, width=w)
    line(f, -.05+.65j, -.1+.6j, width=w)
    line(f, -.15+.75j, -.15+.8j, width=w2)
    line(f, -.05+.8j, -.15+.8j, width=w2)
    line(f, -.05+.8j, -.05+.75j, width=w2)
    line(f, .2+.7j, .1+.6j, width=w)
    line(f, .2+.7j, .15+.75j, width=w)
    line(f, .09+.69j, .15+.75j, width=w)
    line(f, .09+.61j, .1+.6j, width=w)
    line(f, .15+.75j, .15+.8j, width=w2)
    line(f, .05+.8j, .15+.8j, width=w2)
    line(f, .05+.8j, .05+.77j, width=w2)
    line(f, .07+.75j, .15+.75j, width=w2)
    line(f, -.15+.75j, -.05+.75j, width=w2)
    fx = link(f, lambda x: x*-.25j-.05+.7j) # put root to link
    w3=w2
    for i in range(n):
        root(fx, w3)
        if i == 0:
            line(fx, .4j, -.2+.2j, width=w3/sq2)
            line(fx, .4j, -.1+.5j, width=w3/sq2)
            line(fx, -.2+.2j, -.25+.25j, width=w3/sq2)
            line(fx, -.27+.33j, -.1+.5j, width=w3/sq2)
        else:
            root(link(fx, lambda x: x*(.5+.5j)-.1+.3j), w3)
        fx = link(fx, lambda x: x*(.5-.5j)+.1+.3j)
        w3 /= sq2
    #spiral(fx, 2, 1) # INFINITE RECURSION

def sector(fo, depth, width):
    # show square
    # line(fo, -.5, .5, green)
    # line(fo, -.5+1j, .5+1j, green)
    # line(fo, -.5, -.5+1j, green)
    # line(fo, .5, .5+1j, green)
    f = fo
    n = 0
    dpth = 16
    w = width
    for _ in range(depth):
        arch(f, int(dpth), w)
        for s in range(n):
           arch(lambda x: f(x+(s+1)*.4), int(dpth), w)
           arch(lambda x: f(x-(s+1)*.4), int(dpth), w)
        fx = link(f, lambda x: x-n*.4)
        line(fx, -.3+.5j, -.1+.5j, '#ffff00', width=w*sq2)
        line(fx, -.3+.6j, -.1+.6j, '#ffff00', width=w*sq2)
        line(fx, -.3+.5j, -.3+.6j, '#ffff00', width=w*sq2)
        line(fx, -.1+.5j, -.1+.6j, '#ffff00', width=w*sq2)
        fx = link(f, lambda x: x+n*.4)
        line(fx, .3+.5j, .1+.5j, '#ff00ff', width=w*sq2)
        line(fx, .3+.6j, .1+.6j, '#ff00ff', width=w*sq2)
        line(fx, .3+.5j, .3+.6j, '#ff00ff', width=w*sq2)
        line(fx, .1+.5j, .1+.6j, '#ff00ff', width=w*sq2)
        w /= sq2
        n = n*2+1
        f = link(f, lambda x: x*.5+.5j)
        dpth -= .5
    line(fo, -.2, .2, '#ff0000', width=width*sq2*2)
    line(fo, -.2+.2j, .2+.2j, '#ff0000', width=width*sq2*2)
    line(fo, -.2, -.2+.2j, '#ff0000', width=width*sq2*2)
    line(fo, .2, .2+.2j, '#ff0000', width=width*sq2*2)
    return

def spiral(fo, n, depth, width):
    f = fo
    dp = depth
    w = width
    for _ in range(n): # left
        sector(f, int(dp), w)
        dp-=.5
        w /= sq2
        f = link(f, lambda x: x*(.5+.5j)-.1+.3j)
    f = fo
    dp = depth
    w = width
    for _ in range(n-1): # right
        dp-=.5
        w /= sq2
        f = link(f, lambda x: x*(.5-.5j)+.1+.3j)
        sector(f, int(dp), w)

# MAIN

def main():
    drawOpen(100, 65)
    pos = lambda x: x*60+50+1j
    # arch(pos, 4, 1) # show one arch
    # sector(pos, 4, .5) # show one secrot
    # spiral(pos, 10, 5, .5)
    spiral(pos, 20, 8, .5)
    drawClose()

if __name__ == '__main__':
    main()
