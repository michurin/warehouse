#!/usr/bin/env python3


"""
However, we use mirror approach like in tiling.py, but do it more careful:
    - We use true mirror transformation (z/|z|² instead if 1/z)
    - We use unit circle
Background:
    - Low of sin in form: R=[BC]=sin(A)/sin(A+B), r=[AC]=sin(B)/sin(A+B), assuming [AB]=1, A=π/7, B=π/2-π/7-π/3
"""

import itertools

import matplotlib.pyplot as plt

import numpy as np


def flat_shape(x1, x2, y1, y2, steps_per_unit):
    return (int((y2-y1)*steps_per_unit + 1), int((x2-x1)*steps_per_unit + 1))


def flat(x1, x2, y1, y2, steps_per_unit):
    sy, sx = flat_shape(x1, x2, y1, y2, steps_per_unit)
    a = np.linspace(x1, x2, sx)[np.newaxis, ...]
    b = np.linspace(y2, y1, sy)[..., np.newaxis] * 1j
    return (a + b).astype(np.cdouble)


def radiuses(n, m):  # n-sided polygons, m polygons in each corner
    a = np.pi/n
    b = np.pi/2 - a - np.pi/m
    d = np.sin(a + b)
    return (
            np.sin(b)/d,  # radius of seed
            np.sin(a)/d,  # radius of mirror
            np.exp(np.linspace(0, 2*np.pi, n, False) * 1j))  # centers of merrors


def main_simplest():
    """
    The most naive implementation
    Consider `for level in...: for ... in product(repeat(x, level))` as a breadth-first search
    """
    r, rm, cm = radiuses(3, 7)
    size_args = -1, 1, -1, 1, 200
    result = np.zeros(flat_shape(*size_args), dtype=np.cdouble)

    for level in range(5):
        for seq in itertools.product(*itertools.repeat(cm, level)):
            print(f'level={level}, seq={seq}')
            z = flat(*size_args)
            for c in seq:
                z -= c
                z = rm * rm * z / np.power(np.abs(z), 2) + c
            m = np.logical_and(np.abs(z) <= r, np.all(np.abs(z[..., np.newaxis] - cm[np.newaxis, np.newaxis, ...]) >= rm, axis=2))  # map: we are inside central seed
            result[m] = z[m]

    plt.imshow(np.abs(result), extent=size_args[:4])
    plt.colorbar()
    plt.title(r'Poincare')
    plt.show()


def main():
    """
    TODO: is to be optimized
    """
    # settings
    schl = 7, 3  # Schläfli symbol. Try (7, 4), (3, 7)...
    levels = 5

    # params in assuption that all mirrors centers are on unit circle
    r, rm, cm = radiuses(*schl)
    rm2 = rm * rm
    rw = np.sqrt(1 - rm2)  # world's radius

    # scale to make universe radius 1
    r /= rw
    rm /= rw
    cm /= rw
    rw = 1  # rw/rw
    rm2 = rm * rm

    size_args = -1, 1, -1, 1, 100
    result = np.zeros(flat_shape(*size_args), dtype=np.cdouble)
    result_mask = np.abs(flat(*size_args)) < rw

    for level in range(levels):
        for seq in itertools.product(*itertools.repeat(cm, level)):
            print(f'level={level}, seq={seq}')
            z = flat(*size_args)
            for c in seq:
                z -= c
                a = np.abs(z)
                t = np.logical_and(result_mask, a < rm)
                z[t] = rm2 * z[t] / np.power(a[t], 2) + c
            m = np.logical_and(result_mask, np.logical_and(np.abs(z) <= r, np.all(np.abs(z[..., np.newaxis] - cm[np.newaxis, np.newaxis, ...]) >= rm, axis=2)))  # map: we are inside central seed
            result[m] = z[m]
            result_mask[m] = False

    plt.imshow(np.abs(result), extent=size_args[:4])
    plt.colorbar()
    plt.title(f'Schläfli symbol {{{", ".join(map(str, schl))}}}')
    plt.show()


if __name__ == '__main__':
    main()
