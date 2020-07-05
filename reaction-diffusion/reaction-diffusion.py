#!/usr/bin/env python3

"""
More info about RDS:
https://en.wikipedia.org/wiki/Reaction%E2%80%93diffusion_system

Reaction:
    - consider two components: A and B
    - Da, Db — mass diffusivity or diffusion coefficient
    - Reaction rule: B-B-A -> B-B-B (A turns to B if A meets two Bs)
    - F and K — artificial in/out flows

Model:
    dA = (Da · ∆A - A·B² + F(1 - A)) · dt
    dB = (Db · ∆B + A·B² + (K + F) · B) · dt
    Where: F — "feed", K — "kill"

The way to melt images together into video clip:
ffmpeg -i image-%05d.png -c:v libx264 -r 30 -pix_fmt yuv420p out.mp4
With sound to publish on YouTube:
ffmpeg -ar 48000 -ac 2 -f s16le -i /dev/zero -i image-%05d.png -shortest -c:a aac -c:v libx264 -r 30 -pix_fmt yuv420p -strict experimental out.mp4
"""

import matplotlib.pyplot as plt

import numpy as np

import scipy.ndimage as nd


def init(size, spot_size):
    d = (size, size)
    a = np.ones(d, dtype=float)
    b = np.zeros(d, dtype=float)
    p = int(size / 2) - spot_size
    q = int(size / 2) + spot_size
    b[p:q, p:q] = 1
    return a, b


def images(a, b, Da, Db, F, K, dt, steps):
    n = 0
    for s in steps:
        while n < s:
            abb = a * b * b
            a += (Da * nd.filters.laplace(a) - abb + F * (1.0 - a)) * dt
            b += (Db * nd.filters.laplace(b) + abb - (K + F) * b) * dt
            n += 1
        yield b


def main():
    a, b = init(500, 10)
    for n, b in enumerate(images(a, b, .095, .055, .045, .062, 1, range(0, 100000, 50))):
        # plt.imshow(b)
        # plt.show()
        print(f'step={n}')
        plt.imsave(f'image-{n:05d}.png', b, cmap=plt.cm.viridis)


if __name__ == '__main__':
    main()
