package cachemanage

import (
	"github.com/easy-takeout/easy-takeout/backend/definition"
	"github.com/jinzhu/gorm"
	"gopkg.in/redis.v5"
)

type GormClient struct {
	db *gorm.DB
}

func CreateGormClient(dbtype, address string, args map[string]interface{}) (*GormClient, error) {
	db, err := definition.Connect(dbtype, address, args)
	if err != nil {
		return nil, err
	}
	return &GormClient{db}, nil
}

func (orm *GormClient) HGet(key, field string) *redis.StringCmd {
	if key == "" || field == "" {
		return redis.NewStringResult("", redis.Nil)
	}
	var hashData definition.CacheManageHash
	orm.db.Where(&definition.CacheManageHash{Key: key}).First(&hashData)
	if hashData.ID == 0 {
		return redis.NewStringResult("", redis.Nil)
	}
	var fieldData definition.CacheManageHashField
	orm.db.Where(&definition.CacheManageHashField{
		Field:  field,
		HashId: hashData.ID,
	}).First(&fieldData)
	if fieldData.Field == "" {
		return redis.NewStringResult("", redis.Nil)
	}
	return redis.NewStringResult(fieldData.Value, nil)
}

func (orm *GormClient) HSet(key, field string, value interface{}) *redis.BoolCmd {
	if key == "" || field == "" {
		return redis.NewBoolResult(false, redis.Nil)
	}
	var hashData definition.CacheManageHash
	orm.db.FirstOrCreate(&hashData, definition.CacheManageHash{Key: key})

	var fieldData definition.CacheManageHashField
	orm.db.FirstOrInit(&fieldData, definition.CacheManageHashField{
		Field:  field,
		HashId: hashData.ID,
	})
	result := true
	if fieldData.ID == 0 {
		result = false
	}
	fieldData.Value = value.(string)
	orm.db.Save(&fieldData)
	return redis.NewBoolResult(result, nil)
}

func (orm *GormClient) Close() error {
	return orm.db.Close()
}
