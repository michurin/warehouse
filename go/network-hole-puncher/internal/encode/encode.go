package encode

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func Pack(d []byte) ([]byte, error) {
	data := append(make([]byte, 32), d...)
	_, err := rand.Read(data[28:32])
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum224(data[28:])
	copy(data[:28], sum[:])
	return data, nil
}

func Unpack(d []byte) ([]byte, error) {
	if len(d) < 32 {
		return nil, fmt.Errorf("Invalid length: %d", len(d))
	}
	sum := sha256.Sum224(d[28:])
	for i, v := range sum {
		if v != d[i] {
			return nil, fmt.Errorf("Invalid checksum")
		}
	}
	return d[32:], nil
}
