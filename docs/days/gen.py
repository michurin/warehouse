#!/usr/bin/env python3

'''
Why this way: to have one style per day with clear semantic. And to be able to change CSS style without touching names and semantics.
'''

FG=(
    ('', 'inherit'),
    ('holiday', '#900'),
    ('vacation', '#090'),
    ('special', '#f0f'),
    ('dayoff', '#555'),
)

BG=(
    ('', 'inherit'),
    ('blue', '#A3D9FF'),
    ('gray', '#7E6B8F'),
    ('green', '#96E6B3'),
    ('red', '#DA3E52'),
    ('yellow', '#F2E94E'),
)

def css_class(a, b):
    if a == '':
        if b == '':
            return ''
        return b
    if b == '':
        return a
    return a + '-' + b

for bg, bg_color in BG:
    for fg, fg_color in FG:
        cls = css_class(fg, bg)
        if cls == '':
            continue
        print('#body div.{0} {{color: {1};}}'.format(cls, fg_color))
        print('#body div.{0} div {{background-color: {1};}}'.format(cls, bg_color))

print('const classChain = {')
for bg, bg_color in BG:
    prev_cls = css_class(FG[-1][0], bg)
    for fg, fg_color in FG:
        cls = css_class(fg, bg)
        print("  '{0}': '{1}',".format(prev_cls, cls))
        prev_cls = cls
print('};')

print('const bgChain = {')
for fg, fg_color in FG:
    prev_cls = css_class(fg, BG[-1][0])
    for bg, bg_color in BG:
        cls = css_class(fg, bg)
        print("  '{0}': '{1}',".format(prev_cls, cls))
        prev_cls = cls
print('};')
