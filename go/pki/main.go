package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func justOneBlockEncrypton() {
	c, err := aes.NewCipher([]byte("0123456701234567"))
	noerr(err)
	fmt.Println(c.BlockSize()) // 16

	msg := []byte("MICHURIN--ALEXEY") // 16
	enc := make([]byte, 16)
	c.Encrypt(enc, msg)
	fmt.Printf("%q (%d)\n", enc, len(enc)) // encrypting

	dec := make([]byte, 16)
	c.Decrypt(dec, enc)
	fmt.Printf("%q\n", dec) // getting back
}

func messageEncyption() {
	msg := []byte("Michurin")
	c, err := aes.NewCipher([]byte("0123456701234567"))
	noerr(err)

	stream, err := cipher.NewGCM(c) // AEAD better than CFB
	noerr(err)
	fmt.Println(stream.NonceSize())                                    // 12
	enc := stream.Seal(nil, []byte("123123123123"), msg, []byte("xx")) // seal works like append, nonce must be unique, nonce must have correct size
	fmt.Printf("%q (%d)\n", enc, len(enc))

	dec, err := stream.Open(nil, []byte("123123123123"), enc, []byte("xx"))
	noerr(err)
	fmt.Printf("%q\n", dec)
}

func main() {
	justOneBlockEncrypton()
	messageEncyption()
}

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}
