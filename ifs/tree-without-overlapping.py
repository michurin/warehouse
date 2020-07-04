#!/usr/bin/env python3


import math
from collections import namedtuple, deque


V = namedtuple('V', ('o', 'd'))  # normal
B = namedtuple('B', ('o', 'd'))  # border line
L = namedtuple('L', ('o', 'd'))  # square (link)


L45 = .5+.5j
R45 = .5-.5j


def top(q, v):
    q.append(v)


def arm(q, v, p, s):
    for _ in range(s):
        v = V(v.o + v.d, v.d * [R45, L45][p])
        p = (p+1) % 2
        q.append(v)


def border(q, v):
    q.append(B(v.o + v.d*(3-1j), v.d*2j))


def fill(q, v, n):
    so = v.o + v.d
    m = v.d
    c = 1
    for _ in range(n//2):
        s = so
        for k in range(c):
            q.append(V(s + m * (1+.5j), m * (.25-.25j)))
            q.append(V(s + m * (1-.5j), m * (.25+.25j)))
            q.append(L(s + m * (1.25-.25j), m * .25j))  # link
            q.append(V(s + m * (1.25-.25j), m * .25))
            q.append(V(s + m * (1.25+.25j), m * .25))
            s += m*1j
        c = 2*c + 1
        so += m * (1-.5j)
        m *= .5


def xy(n, c):
    return f'x{n}="{c.real}" y{n}="{c.imag}"'


def coords(v, *p):
    return ' '.join(f'{c.real},{c.imag}' for c in (v.o + v.d * z for z in p))


def colors_dark_example(m):
    ca = 70 - int(math.sqrt(m)*50)
    bg = f'rgb({ca}, {ca}, {ca})'
    return bg, '#555555'


def colors(m):
    if m < .1:
        return '#ff0000', '#ff0000'
    if m < .35:
        return '#990000', '#ff0000'
    return '#555555', '#777777'


def draw(q, w, h, size_mag):
    print(f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}">')
    print('<rect width="100%" height="100%" fill="#000000"/>')
    while len(q) > 0:
        v = q.popleft()
        if type(v) is B:
            _, fg = colors(0)
            sw = .02
            print(f'<line {xy(1, v.o)} {xy(2, v.o+v.d)} stroke="{fg}" stroke-linecap="round" stroke-linejoin="round" stroke-width="{sw}" />')
        else:
            m = abs(v.d)/size_mag
            bg, fg = colors(m)  # see colors_dark_example
            sw = max(min(m, .2), .02)
            if type(v) is L:
                pp = (.5+.5j, 1.5+.5j, 1.5-.5j, .5-.5j)
            else:
                pp = (.5+.5j, 1+.5j, 1-.5j, .5-.5j)
            print(f'<polygon points="{coords(v, *pp)}" fill="{bg}" stroke="{fg}" stroke-linecap="round" stroke-linejoin="round" stroke-width="{sw}" />')
    print('</svg>')


def main():
    q = deque()
    v = V(50+80j, -20j)
    n = 22

    top(q, v)
    arm(q, v, 0, n-1)
    arm(q, v, 1, n-1)
    fill(q, v, n)
    border(q, v)
    vr = v
    vl = v
    for t in range(n-2, 0, -1):
        vr = V(vr.o + vr.d, vr.d * R45)
        vl = V(vl.o + vl.d, vl.d * L45)
        fill(q, vl, t)
        fill(q, vr, t)
        arm(q, vr, 0, t)
        arm(q, vl, 1, t)
        border(q, vr)
        border(q, vl)

    draw(q, 100, 100, 20)


if __name__ == '__main__':
    main()
