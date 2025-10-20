package util

import (
	"context"
	"errors"
	"pura-agung-kertajaya-backend/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenUtil struct {
	SecretKey string
	Redis     *redis.Client
}

func NewTokenUtil(secretKey string, redisClient *redis.Client) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
		Redis:     redisClient,
	}
}

func (t *TokenUtil) CreateToken(ctx context.Context, auth *model.Auth) (string, string, error) {
	jti := uuid.New().String()

	exp := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"id":  auth.ID,
		"exp": exp.Unix(),
		"jti": jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", "", err
	}

	err = t.Redis.SetEx(ctx, "session:"+jti, auth.ID, 24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}

	return signedToken, jti, nil
}

func (t *TokenUtil) ParseToken(ctx context.Context, jwtToken string) (*model.Auth, string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", errors.New("invalid claims")
	}

	jti, _ := claims["jti"].(string)
	id, _ := claims["id"].(float64)

	exists, err := t.Redis.Exists(ctx, "session:"+jti).Result()
	if err != nil {
		return nil, "", err
	}
	if exists == 0 {
		return nil, "", errors.New("token revoked or expired")
	}

	return &model.Auth{ID: int(id)}, jti, nil
}

func (t *TokenUtil) RevokeToken(ctx context.Context, jti string) error {
	return t.Redis.Del(ctx, "session:"+jti).Err()
}
