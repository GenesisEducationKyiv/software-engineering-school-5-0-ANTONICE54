package token

import (
	"context"

	"github.com/google/uuid"
)

type UUIDManager struct {
}

func NewUUIDManager() *UUIDManager {
	return &UUIDManager{}
}

func (m *UUIDManager) Generate(_ context.Context) string {
	return uuid.New().String()
}

func (m *UUIDManager) Validate(_ context.Context, token string) bool {
	_, err := uuid.Parse(token)
	if err != nil {
		return false
	}
	return true
}
