package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ssov1 "github.com/zaketn/sso-contracts/gen/go/sso"
	"github.com/zaketn/sso/tests/suite"
	"math/rand"
	"testing"
	"time"
)

const (
	emptyAppId = 0
	appId      = 1
	appSecret  = "test-secret"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, rand.Intn(64))

	registerReq := &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	}

	registerResp, err := st.AuthClient.Register(ctx, registerReq)
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.GetUserId())

	loginReq := &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	}

	loginResp, err := st.AuthClient.Login(ctx, loginReq)

	loginTime := time.Now()

	token := loginResp.GetToken()

	require.NoError(t, err)
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, registerResp.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}
