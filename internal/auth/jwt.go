package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"sync"
	"time"
)

type JWTMiddleware struct {
	accessSecret  []byte
	refreshSecret []byte
	userStore     *UserStore
	blacklist     *TokenBlacklist
}

type UserStore struct {
	users map[string]string // username -> password hash
	mu    sync.RWMutex
}

type TokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

func (m *JWTMiddleware) ValidateRequest(r *http.Request) (string, error) {
	return m.extractAndValidateToken(r)
}

func (m *JWTMiddleware) IsBlacklisted(token string) bool {
	return m.blacklist.IsBlacklisted(token)
}

func NewJWTMiddleware(accessSecret, refreshSecret string) *JWTMiddleware {
	return &JWTMiddleware{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		userStore:     NewUserStore(),
		blacklist:     NewTokenBlacklist(),
	}
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]string),
	}
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
}

func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/register" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := m.extractAndValidateToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if m.blacklist.IsBlacklisted(token) {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *JWTMiddleware) Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := m.userStore.AddUser(creds.Username, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (m *JWTMiddleware) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !m.userStore.ValidateUser(creds.Username, creds.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokens, err := m.generateTokenPair(creds.Username)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (m *JWTMiddleware) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshReq RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(refreshReq.RefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.refreshSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	tokens, err := m.generateTokenPair(claims.Username)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (m *JWTMiddleware) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := m.extractAndValidateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	m.blacklist.Add(token, time.Now().Add(24*time.Hour))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (m *JWTMiddleware) generateTokenPair(username string) (*TokenResponse, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessTokenString, err := accessToken.SignedString(m.accessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	refreshTokenString, err := refreshToken.SignedString(m.refreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

func (m *JWTMiddleware) extractAndValidateToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	token, err := jwt.ParseWithClaims(bearerToken[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.accessSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return bearerToken[1], nil
}

func (s *UserStore) AddUser(username, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return errors.New("user already exists")
	}

	s.users[username] = password
	return nil
}

func (s *UserStore) ValidateUser(username, password string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storedPassword, exists := s.users[username]
	return exists && storedPassword == password
}

func (b *TokenBlacklist) Add(token string, expiresAt time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens[token] = expiresAt
}

func (b *TokenBlacklist) IsBlacklisted(token string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	expiry, exists := b.tokens[token]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		b.mu.RUnlock()
		b.mu.Lock()
		delete(b.tokens, token)
		b.mu.Unlock()
		return false
	}

	return true
}
