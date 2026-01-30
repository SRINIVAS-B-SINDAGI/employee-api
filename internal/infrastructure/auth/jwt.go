package auth

import (
	"time"

	"github.com/SRINIVAS-B-SINDAGI/employee-api/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg config.JWTConfig) *JWTManager {

	return &JWTManager{
		secret:     []byte(cfg.Secret),
		expiration: cfg.Expiration,
		issuer:     cfg.Issuer,
	}
}

func (m *JWTManager) GenerateToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
