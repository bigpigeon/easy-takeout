package cachemanage

import (
	"github.com/satori/go.uuid"
	"gopkg.in/redis.v5"
)

const (
	loginSessionKey = "login"
)

/*
* to hard to implement redis.Cmdable, so used  Client to replace
 */
type Manage struct {
	client Client
}

type Client interface {
	HSet(key, field, value string) *redis.BoolCmd
	HGet(key, field string) *redis.StringCmd
	Close() error
}

func (m *Manage) Close() error {
	return m.client.Close()
}

func Create(dbtype, address string, args map[string]interface{}) (*Manage, error) {
	var client Client
	var err error
	switch dbtype {
	case "redis":
		client, err = CreateRedisClient(address, args["password"].(string), args["db"].(int))
	default:
		client, err = CreateGormClient(dbtype, address, args)
	}
	if err != nil {
		return nil, err
	}
	return &Manage{client}, nil
}

func (m *Manage) GenerateLoginSession(name string) (string, error) {
	session := uuid.NewV4()
	err := m.client.HSet(loginSessionKey, name, string(session.Bytes())).Err()
	if err != nil {
		return "", err
	}
	return session.String(), nil
}

func (m *Manage) GetLoginSession(name string) (string, error) {
	sessionBytes, err := m.client.HGet(loginSessionKey, name).Bytes()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	session, err := uuid.FromBytes(sessionBytes)
	if err != nil {
		return "", err
	}

	return session.String(), nil
}
