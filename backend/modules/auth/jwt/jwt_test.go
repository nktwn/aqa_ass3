package jwt

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndVerifyTokens(t *testing.T) {
	t.Parallel()

	jwtService := NewJSONWebToken("test-secret-key")

	accessToken, refreshToken, err := jwtService.GenerateJSONWebTokens(1, "testuser", 1)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	accessClaims, err := jwtService.VerifyToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(1), accessClaims.UserID)
	assert.Equal(t, "testuser", accessClaims.Username)
	assert.Equal(t, 1, accessClaims.Role)

	refreshClaims, err := jwtService.VerifyToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, int64(1), refreshClaims.UserID)
	assert.Equal(t, "testuser", refreshClaims.Username)
	assert.Equal(t, 1, refreshClaims.Role)
}

func TestRefreshAccessToken(t *testing.T) {
	t.Parallel()

	jwtService := NewJSONWebToken("test-secret-key")
	_, refreshToken, err := jwtService.GenerateJSONWebTokens(5, "refresh-user", 0)
	require.NoError(t, err)

	newAccessToken, err := jwtService.RefreshAccessToken(refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, newAccessToken)

	claims, err := jwtService.VerifyToken(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(5), claims.UserID)
	assert.Equal(t, "refresh-user", claims.Username)
	assert.Equal(t, 0, claims.Role)
}

func TestVerifyTokenRejectsInvalidAndExpiredTokens(t *testing.T) {
	t.Parallel()

	jwtService := NewJSONWebToken("test-secret-key")

	_, err := jwtService.VerifyToken("invalid.token.here")
	assert.ErrorIs(t, err, ErrInvalidToken)

	expiredToken, err := jwtService.generateJSONWebToken(1, 0, "expired-user", time.Now().Add(-1*time.Hour))
	require.NoError(t, err)

	_, err = jwtService.VerifyToken(expiredToken)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestRefreshAccessTokenRejectsInvalidToken(t *testing.T) {
	t.Parallel()

	jwtService := NewJSONWebToken("test-secret-key")
	_, err := jwtService.RefreshAccessToken("invalid.refresh.token")
	assert.Error(t, err)
}

func TestVerifyTokenRejectsTamperedToken(t *testing.T) {
	t.Parallel()

	jwtService := NewJSONWebToken("test-secret-key")
	accessToken, _, err := jwtService.GenerateJSONWebTokens(1, "testuser", 1)
	require.NoError(t, err)

	tamperedToken := accessToken[:len(accessToken)-1] + "x"
	_, err = jwtService.VerifyToken(tamperedToken)
	assert.ErrorIs(t, err, ErrInvalidToken)
}


func TestExtractTokenFromHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		header        string
		expectedToken string
		expectedError bool
	}{
		{
			name:          "valid bearer token",
			header:        "Bearer valid-token",
			expectedToken: "valid-token",
		},
		{
			name:          "missing bearer prefix",
			header:        "valid-token",
			expectedError: true,
		},
		{
			name:          "empty header",
			header:        "",
			expectedError: true,
		},
		{
			name:          "invalid format",
			header:        "Bearer token extra",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			token, err := ExtractTokenFromHeader(req)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}
