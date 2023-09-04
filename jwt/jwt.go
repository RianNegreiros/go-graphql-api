package jwt

import (
	"context"
	"fmt"
	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/models"
	_ "github.com/lestrrat-go/jwx"
	"github.com/lestrrat-go/jwx/jwa"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"time"
)

var signatureType = jwa.HS256

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{Conf: conf}
}

func (s *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (models.AuthToken, error) {
	token, err := jwtGo.ParseRequest(
		r,
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(s.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(s.Conf.JWT.Secret)),
	)
	if err != nil {
		return models.AuthToken{}, models.ErrInvalidToken
	}

	return buildToken(token), nil
}

func buildToken(token jwtGo.Token) models.AuthToken {
	return models.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}
}

func (s *TokenService) ParseToken(ctx context.Context, payload string) (models.AuthToken, error) {
	token, err := jwtGo.Parse(
		[]byte(payload),
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(s.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(s.Conf.JWT.Secret)),
	)
	if err != nil {
		return models.AuthToken{}, models.ErrInvalidToken
	}

	return buildToken(token), nil
}

func (s *TokenService) CreateRefreshToken(ctx context.Context, user models.User, tokenID string) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, RefreshTokenLifeTime, s.Conf); err != nil {
		return "", err
	}

	if err := t.Set(jwtGo.JwtIDKey, tokenID); err != nil {
		return "", fmt.Errorf("failed to set jwt id: %w", err)
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(s.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return string(token), nil
}

func (s *TokenService) CreateAccessToken(ctx context.Context, user models.User) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, AccessTokenLifeTime, s.Conf); err != nil {
		return "", err
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(s.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return string(token), nil
}

func setDefaultToken(t jwtGo.Token, user models.User, lifetime time.Duration, conf *config.Config) error {
	if err := t.Set(jwtGo.SubjectKey, user.ID); err != nil {
		return fmt.Errorf("failed to set jwt subject: %w", err)
	}

	if err := t.Set(jwtGo.IssuerKey, conf.JWT.Issuer); err != nil {
		return fmt.Errorf("failed to set jwt issuer at: %w", err)
	}

	if err := t.Set(jwtGo.IssuedAtKey, time.Now().Unix()); err != nil {
		return fmt.Errorf("failed to set jwt issued at: %w", err)
	}

	if err := t.Set(jwtGo.ExpirationKey, time.Now().Add(lifetime)); err != nil {
		return fmt.Errorf("failed to set jwt expiration: %w", err)
	}

	return nil
}
