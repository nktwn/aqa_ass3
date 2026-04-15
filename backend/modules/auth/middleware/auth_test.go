package middleware

import (
	"diploma/modules/auth/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockJWT struct {
	mock.Mock
}

func (m *mockJWT) GenerateJSONWebTokens(id int64, username string, role int) (string, string, error) {
	args := m.Called(id, username, role)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockJWT) RefreshAccessToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *mockJWT) VerifyToken(tokenStr string) (*jwt.Claims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

type AuthMiddlewareTestSuite struct {
	suite.Suite
	router     *gin.Engine
	middleware *AuthMiddleware
	mockJWT    *mockJWT
}

func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.New()
	s.mockJWT = new(mockJWT)
	s.middleware = NewAuthMiddleware(s.mockJWT)

	// Set up a test route with the auth middleware
	s.router.Use(s.middleware.AuthMiddleware())
	s.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid.token.here")

	// Set up mock expectations
	claims := &jwt.Claims{
		UserID:   1,
		Username: "testuser",
		Role:     0,
	}
	s.mockJWT.On("VerifyToken", "valid.token.here").Return(claims, nil).Once()

	// Perform the request
	s.router.ServeHTTP(w, req)

	// Assert response
	s.Equal(http.StatusOK, w.Code)
	s.mockJWT.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_NoAuthHeader() {
	// Create a test request without Authorization header
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Perform the request
	s.router.ServeHTTP(w, req)

	// Assert response
	s.Equal(http.StatusUnauthorized, w.Code)
	s.mockJWT.AssertNotCalled(s.T(), "VerifyToken")
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidAuthHeader() {
	// Create a test request with invalid Authorization header format
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	// Perform the request
	s.router.ServeHTTP(w, req)

	// Assert response
	s.Equal(http.StatusUnauthorized, w.Code)
	s.mockJWT.AssertNotCalled(s.T(), "VerifyToken")
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	// Set up mock expectations
	s.mockJWT.On("VerifyToken", "invalid.token.here").Return(nil, jwt.ErrInvalidToken).Once()

	// Perform the request
	s.router.ServeHTTP(w, req)

	// Assert response
	s.Equal(http.StatusUnauthorized, w.Code)
	s.mockJWT.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestAuthMiddleware_ExpiredToken() {
	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer expired.token.here")

	// Set up mock expectations
	s.mockJWT.On("VerifyToken", "expired.token.here").Return(nil, jwt.ErrTokenExpired).Once()

	// Perform the request
	s.router.ServeHTTP(w, req)

	// Assert response
	s.Equal(http.StatusUnauthorized, w.Code)
	s.mockJWT.AssertExpectations(s.T())
}
