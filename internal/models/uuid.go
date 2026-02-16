package models

import "github.com/google/uuid"

type UUID string

func NewUUID() UUID {
	return UUID(uuid.New().String())
}

func (u UUID) String() string {
	return string(u)
}

func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return UUID(id.String()), nil
}
