package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/ScaleFT/sshkeys"
	"golang.org/x/crypto/ssh"

	"getsturdy.com/api/pkg/crypto"
)

func genPk() (*rsa.PrivateKey, error) {
	bitSize := 4096

	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, fmt.Errorf("could not generate new private key: %w", err)
	}

	return privateKey, nil
}

func GenerateRsaKeypair() (public crypto.PublicKey, private crypto.PrivateKey, err error) {
	privateKey, err := genPk()
	if err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create public key: %w", err)
	}
	public = crypto.PublicKey(fmt.Sprintf("%s", ssh.MarshalAuthorizedKey(pub)))

	var buf bytes.Buffer
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&buf, privateKeyPEM); err != nil {
		return "", "", fmt.Errorf("failed to encode private key: %w", err)
	}
	private = crypto.PrivateKey(buf.String())
	return
}

func GenerateRsaSSHKeypair() (public crypto.PublicKey, private crypto.PrivateKey, err error) {
	privateKey, err := genPk()
	if err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create public key: %w", err)
	}
	public = crypto.PublicKey(fmt.Sprintf("%s", ssh.MarshalAuthorizedKey(pub)))

	privk, err := sshkeys.Marshal(privateKey, &sshkeys.MarshalOptions{Format: sshkeys.FormatOpenSSHv1})
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}
	private = crypto.PrivateKey(privk)
	return
}
