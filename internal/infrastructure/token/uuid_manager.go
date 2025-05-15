package token

import "github.com/google/uuid"

type UUIDManager struct {
}

func NewUUIDManager() *UUIDManager {
	return &UUIDManager{}
}

func (m *UUIDManager) Generate() string {
	return uuid.New().String()
}

func (m *UUIDManager) Validate(token string) bool {
	_, err := uuid.Parse(token)
	if err != nil {
		return false
	}
	return true
}
