#!/usr/bin/env python3


# 1D random midpoint displacement
# based on
# [1] https://web.williams.edu/Mathematics/sjmiller/public_html/hudson/Dickerson_Terrain.pdf
# [2] https://www.uni-konstanz.de/mmsp/pubsys/publishedFiles/Saupe88c.pdf
#
# 2D
# http://www-evasion.imag.fr/~Fabrice.Neyret/images/fluids-nuages/waves/Jonathan/articlesCG/fourier-synthesis-of-ocean-scenes-87.pdf

import numpy as np
import matplotlib.pyplot as plt


def main():
    R = 2. # 1.4  # D = (5 - R)/2; (2.11) in [2]; R=2 - brown noise
    N = 50
    x = np.linspace(0, 1, 500)
    amp = np.power(np.arange(1, N+1), -R)  # magnitude
    amp = amp * np.random.uniform(0, 1, N)  # rand
    amp = amp * np.exp(1j*np.random.uniform(0, 2*np.pi, N))  # rand phase
    amp = amp[np.newaxis, ...]
    freq = np.arange(1, N+1)[np.newaxis, ...]
    y = np.exp(x[..., np.newaxis] * 2j * np.pi * freq)
    y = np.sum(y * amp, axis=1)
    plt.plot(x, y.real, x, y.imag)
    plt.show()


if __name__ == '__main__':
    main()
