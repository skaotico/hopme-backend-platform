package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implementa port.CacheRepository
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache crea una nueva instancia de RedisCache
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// BlacklistJWT almacena un JWT revocado
func (r *RedisCache) BlacklistJWT(ctx context.Context, jti string, ttlSeconds int) error {
	key := fmt.Sprintf("auth:jwt:blacklist:%s", jti)
	err := r.client.Set(ctx, key, "revoked", time.Duration(ttlSeconds)*time.Second).Err()
	return err
}

// IsJWTBlacklisted verifica si un JWT está en la blacklist
func (r *RedisCache) IsJWTBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("auth:jwt:blacklist:%s", jti)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "revoked", nil
}

// IncrLoginAttempt incrementa el contador de intentos de login
func (r *RedisCache) IncrLoginAttempt(ctx context.Context, ip string, ttlSeconds int) (int, error) {
	key := fmt.Sprintf("auth:ratelimit:login:%s", ip)
	
	// Usamos un pipeline para asegurar que el INCR y el EXPIRE se hagan de forma atómica
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	// Solo queremos hacer EXPIRE si es el primer intento o queremos extenderlo (aquí lo extendemos cada vez)
	pipe.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	return int(incr.Val()), nil
}

// GetProfile obtiene el caché del perfil del usuario
func (r *RedisCache) GetProfile(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("auth:cache:profile:%s", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // No encontrado
	}
	return val, err
}

// SetProfile almacena el perfil del usuario en caché
func (r *RedisCache) SetProfile(ctx context.Context, userID string, profileJSON string, ttlSeconds int) error {
	key := fmt.Sprintf("auth:cache:profile:%s", userID)
	err := r.client.Set(ctx, key, profileJSON, time.Duration(ttlSeconds)*time.Second).Err()
	return err
}

// DeleteProfile invalida el caché del perfil
func (r *RedisCache) DeleteProfile(ctx context.Context, userID string) error {
	key := fmt.Sprintf("auth:cache:profile:%s", userID)
	err := r.client.Del(ctx, key).Err()
	return err
}
