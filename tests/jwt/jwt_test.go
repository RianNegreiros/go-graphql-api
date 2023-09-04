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

		now := time.Now()

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithIssuer(conf.JWT.Issuer),
			jwtGo.WithVerify(jwa.HS256, []byte(conf.JWT.Secret)),
		)
		require.NoError(t, err)

		require.Equal(t, user.ID, tok.Subject())
		require.Equal(t, now.Add(jwt.AccessTokenLifeTime).Unix(), tok.Expiration().Unix())
	})
}
