package token

import (
	"testing"
	"time"

	"github.com/dmsbyg/auth-service-demo/utils"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	testCases := []struct {
		desc      string
		secretKey string
		duration  time.Duration
		assert    func(maker *jwtMaker, err error)
	}{
		{
			desc:      "successfully create jwtMaker",
			secretKey: utils.RandomString(32),
			duration:  time.Hour,
			assert: func(maker *jwtMaker, err error) {
				require.NoError(t, err)
				require.NotNil(t, maker)
			},
		},
		{
			desc:      "secret length is too short",
			secretKey: utils.RandomString(3),
			duration:  time.Hour,
			assert: func(maker *jwtMaker, err error) {
				require.Error(t, err)
				require.Nil(t, maker)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			maker, err := NewJWTMaker(tC.secretKey, tC.duration)
			tC.assert(maker, err)
		})
	}
}

func TestCreateToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32), time.Minute)
	require.NoError(t, err)

	t.Run("create and verify valid token", func(t *testing.T) {
		duration := time.Minute
		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		userID := utils.RandomString(10)
		userEmail := utils.RandomString(10)

		token, err := maker.CreateToken(userID, userEmail)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := maker.Verify(token)
		require.NotNil(t, payload)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		require.Equal(t, userID, payload.ID)
		require.Equal(t, userEmail, payload.Email)

		require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
		require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
	})

	t.Run("must error when invalid user param given", func(t *testing.T) {
		userID := ""
		userEmail := utils.RandomString(10)

		token, err := maker.CreateToken(userID, userEmail)
		require.Error(t, err)
		require.Empty(t, token)
	})

	t.Run("create and verify valid token with default duration ", func(t *testing.T) {
		userID := utils.RandomString(10)
		userEmail := utils.RandomString(10)

		duration := time.Minute
		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		token, err := maker.CreateToken(userID, userEmail)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := maker.Verify(token)
		require.NotNil(t, payload)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		require.Equal(t, payload.ID, userID)
		require.Equal(t, payload.Email, userEmail)

		require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
		require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
	})

	t.Run("expired token", func(t *testing.T) {
		userName := utils.RandomString(10)
		userEmail := utils.RandomString(10)

		duration := -time.Minute

		maker, err = NewJWTMaker(utils.RandomString(secretKeyMinLength), duration)
		require.NoError(t, err)

		token, err := maker.CreateToken(userName, userEmail)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := maker.Verify(token)
		require.Error(t, err)
		require.Zero(t, payload)
	})

	t.Run("invalid token", func(t *testing.T) {
		token := utils.RandomString(10)
		payload, err := maker.Verify(token)
		require.Error(t, err)
		require.Zero(t, payload)
	})
}
