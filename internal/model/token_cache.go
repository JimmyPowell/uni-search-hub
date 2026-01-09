package model

import (
	"fmt"
	"time"
	"uni-search-hub/pkg/common"
	"uni-search-hub/pkg/crypto"
	"uni-search-hub/pkg/database"
)

func cacheSetToken(token Token) error {
	key := crypto.GenerateHMAC(token.Key)
	token.Clean()
	err := database.RedisHSetObj(fmt.Sprintf("token:%s", key), &token, time.Duration(database.RedisKeyCacheSeconds())*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func cacheDeleteToken(key string) error {
	key = crypto.GenerateHMAC(key)
	err := database.RedisDelKey(fmt.Sprintf("token:%s", key))
	if err != nil {
		return err
	}
	return nil
}

func cacheIncrTokenQuota(key string, increment int64) error {
	key = crypto.GenerateHMAC(key)
	err := database.RedisHIncrBy(fmt.Sprintf("token:%s", key), common.TokenFiledRemainQuota, increment)
	if err != nil {
		return err
	}
	return nil
}

func cacheDecrTokenQuota(key string, decrement int64) error {
	return cacheIncrTokenQuota(key, -decrement)
}

func cacheSetTokenField(key string, field string, value string) error {
	key = crypto.GenerateHMAC(key)
	err := database.RedisHSetField(fmt.Sprintf("token:%s", key), field, value)
	if err != nil {
		return err
	}
	return nil
}

// CacheGetTokenByKey 从缓存中获取 token，如果缓存中不存在，则从数据库中获取
func cacheGetTokenByKey(key string) (*Token, error) {
	hmacKey := crypto.GenerateHMAC(key)
	if !database.RedisEnabled {
		return nil, fmt.Errorf("redis is not enabled")
	}
	var token Token
	err := database.RedisHGetObj(fmt.Sprintf("token:%s", hmacKey), &token)
	if err != nil {
		return nil, err
	}
	token.Key = key
	return &token, nil
}
