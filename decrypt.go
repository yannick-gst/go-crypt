package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

func decryptChunk(gcm cipher.AEAD, nonce []byte, chunk []byte, out *os.File) error {
	decryptedChunk, err := gcm.Open(nil, nonce, chunk, nil)
	if err != nil {
		return err
	}
	if _, err := out.Write(decryptedChunk); err != nil {
		return err
	}
	return nil
}

func decryptFile(inputFile, outputFile *string, password string) error {
	// Open input file stream
	ifs, err := os.Open(*inputFile)
	if err != nil {
		return err
	}
	defer ifs.Close()

	// Create output file stream
	ofs, err := os.Create(*outputFile)
	if err != nil {
		return err
	}
	defer ofs.Close()

	// Read the salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(ifs, salt); err != nil {
		return err
	}

	// Generate encryption key and AES cipher block
	key := generateKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a Galois Counter Mode (GCM) instance. GCM encrypts the data, and
	// also generates a checksum that can be used to verify the data during
	// decryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	buffer := make(
		[]byte,
		// The overhead is the maximum difference in length betwen a text and
		// its ciphertext
		chunkSize + gcm.Overhead()) // Buffer in which to write chunks

	// Decrypt input file in chunks
	for {
		// Read the nonce
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(ifs, nonce); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Read and decrypt chunk
		bytesRead, err := ifs.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if bytesRead == 0 {
			break
		}

		if err := decryptChunk(gcm, nonce, buffer[:bytesRead], ofs); err != nil {
			return err
		}
	}

	return nil
}
