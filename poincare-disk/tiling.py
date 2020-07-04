#!/usr/bin/env python3


from collections import deque
import numpy as np


def coords(*p):
    return ' '.join(f'{c.real},{c.imag}' for c in p)


def draw(q, w, h, scale):
    print(f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}">')
    # background? print('<rect width="100%" height="100%" fill="#000000"/>')
    for v in q:
        v = v * scale + w/2 + h/2*1j
        print(f'<polygon points="{coords(*v)}" fill="#000000" stroke-width="0"/>')
    print('</svg>')


def main():
    levels = 5
    vertex = 7  # obviously, have to be 7 or more

    mc = 1/2/np.sin(np.pi/vertex)  # mirror centers radius
    r0 = mc - 1  # radius of the first ("biggest") seed
    angles = np.linspace(0, 2*np.pi, vertex, False)

    queue = deque()

    seeds = np.array([np.exp(angles * 1j) * r0])
    queue.extend(seeds)
    for n in range(levels):
        next_generation = deque()
        for a in angles:  # in fact, we could add axis instead this loop, houwever, I wan to keep only needed reflections see inside
            reflections = np.exp(a * 1j) * (1 / (seeds - mc) + mc)
            d = np.abs(reflections - mc * np.exp(a * 1j))
            mask = np.sum(d < 1, axis=1) >= vertex - 1  # we take only poligons inside mirror circle (one vertex is allowed to touch border) it is oversimlified approach, it vanish only major overlapping
            reflections = reflections[mask]
            queue.extend(reflections)
            next_generation.extend(reflections)
            # if n == 1: break  # uncomment it to make sence why overloppings are not vanished completely

        seeds = np.array(next_generation)

    draw(queue, 100, 100, 100/1.15)


if __name__ == '__main__':
    main()
