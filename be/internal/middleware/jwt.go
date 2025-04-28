package middleware

import (
	"app/be/internal/models"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

var (
	cachedAccessKey  string
	cachedRefreshKey string
	onceAccess       sync.Once
	onceRefresh      sync.Once
)

func GenerateAccessJWT(userID string) (string, error) {
	onceAccess.Do(func() {
		data, err := os.ReadFile("cfg.yml")
		if err == nil {
			var cfg models.JwtCfg
			if err := yaml.Unmarshal(data, &cfg); err == nil {
				cachedAccessKey = cfg.Jwt.SecKey
			}
		}
	})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(cachedAccessKey))
}
func ParseAccessJWT(tokenString string) (string, error) {
	onceAccess.Do(func() {
		data, err := os.ReadFile("cfg.yml")
		if err == nil {
			var cfg models.JwtCfg
			if err := yaml.Unmarshal(data, &cfg); err == nil {
				cachedAccessKey = cfg.Jwt.SecKey
			}
		}
	})
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cachedAccessKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	}
	return "", nil
}
func GenerateRefreshJWT(userID string) (string, error) {
	onceRefresh.Do(func() {
		data, err := os.ReadFile("cfg.yml")
		if err == nil {
			var cfg models.JwtCfg
			if err := yaml.Unmarshal(data, &cfg); err == nil {
				cachedRefreshKey = cfg.Jwt.RefreshKey
			}
		}
	})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(cachedRefreshKey))
}
func ParseRefreshJWT(tokenString string) (string, error) {
	onceRefresh.Do(func() {
		data, err := os.ReadFile("cfg.yml")
		if err == nil {
			var cfg models.JwtCfg
			if err := yaml.Unmarshal(data, &cfg); err == nil {
				cachedRefreshKey = cfg.Jwt.RefreshKey
			}
		}
	})
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cachedRefreshKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	}
	return "", nil
}
