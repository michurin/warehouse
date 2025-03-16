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

def line(f, a, b, color=None):
    if color is None:
        color = '#ffffff'
    print(f'<line {xy(1, f(a))} {xy(2, f(b))} style="stroke:{color};stroke-width:.2px;stroke-linecap:round" />')

def drawClose():
    print('</g></svg>')

green = 'rgba(63, 255, 127, .25)'

# FUNCTIONS

def link(f, g):
    return lambda x: f(g(x))

# SIMPLEX

def arch(f): # relatively to bottom-center of square
    line(f, -.2+.7j, -.1+.6j)
    line(f, -.2+.7j, -.15+.75j)
    line(f, -.05+.65j, -.15+.75j)
    line(f, -.05+.65j, -.1+.6j)
    line(f, -.15+.75j, -.15+.8j)
    line(f, -.05+.8j, -.15+.8j)
    line(f, -.05+.8j, -.05+.65j)
    line(f, .2+.7j, .1+.6j)
    line(f, .2+.7j, .15+.75j)
    line(f, .05+.65j, .15+.75j)
    line(f, .05+.65j, .1+.6j)
    line(f, .15+.75j, .15+.8j)
    line(f, .05+.8j, .15+.8j)
    line(f, .05+.8j, .05+.65j)
    line(f, -.05+.65j, .05+.65j)
    line(f, -.15+.75j, .15+.75j)

def sector(fo, depth):
    # show square
    line(fo, -.5, .5, green)
    line(fo, -.5+1j, .5+1j, green)
    line(fo, -.5, -.5+1j, green)
    line(fo, .5, .5+1j, green)
    f = fo
    n = 0
    for _ in range(depth):
        arch(f)
        for s in range(n):
           arch(lambda x: f(x+(s+1)*.4))
           arch(lambda x: f(x-(s+1)*.4))
        fx = link(f, lambda x: x-n*.4) # left helper
        line(fx, -.3+.5j, -.1+.5j, '#ffff00')
        line(fx, -.3+.6j, -.1+.6j, '#ffff00')
        line(fx, -.3+.5j, -.3+.6j, '#ffff00')
        line(fx, -.1+.5j, -.1+.6j, '#ffff00')
        fx = link(f, lambda x: x+n*.4) # right helper
        line(fx, .3+.5j, .1+.5j, '#ff00ff')
        line(fx, .3+.6j, .1+.6j, '#ff00ff')
        line(fx, .3+.5j, .3+.6j, '#ff00ff')
        line(fx, .1+.5j, .1+.6j, '#ff00ff')
        n = n*2+1
        f = link(f, lambda x: x*.5+.5j)
    line(fo, -.2, .2, '#ff0000') # sector's root
    line(fo, -.2+.2j, .2+.2j, '#ff0000')
    line(fo, -.2, -.2+.2j, '#ff0000')
    line(fo, .2, .2+.2j, '#ff0000')
    return

def spiral(fo, n, depth):
    f = fo
    dp = depth
    for _ in range(n): # left
        sector(f, int(dp))
        dp-=.5
        f = link(f, lambda x: x*(.5+.5j)-.1+.3j)
    f = fo
    dp = depth
    for _ in range(n-1): # right
        dp-=.5
        f = link(f, lambda x: x*(.5-.5j)+.1+.3j)
        sector(f, int(dp))

# MAIN

def main():
    drawOpen(100, 100)
    pos = lambda x: x*60+50+10j
    # arch(pos) # show one arch
    # sector(pos, 4) # show one secrot
    spiral(pos, 16, 7)
    drawClose()

if __name__ == '__main__':
    main()
