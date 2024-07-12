package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptChunk(gcm cipher.AEAD, nonce []byte, chunk []byte, out *os.File) error {
	defer out.Close()
	encryptedChunk := gcm.Seal(nil, nonce, chunk, nil)
	if _, err := out.Write(encryptedChunk); err != nil {
		return err
	}
	return nil
}

func encryptFile(inputFile, outputFile *string, password string) error {
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

	// Generate random salt, encryption key and AES cipher block
	salt, err := generateSalt()
	if err != nil {
		return err
	}

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

	// Write salt to output file
	if _, err := ofs.Write(salt); err != nil {
		return err
	}

	buffer := make([]byte, chunkSize) // Buffer in which to write chunks

	// Encrypt input file in chunks
	for {
		bytesRead, err := ifs.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if bytesRead == 0 {
			break
		}

		// Generate a random nonce for each chunnk. This ensures data
		// confidentiality and helps prevent replay attacks.
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}
		if _, err := ofs.Write(nonce); err != nil {
			return err
		}

		if err := encryptChunk(gcm, nonce, buffer[:bytesRead], ofs); err != nil {
			return err
		}
	}

	return nil
}
