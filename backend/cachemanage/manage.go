package cachemanage

import (
	"github.com/satori/go.uuid"
	"gopkg.in/redis.v5"
)

const (
	loginSessionKey = "login"
)

type Manage struct {
	client *Client
}

func Create(address string) *Manage {
	client := CreateClient(address)
	return &Manage{client}
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
