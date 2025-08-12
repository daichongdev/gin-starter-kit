package auth

import (
	"errors"
	"gin-demo/config"
	"gin-demo/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, email string) (string, time.Time, error) {
	cfg := config.GetConfig()

	expiresAt := time.Now().Add(time.Duration(cfg.JWT.ExpiresHours) * time.Hour)

	claims := &model.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*model.JWTClaims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证Token是否有效
func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查Token是否过期
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// GenerateRefreshToken 生成刷新Token
func GenerateRefreshToken(userID uint, email string) (string, time.Time, error) {
	cfg := config.GetConfig()

	expiresAt := time.Now().Add(time.Duration(cfg.JWT.RefreshExpiresHours) * time.Hour)

	claims := &model.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
