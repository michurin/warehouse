#!/usr/bin/env python3


'''
Generate very simple IFS in SVG format
'''


from collections import namedtuple, deque


V = namedtuple('V', ('o', 'd'))


def calc(q, v, s):
    s -= 1
    if s <= 0:
        return
    q.append(v)
    calc(q, V(v.o + v.d, v.d * (.5+.5j)), s)
    calc(q, V(v.o + v.d, v.d * (.5-.5j)), s)


def xy(n, c):
    return f'x{n}="{c.real}" y{n}="{c.imag}"'


def draw(q, w, h):
    print(f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}">')
    print('<rect width="100%" height="100%" fill="#000000"/>')
    while len(q) > 0:
        v = q.popleft()
        print(f'<line {xy(1, v.o)} {xy(2, v.o+v.d)} style="stroke:rgb(255,255,255);stroke-width:{abs(v.d)/5}" />')
    print('</svg>')


def main():
    q = deque()
    calc(q, V(50+80j, -20j), 12)
    draw(q, 100, 100)


if __name__ == '__main__':
    main()
