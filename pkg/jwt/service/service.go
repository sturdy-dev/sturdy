package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"mash/pkg/jwt"
	"mash/pkg/jwt/keys"
	db_keys "mash/pkg/jwt/keys/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
	jose_jwt "gopkg.in/square/go-jose.v2/jwt"
)

// Known errors.
var (
	ErrInvalidToken = fmt.Errorf("token is invalid")
	ErrTokenExpired = fmt.Errorf("token is expired")
)

const (
	defaultIssuer = "https://getsturdy.com"
)

type Service struct {
	logger   *zap.Logger
	keysRepo db_keys.Repository

	signerInitOnce *sync.Once
	signer         jose.Signer
}

func NewService(logger *zap.Logger, keysRepo db_keys.Repository) *Service {
	return &Service{
		logger:   logger,
		keysRepo: db_keys.NewCache(keysRepo),

		signerInitOnce: &sync.Once{},
	}
}

func (s *Service) initOnce(ctx context.Context) error {
	var err error
	s.signerInitOnce.Do(func() {
		err = s.init(ctx)
	})
	return err
}

// init prepares jwt service for work by generating a jwt signer.
func (s *Service) init(ctx context.Context) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	publicDER, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return fmt.Errorf("failed to marshal encryption key: %w", err)
	}

	key, err := keys.New(publicDER)
	if err != nil {
		return fmt.Errorf("failed to create a key: %w", err)
	}

	if err := s.keysRepo.Create(ctx, key); err != nil {
		return fmt.Errorf("failed to store key in the database: %w", err)
	}

	options := (&jose.SignerOptions{}).
		WithHeader("kid", key.ID).
		WithType("JWT")

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       privateKey,
	}, options)
	if err != nil {
		return err
	}

	s.signer = signer

	return nil
}

func (s *Service) IssueToken(ctx context.Context, subject string, validFor time.Duration, tokenType jwt.TokenType) (*jwt.Token, error) {
	if err := s.initOnce(ctx); err != nil {
		return nil, err
	}

	now := time.Now()
	stdClaims := &jose_jwt.Claims{
		ID:       uuid.New().String(),
		Issuer:   defaultIssuer,
		Subject:  subject,
		IssuedAt: jose_jwt.NewNumericDate(now),
		Expiry:   jose_jwt.NewNumericDate(now.Add(validFor)),
	}
	sturdyClaims := jwtClaims{
		Type: tokenType,
	}

	token, err := jose_jwt.Signed(s.signer).
		Claims(stdClaims).
		Claims(sturdyClaims).
		CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("failed to create a signed token: %w", err)
	}
	return &jwt.Token{
		Token:     token,
		Subject:   stdClaims.Subject,
		ExpiresAt: stdClaims.Expiry.Time(),
		Type:      sturdyClaims.Type,
	}, nil
}

// Verify checks token signature and returns it's meaningful content.
func (s *Service) Verify(ctx context.Context, rawToken string, expectedTypes ...jwt.TokenType) (*jwt.Token, error) {
	jwtoken, err := jose_jwt.ParseSigned(rawToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if len(jwtoken.Headers) == 0 {
		return nil, ErrInvalidToken
	}
	header := jwtoken.Headers[0]

	switch header.Algorithm {
	case string(jose.ES256):
		return s.verify(ctx, rawToken, jwtoken, expectedTypes...)
	default:
		return nil, ErrInvalidToken
	}
}

func (s *Service) get(ctx context.Context, id string) (*ecdsa.PublicKey, error) {
	key, err := s.keysRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	untypedResult, err := x509.ParsePKIXPublicKey(key.PublicDER)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PKIX public key: %w", err)
	}

	switch v := untypedResult.(type) {
	case *ecdsa.PublicKey:
		return v, nil
	default:
		return nil, fmt.Errorf("unknown public key type: %T", v)
	}
}

func validType(typ jwt.TokenType, expected ...jwt.TokenType) bool {
	for _, t := range expected {
		if t == typ {
			return true
		}
	}
	return false
}

// verify verifies token using ECDSA algorithm with a key from the database.
func (s *Service) verify(ctx context.Context, rawToken string, jwtoken *jose_jwt.JSONWebToken, expectedTypes ...jwt.TokenType) (*jwt.Token, error) {
	id := jwtoken.Headers[0].KeyID

	pubicKey, err := s.get(ctx, id)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	stdClaims := &jose_jwt.Claims{}
	sturdyClaims := &jwtClaims{}
	if err := jwtoken.Claims(pubicKey, stdClaims, sturdyClaims); err != nil {
		return nil, ErrInvalidToken
	}

	validateErr := stdClaims.ValidateWithLeeway(jose_jwt.Expected{
		Time:   time.Now(),
		Issuer: defaultIssuer,
	}, time.Second)
	switch {
	case validateErr == nil:
		if !validType(sturdyClaims.Type, expectedTypes...) {
			return nil, ErrInvalidToken
		}
		return &jwt.Token{
			Token:     rawToken,
			Subject:   stdClaims.Subject,
			ExpiresAt: stdClaims.Expiry.Time(),
			Type:      sturdyClaims.Type,
		}, nil
	case errors.Is(validateErr, jose_jwt.ErrExpired):
		return nil, ErrTokenExpired
	default:
		return nil, ErrInvalidToken
	}
}
