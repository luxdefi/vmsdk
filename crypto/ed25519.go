// Copyright (C) 2023-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Package crypto provides functionality for interacting with Ed25519
// public and private keys. Package crypto uses the "crypto/ed25519"
// package from  Go's standard library for the underlying cryptography.
package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"os"

	"github.com/luxdefi/node/utils/formatting/address"
)

type (
	PublicKey  [ed25519.PublicKeySize]byte
	PrivateKey [ed25519.PrivateKeySize]byte
	Signature  [ed25519.SignatureSize]byte
)

const (
	PublicKeyLen  = ed25519.PublicKeySize
	PrivateKeyLen = ed25519.PrivateKeySize
	// PrivateKeySeedLen is defined because ed25519.PrivateKey
	// is formatted as privateKey = seed|publicKey. We use this const
	// to extract the publicKey below.
	PrivateKeySeedLen = ed25519.SeedSize
	SignatureLen      = ed25519.SignatureSize
)

var (
	EmptyPublicKey  = [ed25519.PublicKeySize]byte{}
	EmptyPrivateKey = [ed25519.PrivateKeySize]byte{}
	EmptySignature  = [ed25519.SignatureSize]byte{}
)

// Address returns a Bech32 address from hrp and p.
// This function uses node's FormatBech32 function.
func Address(hrp string, p PublicKey) string {
	// TODO: handle error
	addrString, _ := address.FormatBech32(hrp, p[:])
	return addrString
}

// ParseAddress parses a Bech32 encoded address string and extracts
// its public key. If there is an error reading the address or the hrp
// value is not valid, ParseAddress returns an EmptyPublicKey and error.
func ParseAddress(hrp, saddr string) (PublicKey, error) {
	phrp, pk, err := address.ParseBech32(saddr)
	if err != nil {
		return EmptyPublicKey, err
	}
	if phrp != hrp {
		return EmptyPublicKey, ErrIncorrectHrp
	}
	// The parsed public key may be greater than [PublicKeyLen] because the
	// underlying Bech32 implementation requires bytes to each encode 5 bits
	// instead of 8 (and we must pad the input to ensure we fill all bytes):
	// https://github.com/btcsuite/btcd/blob/902f797b0c4b3af3f7196d2f5d2343931d1b2bdf/btcutil/bech32/bech32.go#L325-L331
	if len(pk) < PublicKeyLen {
		return EmptyPublicKey, ErrInvalidPublicKey
	}
	return PublicKey(pk[:PublicKeyLen]), nil
}

// GeneratePrivateKey returns a Ed25519 PrivateKey.
// It uses the crypto/ed25519 go package to generate the key.
func GeneratePrivateKey() (PrivateKey, error) {
	_, k, err := ed25519.GenerateKey(nil)
	if err != nil {
		return EmptyPrivateKey, err
	}
	return PrivateKey(k), nil
}

// PublicKey returns a PublicKey associated with the Ed25519 PrivateKey p.
// The PublicKey is the last 32 bytes of p.
func (p PrivateKey) PublicKey() PublicKey {
	return PublicKey(p[PrivateKeySeedLen:])
}

// ToHex converts a PrivateKey to a hex string.
func (p PrivateKey) ToHex() string {
	return hex.EncodeToString(p[:])
}

// Save writes [PrivateKey] to a file [filename]. If filename does
// not exist, it creates a new file with read/write permissions (0o600).
func (p PrivateKey) Save(filename string) error {
	return os.WriteFile(filename, p[:], 0o600)
}

// LoadKey returns a PrivateKey from a file filename.
// If there is an error reading the file, or the file contains an
// invalid PrivateKey, LoadKey returns an EmptyPrivateKey and an error.
func LoadKey(filename string) (PrivateKey, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return EmptyPrivateKey, err
	}
	if len(bytes) != ed25519.PrivateKeySize {
		return EmptyPrivateKey, ErrInvalidPrivateKey
	}
	return PrivateKey(bytes), nil
}

// Sign returns a valid signature for msg using pk.
func Sign(msg []byte, pk PrivateKey) Signature {
	sig := ed25519.Sign(pk[:], msg)
	return Signature(sig)
}

// Verify returns whether s is a valid signature of msg by p.
func Verify(msg []byte, p PublicKey, s Signature) bool {
	return ed25519.Verify(p[:], msg, s[:])
}

// HexToKey Converts a hexadecimal encoded key into a PrivateKey. Returns
// an EmptyPrivateKey and error if key is invalid.
func HexToKey(key string) (PrivateKey, error) {
	bytes, err := hex.DecodeString(key)
	if err != nil {
		return EmptyPrivateKey, err
	}
	if len(bytes) != ed25519.PrivateKeySize {
		return EmptyPrivateKey, ErrInvalidPrivateKey
	}
	return PrivateKey(bytes), nil
}
