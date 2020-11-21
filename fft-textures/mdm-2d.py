#!/usr/bin/env python3

import numpy as np
import matplotlib.pyplot as plt

def main():
    N = 401
    noise = np.random.uniform(-1, 1, (N, N))  # it is important that noise is balased around 0
    noise_fft = np.fft.fft2(noise)
    t = np.arange(N, dtype=np.float)
    r = np.hypot(t[np.newaxis, ...], t[..., np.newaxis])
    #mask = np.logical_and(r > 12, r < 20).astype(float)
    mask = np.maximum(16 - np.power(r - 16, 2.), 0)
    noise_fft_masked = noise_fft * mask
    vawes = np.fft.ifft2(noise_fft_masked)

    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(8, 8))

    ax1.imshow(noise)
    ax1.set_title('noise')

    ax2.imshow(np.abs(noise_fft_masked))
    ax2.set_title('abs(mask*fft)')

    ax3.imshow(vawes.real)
    ax3.set_title('ifft(mask*fft)')

    ax4.imshow(np.sin(vawes.real * 10 * np.pi))
    ax4.set_title('sin(ifft(mask*fft))')

    plt.tight_layout()

    plt.show()

if __name__ == '__main__':
    main()
