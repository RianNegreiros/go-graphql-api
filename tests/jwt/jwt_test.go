package jwt

import (
	"context"
	"github.com/RianNegreiros/go-graphql-api/config"
	"github.com/RianNegreiros/go-graphql-api/jwt"
	"github.com/RianNegreiros/go-graphql-api/models"
	"github.com/lestrrat-go/jwx/jwa"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	conf         *config.Config
	tokenService *jwt.TokenService
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env.test")
	conf = config.New()

	tokenService = jwt.NewTokenService(conf)

	os.Exit(m.Run())
}

func TestTokenService_CreateAccessToken(t *testing.T) {
	t.Run("should create access token", func(t *testing.T) {
		ctx := context.Background()
		user := models.User{
			ID: "1",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		jwt.Now = func() time.Time {
			return time.Now()
		}

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithIssuer(conf.JWT.Issuer),
			jwtGo.WithVerify(jwa.HS256, []byte(conf.JWT.Secret)),
		)
		require.NoError(t, err)

		require.Equal(t, user.ID, tok.Subject())
		require.Equal(t, jwt.Now().Add(jwt.AccessTokenLifeTime).Unix(), tok.Expiration().Unix())

		teardownTimeNow(t)
	})
}

func TestTokenService_CreateRefreshToken(t *testing.T) {
	t.Run("should create refresh token", func(t *testing.T) {
		ctx := context.Background()
		user := models.User{
			ID: "1",
		}

		token, err := tokenService.CreateRefreshToken(ctx, user, "2")
		require.NoError(t, err)

		jwt.Now = func() time.Time {
			return time.Now()
		}

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithIssuer(conf.JWT.Issuer),
			jwtGo.WithVerify(jwa.HS256, []byte(conf.JWT.Secret)),
		)
		require.NoError(t, err)

		require.Equal(t, user.ID, tok.Subject())
		require.Equal(t, "2", tok.JwtID())
		require.Equal(t, jwt.Now().Add(jwt.RefreshTokenLifeTime).Unix(), tok.Expiration().Unix())

		teardownTimeNow(t)
	})
}

func TestTokenService_ParseToken(t *testing.T) {
	t.Run("should parse valid token", func(t *testing.T) {
		ctx := context.Background()
		user := models.User{
			ID: "1",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token)
		require.NoError(t, err)

		require.Equal(t, user.ID, tok.Sub)
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		ctx := context.Background()
		user := models.User{
			ID: "1",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token+"invalid")
		require.Error(t, err)
		require.Equal(t, models.ErrInvalidToken, err)
		require.Equal(t, models.AuthToken{}, tok)
	})

	t.Run("should return error when token is expired", func(t *testing.T) {
		ctx := context.Background()
		user := models.User{
			ID: "1",
		}

		jwt.Now = func() time.Time {
			return time.Now().Add(-jwt.AccessTokenLifeTime * 5)
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		_, err = tokenService.ParseToken(ctx, token)
		require.ErrorIs(t, err, models.ErrInvalidToken)

		teardownTimeNow(t)
	})
}

func teardownTimeNow(t *testing.T) {
	t.Helper()

	jwt.Now = func() time.Time {
		return time.Now()
	}
}
