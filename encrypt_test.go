package main

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	plain := "test"
	key, err := generateKey()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	encrypted, err := encrypt([]byte(plain), []byte(key))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(encrypted) != 44 {
		t.Fatalf("encrypted string length is invalid: %d", len(encrypted))
	}

	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if plain != decrypted {
		t.Fatalf("expected: '%s', but '%s'", plain, decrypted)
	}
}

func TestGenerateKey(t *testing.T) {
	first, err := generateKey()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	second, err := generateKey()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if bytes.Equal(first, second) {
		t.Fatal("generateKey is not random")
	}
}

func TestPadPKCS7(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected []byte
	}{
		{
			[]byte{},
			[]byte{16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16},
		},
		{
			[]byte{1},
			[]byte{1, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15},
		},
		{
			[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	for _, testCase := range testCases {
		actual := padPKCS7(testCase.data)
		if !bytes.Equal(actual, testCase.expected) {
			t.Fatalf("expected: '%v', but '%v'", testCase.expected, actual)
		}
	}
}
