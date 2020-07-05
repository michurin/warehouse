#!/usr/bin/env python3


import matplotlib.pyplot as plt

import numpy as np


def flat(x1, x2, y1, y2, steps_per_unit):
    a = np.linspace(x1, x2, int((x2-x1)*steps_per_unit + 1))[np.newaxis, ...]
    b = np.linspace(y2, y1, int((y2-y1)*steps_per_unit + 1))[..., np.newaxis] * 1j
    return a + b


def main():
    z = flat(-4, 4, -4, 4, 200)

    modX = 1  # tilling rects size
    modY = 1
    N = 5  # number of A-spirals
    M = 9  # number of B-spirals
    F = np.pi * 2j / (M*modY*1j+N*modX)

    with np.errstate(invalid='ignore', divide='ignore'):
        f = np.log((z+1)/(z-1)) / F

    np.nan_to_num(f, copy=False)  # f[z-1 == 0] = 0; f[z+1 == 0] = 0

    x = np.real(f) % modX
    y = np.imag(f) % modY
    data = np.logical_xor(np.logical_and(x > .25, x < .75), np.logical_and(y > .25, y < .75)).astype(int)

    plt.imshow(data, extent=(-4, 4, -4, 4))
    plt.colorbar().set_label('val')
    plt.title(r'Möbius transformation: $\frac{\ln(M)}{F}$; $M=\frac{z+1}{z-1}$')
    plt.show()


if __name__ == '__main__':
    main()
