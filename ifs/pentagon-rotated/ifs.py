#!/usr/bin/env python3

from collections import namedtuple
from math import pi, sin, cos

# TYPES

Vector = namedtuple('Vector', ('x', 'y', 'r', 'a'))

# CONSTANTS

A = 22.38617755919675     # 0.3907125053440511 rad # f = lambda x: 2*sin(x)-sin(pi*.4-x) = 0
K = 0.40044657145607854   # 1/(2*cos(a)+cos(pi*.4-a))
B = 13.613822440803249    # atan2(sin(A*pi/180), 1/K - cos(A*pi/180))/pi*180

# DRAW

SVG_BODY='''<svg height="80%" style="background-color:black" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10">
<rect width="100%" height="100%" fill="black"/>
<g stroke="white" fill="none"
   stroke-width=".02"
   stroke-linecap="round" stroke-linejoin="round"
   transform="scale(1 -1) translate(0 -10)">
{}
</g>
</svg>'''

def color(c, *v):
    return '<g stroke="{}">\n{}\n</g>'.format(c, '\n'.join(v))

def line(x1, y1, x2, y2):
    return '<line x1="{}" y1="{}" x2="{}" y2="{}"/>'.format(x1, y1, x2, y2)

def vector(v):
    a = v.a * pi / 180.
    return line(v.x, v.y, v.x + v.r * cos(a), v.y + v.r * sin(a))

# GEOMETRY

def rotate(v, a):
    return Vector(v.x, v.y, v.r, v.a+a)

def scale(v, k):
    return Vector(v.x, v.y, v.r*k, v.a)

def put_to_end(s, p):
    a = p.a * pi / 180.
    return Vector(p.x + p.r * cos(a), p.y + p.r * sin(a), s.r, s.a)

def pentogram(out, level, base):
    if level >= len(out):
        return
    a = scale(rotate(base, A), K)
    out[level].append(a)
    pentogram(out, level+1, a)
    for _ in range(4):
        a = put_to_end(rotate(a, -72), a)
        out[level].append(a)
        pentogram(out, level+1, a)

# MAIN

def main():
    deep = 7
    v = Vector(1, 5, 12, B-A+36)
    out = list([] for _ in range(deep))
    pentogram(out, 0, v)
    print(SVG_BODY.format('\n'.join((
        # steps
        color('#008800', *map(vector, out[0])),
        color('#00ff00', *map(vector, out[1])),
        color('#ffff00', *map(vector, out[2])),
        color('#ffffff', *map(vector, out[3])),
        # full result
        #color('#ffffff', *map(vector, out[deep-1])),
        # just show initial vector
        #color('#ff0000', vector(v)),
    ))))

if __name__ == "__main__":
    main()
