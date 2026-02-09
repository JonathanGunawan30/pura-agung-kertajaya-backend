package util

import (
	"context"
	"errors"
	"pura-agung-kertajaya-backend/internal/model"
	"strconv"
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
		"id":   auth.ID,
		"role": auth.Role,
		"exp":  exp.Unix(),
		"jti":  jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", "", err
	}

	key := "session:" + jti

	data := map[string]string{
		"user_id":   auth.ID,
		"email":     auth.Email,
		"role":      auth.Role,
		"issued_at": strconv.FormatInt(time.Now().Unix(), 10),
	}

	if err = t.Redis.HSet(ctx, key, data).Err(); err != nil {
		return "", "", err
	}

	if err = t.Redis.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
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
	id, _ := claims["id"].(string)
	role, _ := claims["role"].(string)

	key := "session:" + jti
	exists, err := t.Redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, "", err
	}
	if exists == 0 {
		return nil, "", errors.New("session expired or revoked")
	}

	return &model.Auth{
		ID:   id,
		Role: role,
	}, jti, nil
}

func (t *TokenUtil) RevokeToken(ctx context.Context, jti string) error {
	return t.Redis.Del(ctx, "session:"+jti).Err()
}
