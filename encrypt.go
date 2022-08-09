package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Encrypt and compress shellcode and save as file
func encryptShellcode(shellcode []byte, keyString string) {
	var buf bytes.Buffer
	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		fmt.Printf("[*] Error creating AES cipher block. Error message: %s\n", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, shellcode, nil)

	gz := gzip.NewWriter(&buf)
	gz.Write(ciphertext)
	gz.Close()
	outputFile := filepath.FromSlash("pkgs/encrypted_shellcode.bin")
	ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
}
