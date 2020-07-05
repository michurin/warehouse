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
    modY = np.sqrt(3)
    N = 4  # number of A-spirals
    M = 9  # number of B-spirals
    F = np.pi * 2j / (M*modY*1j+N*modX)

    with np.errstate(invalid='ignore', divide='ignore'):
        f = np.log((z+1)/(z-1)) / F

    np.nan_to_num(f, copy=False)  # f[z-1 == 0] = 0; f[z+1 == 0] = 0

    f += modX / 2  # small shift for symetry

    t = np.real(f) % modX + 1j * (np.imag(f) % modY)  # (x, y) inside tile

    cc = np.array((modX * .5, modY * .5j, modX + modY * .5j, modX * .5 + modY * 1j))[np.newaxis, np.newaxis, :]  # centers of hexagons in rectangles

    t = t[..., np.newaxis]

    data = np.amin(np.abs(t-cc), axis=2) + np.minimum(np.sin(np.real(f) % (modX * 2) * np.pi), 1.2) / 5  # main image + long wave

    '''
    # a colormap and a normalization instance
    cmap = plt.cm.viridis
    norm = plt.Normalize(vmin=np.amin(data), vmax=np.amax(data))
    # map the normalized data to colors
    # image is now RGBA (WxHx4)
    image = cmap(norm(data))
    # save the image
    plt.imsave('image.png', image)
    '''
    plt.imsave('image.png', data, cmap=plt.cm.viridis)  # all above in one

    plt.imshow(data, extent=(-4, 4, -4, 4))
    plt.colorbar().set_label('val')
    plt.title(r'Möbius transformation: $\frac{\ln(M)}{F}$; $M=\frac{z+1}{z-1}$; $F=\frac{2 \pi i}{S_y i + S_x}$')
    plt.show()


if __name__ == '__main__':
    main()
