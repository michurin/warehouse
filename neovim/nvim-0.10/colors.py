#!/usr/bin/env python3

import re

def fix(x): # https://www.ditig.com/publications/256-colors-cheat-sheet
    try:
        v = int(x)
    except ValueError:
        return x
    if v < 7:
        if v & 1:
            r = 128
        else:
            r = 0
        if v & 2:
            g = 128
        else:
            g = 0
        if v & 4:
            b = 128
        else:
            b = 0
        return f'#{r:02x}{g:02x}{b:02x}'
    if v == 7:
        return '#c0c0c0'
    if v == 8:
        return '#808080'
    if v < 16:

        if v & 1:
            r = 255
        else:
            r = 0
        if v & 2:
            g = 255
        else:
            g = 0
        if v & 4:
            b = 255
        else:
            b = 0
        return f'#{r:02x}{g:02x}{b:02x}'
    if v < 232:
        cc = [0, 95, 135, 175, 215, 255]
        v -= 16
        b = cc[v % 6]
        v //= 6
        g = cc[v % 6]
        v //= 6
        r = cc[v % 6]
        return f'#{r:02x}{g:02x}{b:02x}'
    if v > 255:
        raise Exception('Invalid integer (big)')
    else:
        v -= 232
        v *= 10
        v += 8
        return f'#{v:02x}{v:02x}{v:02x}'

def build(pairs):
    for k in ('gui', 'guifg', 'guibg', 'cterm', 'ctermfg', 'ctermbg'):
        if k in pairs.keys():
            continue
        sk = {'gui': 'cterm', 'guifg': 'ctermfg', 'guibg': 'ctermbg'}.get(k)
        if sk is None:
            continue
        if sk not in pairs.keys():
            continue
        pairs[k] = fix(pairs[sk])
    return ' '.join(k+'='+pairs[k] for k in sorted(pairs.keys()))

def main():
    #print(fix('15'))
    #return
    with open('init.vim-0.9') as fh:
        for line in fh:
            line = line.rstrip()
            m = re.match(r'^\s*(highlight)\s+(\S+)((\s+[a-z]+=\S+)+)$', line)
            if m:
                name = m.group(2)
                pairs = dict(x.split('=') for x in m.group(3).split())
                line = m.group(1) + ' ' + name + ' ' + build(pairs)
            print(line)

if __name__ == '__main__':
    main()
