package edkey

// Originally from https://github.com/mikesmitty/edkey/
// MIT License
// Copyright (c) 2017 Michael Smith

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func MarshalED25519PrivateKey(key ed25519.PrivateKey) ([]byte, error) {
	var w struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	}

	// Fill out the private key fields
	pk1 := struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}{}

	// Set our check ints
	ci := rand.Uint32()
	pk1.Check1 = ci
	pk1.Check2 = ci

	// Set our key type
	pk1.Keytype = ssh.KeyAlgoED25519

	// Add the pubkey to the optionally-encrypted block
	publicKey, ok := key.Public().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast to ed25519.PublicKey")
	}

	pk1.Pub = publicKey

	// Add our private key
	pk1.Priv = key

	// Might be useful to put something in here at some point
	pk1.Comment = ""

	// Add some padding to match the encryption block size within PrivKeyBlock (without Pad field)
	// 8 doesn't match the documentation, but that's what ssh-keygen uses for unencrypted keys. *shrug*
	bs := 8
	blockLen := len(ssh.Marshal(pk1))
	padLen := (bs - (blockLen % bs)) % bs
	pk1.Pad = make([]byte, padLen)

	// Padding is a sequence of bytes like: 1, 2, 3...
	for i := 0; i < padLen; i++ {
		pk1.Pad[i] = byte(i + 1)
	}

	// Generate the pubkey prefix "\0\0\0\nssh-ed25519\0\0\0 "
	prefix := []byte{0x0, 0x0, 0x0, 0x0b}
	prefix = append(prefix, []byte(ssh.KeyAlgoED25519)...)
	prefix = append(prefix, []byte{0x0, 0x0, 0x0, 0x20}...)

	// Only going to support unencrypted keys for now
	w.CipherName = "none"
	w.KdfName = "none"
	w.KdfOpts = ""
	w.NumKeys = 1
	w.PubKey = append(prefix, publicKey...)
	w.PrivKeyBlock = ssh.Marshal(pk1)

	// Add key header (followed by a null byte)
	magic := append([]byte("openssh-key-v1"), 0)
	magic = append(magic, ssh.Marshal(w)...)

	return magic, nil
}
