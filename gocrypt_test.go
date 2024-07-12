package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

func createTestFile() (string, []byte, error) {
	testFile, err := os.CreateTemp("./", "gocrypt_test_*.bin")
	if err != nil {
		return "", nil, err
	}
	defer testFile.Close()

	testData := make([]byte, 1024)
	if _, err := io.ReadFull(rand.Reader, testData); err != nil {
		return "", nil, err
	}

	if _, err := testFile.Write(testData); err != nil {
		return "", nil, err
	}

	return testFile.Name(), testData, nil
}

func TestEncryptAndDecrypt(t *testing.T) {
	testFileName, testData, err := createTestFile()
	if err != nil {
		t.Error(err)
		return
	}

	testPassword := "testing_123456789"
	encryptedFileName := testFileName + "_encrypted"
	if err := encryptFile(&testFileName, &encryptedFileName, testPassword); err != nil {
		t.Error(err)
	}

	decryptedFileName := testFileName + "_decrypted"
	if err := decryptFile(&encryptedFileName, &decryptedFileName, testPassword); err != nil {
		t.Error(err)
	}

	encryptedFile, err := os.Open(encryptedFileName)
	if err != nil {
		t.Error(err)
	}

	decryptedFile, err := os.Open(decryptedFileName)
	if err != nil {
		t.Error(err)
	}

	decryptedData := make([]byte, len(testData))
	if _, err := decryptedFile.Read(decryptedData); err != nil {
		t.Error(err)
	}

	defer func() {
		encryptedFile.Close()
		decryptedFile.Close()
		os.Remove(testFileName)
		os.Remove(encryptedFileName)
		os.Remove(decryptedFileName)
	} ()

	if hex.EncodeToString(decryptedData) != hex.EncodeToString(testData) {
		t.Error("Failed to match encrypted file against decrypted file.")
	}
}
