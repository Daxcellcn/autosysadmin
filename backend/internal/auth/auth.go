// backend/internal/auth/auth.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password string) (string, error)
	Login(email, password string) (string, string, error)
	VerifyToken(tokenString string) (*Claims, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type authService struct {
	secretKey        []byte
	refreshSecretKey []byte
	users           map[string]User // In-memory store for demo, replace with DB
}

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Roles        []string
}

func NewAuthService(secretKey, refreshSecretKey string) AuthService {
	return &authService{
		secretKey:        []byte(secretKey),
		refreshSecretKey: []byte(refreshSecretKey),
		users:           make(map[string]User),
	}
}

func (s *authService) Register(email, password string) (string, error) {
	if _, exists := s.users[email]; exists {
		return "", errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := User{
		ID:           generateUUID(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Roles:        []string{"user"},
	}

	s.users[email] = user
	return user.ID, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, exists := s.users[email]
	if !exists {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.refreshSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	user, exists := s.users[claims.Email]
	if !exists {
		return "", "", errors.New("user not found")
	}

	newAccessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) generateAccessToken(user *User) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "autosysadmin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *authService) generateRefreshToken(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour * 7) // 7 days
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "autosysadmin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecretKey)
}

func generateUUID() string {
	// In production, use github.com/google/uuid
	return "generated-uuid"
}